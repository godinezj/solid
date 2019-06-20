package actions

import (
	"bitbucket.org/godinezj/solid/log"
	"bitbucket.org/godinezj/solid/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo-pop/pop/popmw"
	"github.com/gobuffalo/envy"
	contenttype "github.com/gobuffalo/mw-contenttype"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/rs/cors"
	"github.com/unrolled/secure"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env: ENV,
			// TODO replace buffalo sessions with OAuth.
			// This line is commented out below because buffalo uses
			// null sessions for APIS by default.
			// SessionStore: sessions.Null{},
			PreWares: []buffalo.PreWare{
				cors.Default().Handler,
			},
			SessionName: "_solid_session",
		})
		// Automatically redirect to SSL
		app.Use(forceSSL())

		// injects log into app logger
		log.Log = app.Logger

		// Set the request content type to JSON
		app.Use(contenttype.Set("application/json"))

		if ENV == "development" {
			app.Use(paramlogger.ParameterLogger)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		app.GET("/", HomeHandler)

		app.POST("/users", Create)
		app.POST("/login", Login)
		app.POST("/forgot_password", GenPassResetToken)
		app.POST("/reset_password", ValidatePassResetToken)
		app.POST("/vpn/create", VPNCreate)
		app.POST("/vpn/show", VPNCreate)
	}

	return app
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
