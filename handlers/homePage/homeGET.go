package homePage

import (
	"net/http"
	"time"

	"github.com/mgenware/go-triton/app"
)

var indexView = app.TemplateManager.MustParseLocalizedView("home.html")

func HomeGET(w http.ResponseWriter, r *http.Request) {
	ctx, tm, resp := app.HTMLResponse(w, r)

	pageData := &HomePageData{Time: time.Now().String()}
	pageHTML := indexView.MustExecuteToString(ctx, pageData)

	d := tm.NewMasterPageData(tm.NewLocalizedTitle(ctx, "home"), pageHTML)
	resp.MustComplete(d)
}
