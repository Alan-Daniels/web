package buildinggowebapps

var setupGo = `$ go mod init github.com/Alan-Daniels/web
  go: creating new go.mod: module github.com/Alan-Daniels/web`

var setupNix = `$ nix flake init -t github:nix-community/gomod2nix#app
  wrote: /home/[...]/web/.gitignore
  wrote: /home/[...]/web/default.nix
  wrote: /home/[...]/web/flake.nix
  refusing to overwrite existing file '/home/[...]/web/go.mod'
   please merge it manually with '/nix/store/cq27b0ilj9y3ld9h5wrhxjnih3bij2kd-source/templates/app/go.mod'
  wrote: /home/[...]/web/go.mod
  wrote: /home/[...]/web/gomod2nix.toml
  wrote: /home/[...]/web/main.go
  wrote: /home/[...]/web/shell.nix
  error: Encountered 1 conflicts - see above`

var setupGit = `$ git init
  Initialized empty Git repository in /home/[...]/web/.git/`

var setupAir = `$ air init
  
    __    _   ___
   / /\  | | | |_)
  /_/--\ |_| |_| \_ (devel), built with Go go1.22.5
  
  .air.toml file created to the current directory with the default settings`

var codeDevSh = `#!/usr/bin/env bash

(trap 'kill 0' SIGINT; 
air &
templ generate --watch --proxy="http://127.0.0.1:8080" --open-browser=false &
npx tailwindcss -i ./internal/input.css -o ./assets/styles.css --minify --watch
)`

var tree = `$ tree -a
  .
  ├── .air.toml
  ├── default.nix
  ├── dev.sh
  ├── flake.nix
  ├── .git
  │   └── [...]
  ├── .gitignore
  ├── assets
  │   └── favicon.ico
  ├── go.mod
  ├── gomod2nix.toml
  ├── go.sum
  ├── internal
  │   ├── commit.txt
  │   ├── home.templ
  │   └── input.css
  ├── main.go
  ├── shell.nix
  └── tools.go`

var workingWithEchoImports = `package main

import (
	"github.com/labstack/echo/v4"
	"github.com/Alan-Daniels/web/internal"
	
	[...]
)`

var workingWithEchoMain = `func main() {
	flSSL := flag.Bool("ssl", false, "whether to start with ssl")
	webroot := flag.String("root", ".", "where the files be ;)")
	flag.Parse()

	app := echo.New()
	app.Static("/assets", (*webroot)+"/assets")
	app.File("/favicon.ico", (*webroot)+"/assets/favicon.ico")
	app.GET("/", ComponentHandler(internal.Home))

	[...]

	if *flSSL {
		customHTTPServer(app, *webroot)
	} else {
		app.Logger.Fatal(app.Start(":8080"))
	}
}
`

var sslSupport = `func customHTTPServer(e *echo.Echo, webroot string) {
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	autoTLSManager := autocert.Manager{
		Prompt: autocert.AcceptTOS,
		// Cache certificates to avoid issues with
		// rate limits (https://letsencrypt.org/docs/rate-limits)
		Cache:      autocert.DirCache(webroot + "./.cache"),
		HostPolicy: autocert.HostWhitelist("your-domain.here"),
	}
	s := http.Server{
		Addr:    ":4343",
		Handler: e, // set Echo as handler
		TLSConfig: &tls.Config{
			GetCertificate: autoTLSManager.GetCertificate,
			NextProtos:     []string{acme.ALPNProto},
		},
		//ReadTimeout: 30 * time.Second, // use custom timeouts
	}
	go http.ListenAndServe(":8080", autoTLSManager.HTTPHandler(nil))
	if err := s.ListenAndServeTLS("", ""); err != http.ErrServerClosed {
		e.Logger.Fatal(err)
	}
}`

var workingWithTempl = `package internal

import _ "embed"

//go:embed commit.txt
var Commit string

templ Page() {
	<!DOCTYPE html>
	<html lang="en">
		<head>
			<meta charset="UTF-8"/>
			<title>Your Site</title>
			<link rel="stylesheet" href={ "/assets/styles.css?ver=" + Commit }/>
		</head>
		<body>
			{ children... }
		</body>
	</html>
}

templ Home() {
	@Page() {
		<h1>Hello world!</h1>
	}
}
`

var templEchoGlue = `func ComponentHandler(comp func() templ.Component) func(e echo.Context) error {
	return func(e echo.Context) error {
		return Render(e, http.StatusOK, comp())
	}
}

// This custom Render replaces Echo's echo.Context.Render()
// with templ's templ.Component.Render().
func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}`
