package client

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/hostrouter"
	"github.com/lpar/gzipped"
	dendriteUserApi "github.com/matrix-org/dendrite/userapi/api"
	"github.com/unrolled/secure"
)

func (c *Client) Routes() {

	compressor := middleware.NewCompressor(5, "text/html", "text/css")
	compressor.SetEncoder("nop", func(w io.Writer, _ int) io.Writer {
		return w
	})

	//c.Router.Use(middleware.ThrottleBacklog(10, 50, time.Second*10))
	c.Router.Use(middleware.RequestID)
	c.Router.Use(middleware.RealIP)
	c.Router.Use(middleware.Logger)
	c.Router.Use(c.Recoverer)
	c.Router.Use(middleware.StripSlashes)
	c.Router.Use(compressor.Handler)

	c.CORS()
	c.ServeStaticFiles()

	hr := hostrouter.New()

	ad := fmt.Sprintf(`%s:%d`, c.Config.Client.Domain, c.Config.Client.Port)
	sad := fmt.Sprintf(`%s:%d`, c.Config.Client.ShortlinkDomain, c.Config.Client.Port)

	if c.Config.Mode == "production" {
		ad = c.Config.Client.Domain
		sad = c.Config.Client.ShortlinkDomain
	}

	hr.Map(sad, shortlinkRoutes(c))
	hr.Map(ad, routes(c))
	//local dev please ignore
	hr.Map("192.168.1.12:8989", routes(c))

	c.Router.Mount("/", hr)

}

func shortlinkRoutes(c *Client) chi.Router {

	r := chi.NewRouter()

	r.Get("/", c.ShortlinkIndex())
	r.Get("/resolve_domain", c.ResolveShortlinkDomain())
	r.Get("/{id}", c.Shortlink())

	return r
}

func routes(c *Client) chi.Router {

	sop := secure.Options{
		ContentSecurityPolicy: "script-src 'self' 'unsafe-eval' 'unsafe-inline' $NONCE",
		IsDevelopment:         false,
		AllowedHosts: []string{
			fmt.Sprintf(`%s:%d`, c.Config.Client.Domain, c.Config.Client.Port),
			fmt.Sprintf(`%s`, c.Config.Matrix.Server),
			"192.168.1.12:8989",
		},
	}

	if c.Config.Mode == "production" {
		sop.AllowedHosts = []string{c.Config.Client.Domain}
	}

	secureMiddleware := secure.New(sop)

	r := chi.NewRouter()
	r.Use(c.GetLoggedInUser)
	r.Use(c.WASMmimeType)

	r.Route("/api/v0", func(r chi.Router) {
		r.Get("/", c.NotFound)
		r.Post("/login", c.APILogin())
	})

	r.Route("/refresh", func(r chi.Router) {
		r.Get("/", c.NotFound)
		r.Post("/", c.RefreshApp())
		//r.Post("/space/available", c.SpaceUsernameAvailable())
	})

	r.Route("/username", func(r chi.Router) {
		r.Get("/", c.NotFound)
		r.Post("/available", c.UsernameAvailable())
		//r.Post("/space/available", c.SpaceUsernameAvailable())
	})

	r.Route("/signup", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		r.Get("/", c.Error)
		r.Post("/verification/validate", c.ValidateVerificationCode())
		//r.Post("/complete", c.CompleteSignup())
		r.Post("/complete", c.CompleteSignupTemp())
		r.Post("/verification", c.ValidateSignup())
	})

	r.Route("/invite", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		r.Get("/", c.Error)
		r.Post("/valid", c.ValidateInviteCode())
	})

	r.Route("/oauth2", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		r.Route("/discord", func(r chi.Router) {
			r.Get("/authorize", c.HandleDiscordAuthorized())
		})
	})

	r.Route("/login", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		r.Get("/", c.NotFound)
		r.Post("/", c.Login())
		r.Post("/validate", c.ValidateLogin())
	})

	r.Route("/password", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		r.Use(c.GuestsOnly)
		r.Post("/verification", c.VerifyPassword())
		r.Post("/verification/validate", c.ValidatePasswordVerificationCode())
		r.Post("/reset", c.CompletePasswordReset())
	})

	r.Route("/logout", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		//r.Use(c.RequireAuthentication)
		r.Get("/", c.Logout())
	})

	r.Route("/server", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		r.Use(c.RequireAuthentication)
		r.Post("/create", c.ValidateRoomCreation())
		r.Post("/purge", c.PurgeRoom())
	})

	r.Route("/topic", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		r.Use(c.RequireAuthentication)
		r.Post("/create", c.CreateTopicRoom())
	})

	r.Route("/sync", func(r chi.Router) {
		r.Get("/{token}", c.syncConnector())
	})

	r.Route("/link", func(r chi.Router) {
		r.Post("/metadata", c.LinkMetadata())
	})

	r.Route("/", func(r chi.Router) {
		r.Use(secureMiddleware.Handler)
		r.Get("/favicon.ico", c.NotFound)
		r.Get("/about", c.StaticPage())
		r.Get("/*", c.Index())
		//r.Get("/*", c.Dispatch())
	})

	compressor := middleware.NewCompressor(5, "text/html", "text/css")
	compressor.SetEncoder("nop", func(w io.Writer, _ int) io.Writer {
		return w
	})
	r.NotFound(c.NotFound)

	return r
}

type UserSession struct {
	CreatedAt time.Time
	User      struct {
		Account *dendriteUserApi.Account
		Device  *dendriteUserApi.Device
	}
}

func (c *Client) NotFound(w http.ResponseWriter, r *http.Request) {

	us := LoggedInUser(r)
	type NotFoundPage struct {
		LoggedInUser interface{}
		Nonce        string
	}

	nonce := secure.CSPNonce(r.Context())
	pg := NotFoundPage{
		LoggedInUser: us,
		Nonce:        nonce,
	}
	c.Templates.ExecuteTemplate(w, "not-found", pg)
}

func (c *Client) ServeStaticFiles() {
	path := "/static"
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}

	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "static")
	fs := http.StripPrefix(path, gzipped.FileServer(FileSystem{http.Dir(filesDir)}))

	if path != "/" && path[len(path)-1] != '/' {
		c.Router.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	c.Router.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "max-age=31536000")
		fs.ServeHTTP(w, r)
	}))
}

type FileSystem struct {
	fs http.FileSystem
}

func (nfs FileSystem) Open(path string) (http.File, error) {
	f, err := nfs.fs.Open(path)
	if err != nil {
		return nil, err
	}

	s, err := f.Stat()
	if s.IsDir() {
		index := strings.TrimSuffix(path, "/") + "/index.html"
		if _, err := nfs.fs.Open(index); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (c *Client) CORS() {
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Image", "Attachment", "File-Type", "Content-Type", "X-CSRF-Token", "Access-Control-Allow-Origin"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})
	c.Router.Use(cors.Handler)
}

func (c *Client) Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {

				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}

				c.Error(w, r)
				return
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (c *Client) Error(w http.ResponseWriter, r *http.Request) {

	us := LoggedInUser(r)

	type errorPage struct {
		LoggedInUser interface{}
		Nonce        string
	}

	nonce := secure.CSPNonce(r.Context())
	pg := errorPage{
		LoggedInUser: us,
		Nonce:        nonce,
	}

	c.Templates.ExecuteTemplate(w, "error", pg)
}

func (c *Client) RoomTooLarge(w http.ResponseWriter, r *http.Request) {

	us := LoggedInUser(r)

	type errorPage struct {
		LoggedInUser interface{}
		Nonce        string
	}

	nonce := secure.CSPNonce(r.Context())
	pg := errorPage{
		LoggedInUser: us,
		Nonce:        nonce,
	}

	c.Templates.ExecuteTemplate(w, "room-too-large", pg)
}

func (c *Client) StaticPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := LoggedInUser(r)

		url := strings.TrimLeft(r.URL.Path, "/")

		type page struct {
			LoggedInUser interface{}
			Nonce        string
		}
		nonce := secure.CSPNonce(r.Context())

		pg := page{
			LoggedInUser: us,
			Nonce:        nonce,
		}
		c.Templates.ExecuteTemplate(w, url, pg)
	}
}
