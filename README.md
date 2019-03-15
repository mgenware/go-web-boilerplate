# go-triton

<img src="./assets/img/triton.jpg" width="300" height="300"/>

A boilerplate template for Go web applications. Uses Go 1.11 modules.

* Configuration file support.
* Development/production mode (via `config.DevMode`).
* Implemented common HTTP handlers:
  * Not found(404) handler.
  * Panic recovery handler as 500 Internal Server Error.
* Template support (auto reloads template in development mode).
* Auto serves static files in development mode.
* i18n support.

## Main Dependencies
* `github.com/go-chi/chi`: HTTP routing. 
* `github.com/mgenware/go-packagex`: for common helpers like template wrapper, MIME type definitions, etc.
* `golang.org/x/text/language`: HTTP `Accept-Language` header parsing and matching.
* `github.com/uber-go/zap`: Logging.

## Usage
Start in development mode:
```sh
# Start with ./config/dev.json
go run main.go dev
```

Start in production mode:
```sh
# Start with ./config/prod.json
go run main.go prod
```

The two commands simply load a configuration file by the given name, you can also create your own config file like `./config/myName.json` and start the app with it:
```sh
go run main.go myName
```

Or use the `--config` argument to specify a file:
```sh
go run main.go --config /etc/my_server/dev.json
```

## Directory Structure
```
├── appdata             Application generated files, e.g. logs, git ignored
│   └── log
├── assets              Static assets, HTML/JavaScript/CSS/Image files
├── localization        Localization resources
│   └── langs               Localized strings used by your app
├── src                 Go source directory
│   ├── app                 Core app modules, such as template manager, logger, etc.
│   ├── config              Config files
│   │   ├── dev.json
│   │   └── prod.json
│   ├── r               Routes
└── templates           Go HTML template files 
```

### The `r` Directory
This `r` directory contains all routes of your application, and because it will be so commonly used so we shorten in to `r`. In order to follow best practices for package naming ([details](https://blog.golang.org/package-names)), child directories of `r` usually consist of a short name plus a letter indicating the type of the route, e.g. `sysh` for system handlers, `homep` for home page stuff, etc.

## Error handling in HTTP handler
Two styles of error handling are supported, "panic" style and "return" style.

```go
// "panic" style
func formAPI1(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		// `panic` with a string to indicate an expected error (or user error)
		// Expected errors are not logged and served with normal 200 HTTP code.
		panic("The argument \"id\" cannot be empty")
	}
	result, err := systemCall()
	if err != nil {
		// `panic` with an error to indicate an unexpected error (or app error)
		// Unexpected errors are considered fatal and are usually because something went wrong in your application or system, so they are logged and served with 500 (Internal Server Error) code.
		panic(err)
	}
	resp := app.JSONResponse(w, r)
	resp.MustComplete(result)
}

// "return" style
func formAPI2(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	resp := app.JSONResponse(w, r)
	if id == "" {
		// `panic` with a string to indicate an expected error (or user error)
		// Expected errors are not logged and served with normal 200 HTTP code.
		resp.MustFailWithUserError("The argument \"id\" cannot be empty")
		// DON'T FORGET THE `return`
		return
	}
	result, err := systemCall()
	if err != nil {
		// For unexpected error (or app error), call `MustFail`.
		resp.MustFail(err)
		// DON'T FORGET THE `return`
	}
	resp.MustComplete(result)
}
```

## Projects built from go-trion
* [qing](https://github.com/mgenware/qing)