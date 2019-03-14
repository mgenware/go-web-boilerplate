package template

import (
	"database/sql"
	"fmt"
	"go-triton-app/app/cfg"
	"go-triton-app/app/logx"
	"log"
	"net/http"
	"path/filepath"

	"go-triton-app/app/template/localization"

	"github.com/mgenware/go-packagex/httpx"
	"github.com/mgenware/go-packagex/templatex"
	strf "github.com/mgenware/go-string-format"
)

// Manager provides common functions to generate HTML strings.
type Manager struct {
	dir    string
	logger *logx.Logger
	config *cfg.Config

	reloadViewsOnRefresh bool
	log404Error          bool
	panicOnFatalError    bool

	masterView          *LocalizedView
	errorView           *LocalizedView
	LocalizationManager *localization.Manager
}

// MustCreateManager creates an instance of TemplateManager with specified arguments. Note that this function panics when main template loading fails.
func MustCreateManager(
	dir string,
	i18nDir string,
	defaultLang string,
	logger *logx.Logger,
	config *cfg.Config,
) *Manager {
	reloadViewsOnRefresh := config.Debug != nil && config.Debug.ReloadViewsOnRefresh
	if reloadViewsOnRefresh {
		log.Print("⚠️ View dev mode is on")
	}

	// Create the localization manager used by localized template views
	localizationManager, err := localization.NewManagerFromDirectory(i18nDir, defaultLang)
	if err != nil {
		panic(err)
	}

	t := &Manager{
		dir:                  dir,
		LocalizationManager:  localizationManager,
		config:               config,
		logger:               logger,
		reloadViewsOnRefresh: reloadViewsOnRefresh,
		log404Error:          config.HTTP.Log404Error,
		panicOnFatalError:    config.Debug != nil && config.Debug.PanicOnUnexpectedHTMLErrors,
	}

	// Load the master template
	t.masterView = t.MustParseLocalizedView("master.html")
	// Load the error template
	t.errorView = t.MustParseLocalizedView("error.html")

	return t
}

// MustCompleteWithContent finished the response with the given HTML content.
func (m *Manager) MustCompleteWithContent(content []byte, w http.ResponseWriter) {
	httpx.SetResponseContentType(w, httpx.MIMETypeHTMLUTF8)
	w.Write(content)
}

// MustComplete executes the main view template with the specified data and panics if error occurs.
func (m *Manager) MustComplete(r *http.Request, lang string, d *MasterPageData, w http.ResponseWriter) {
	httpx.SetResponseContentType(w, httpx.MIMETypeHTMLUTF8)

	// Setup additional assets, e.g.:
	// data.Header += "<link href=\"/static/main.min.css\" rel=\"stylesheet\"/>"
	// data.Scripts += "<script src=\"/static/main.min.js\"></script>"

	m.masterView.MustExecute(lang, w, d)
}

// MustError executes the main view template with the specified data and panics if error occurs.
func (m *Manager) MustError(r *http.Request, lang string, d *ErrorPageData, w http.ResponseWriter) {
	// Handle unexpected errors
	if !d.Expected {
		if d.Error == sql.ErrNoRows {
			// Consider `sql.ErrNoRows` as 404 not found error
			w.WriteHeader(http.StatusNotFound)

			d.Message = m.LocalizedString(lang, "resourceNotFound")
			if m.config.HTTP.Log404Error {
				m.logger.NotFound("sql", r.URL.String())
			}
		} else {
			// At this point, this should be a 500 server internal error
			w.WriteHeader(http.StatusInternalServerError)

			if d.Error != nil {
				d.Message = d.Error.Error()
			}
			m.logger.Error("fatal-error", "msg", d.Message)
		}
	}
	// Throw unexpected error in dev mode, note that `d.Expected` may change in this method, e.g. a unexpected `sql.errNoRows` can turn into an expected 404 error
	if !d.Expected && m.panicOnFatalError {
		fmt.Println("🙉 This message only appears in dev mode.")
		if d.Error != nil {
			panic(d.Error)
		} else {
			panic(d.Message)
		}
	}
	errorHTML := m.errorView.MustExecuteToString(lang, d)
	htmlData := NewMasterPageData("Error", errorHTML)
	m.MustComplete(r, lang, htmlData, w)
}

// PageTitle returns the given string followed by the localized site name.
func (m *Manager) PageTitle(lang, s string) string {
	return s + " - " + m.LocalizationManager.ValueForKey(lang, "_siteName")
}

// LocalizedPageTitle calls PageTitle with a localized string associated with the specified key.
func (m *Manager) LocalizedPageTitle(lang, key string) string {
	ls := m.LocalizationManager.ValueForKey(lang, key)
	return m.PageTitle(lang, ls)
}

// MustParseLocalizedView creates a new LocalizedView with the given relative path.
func (m *Manager) MustParseLocalizedView(relativePath string) *LocalizedView {
	file := filepath.Join(m.dir, relativePath)
	view := templatex.MustParseView(file, m.reloadViewsOnRefresh)
	return &LocalizedView{view: view, localizationManager: m.LocalizationManager}
}

// MustParseView creates a new View with the given relative path.
func (m *Manager) MustParseView(relativePath string) *templatex.View {
	file := filepath.Join(m.dir, relativePath)
	return templatex.MustParseView(file, m.reloadViewsOnRefresh)
}

// LocalizedString is a convenience function of LocalizationManager.ValueForKey.
func (m *Manager) LocalizedString(lang, key string) string {
	return m.LocalizationManager.ValueForKey(lang, key)
}

// FormatLocalizedString is a convenience function to format a localized string.
func (m *Manager) FormatLocalizedString(lang, key string, a ...interface{}) string {
	return strf.Format(m.LocalizedString(lang, key), a...)
}
