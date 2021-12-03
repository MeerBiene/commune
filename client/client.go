package client

import (
	"commune/config"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"commune/gomatrix"

	"github.com/go-chi/chi"
	"github.com/go-redis/redis"
	"github.com/gorilla/sessions"
	"github.com/robfig/cron/v3"

	syncmatrix "commune/client/sync"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Client struct {
	Config      *config.Config
	Router      *chi.Mux
	HTTP        *http.Server
	Templates   *Template
	Sessions    *sessions.CookieStore
	Store       *redis.Client
	Matrix      *gomatrix.Client
	DefaultUser User
	DB          *DB
	MatrixDB    *DB
	Cron        *cron.Cron
	SyncHub     *syncmatrix.Hub
}

func (c *Client) Activate() {

	log.Println("started server")

	idleConnsClosed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)

		signal.Notify(sigint, os.Interrupt)
		signal.Notify(sigint, syscall.SIGTERM)

		<-sigint

		if err := c.HTTP.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
			log.Printf("Shutdown by user")
		}
		close(idleConnsClosed)
	}()

	if err := c.HTTP.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

type StartRequest struct {
	Config string
}

var CONFIG_FILE string

func Start(s *StartRequest) {

	CONFIG_FILE = s.Config

	db, err := NewDB("")
	if err != nil {
		panic(err)
	}

	matrixDB, err := NewDB("matrix")
	if err != nil {
		panic(err)
	}

	conf, err := config.Read(s.Config)
	if err != nil {
		panic(err)
	}

	tmpl, err := NewTemplate()
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()

	redis := redis.NewClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       conf.Redis.DB,
	})

	server := fmt.Sprintf(`http://%s:%d`, conf.Matrix.Server, conf.Matrix.Port)
	matrix, err := gomatrix.NewClient(server, "", "")
	if err != nil {
		panic(err)
	}

	//log into the default matrix account
	defUser := User{}

	username := conf.Matrix.Username
	password := conf.Matrix.Password

	resp, err := matrix.Login(&gomatrix.ReqLogin{
		Type:     "m.login.password",
		User:     username,
		Password: password,
	})

	if resp != nil {
		defUser.UserID = resp.UserID
		defUser.AccessToken = resp.AccessToken
		matrix.SetCredentials(resp.UserID, resp.AccessToken)
	}

	//default account doesn't exist yet, let's create it
	if err != nil {
		log.Println(err)

		av, err := matrix.RegisterAvailable(&gomatrix.ReqRegisterAvailable{
			Username: username,
		})
		if err != nil {
			log.Println(err)
		}

		log.Println(av)

		if av == nil || !av.Available {
			panic(err)
		}

		type Auth struct {
			Type    string
			Session string
			Admin   bool
		}

		server := fmt.Sprintf(`http://%s:%d`, conf.Matrix.Server, conf.Matrix.Port)
		matrix.Prefix = "/_synapse/admin/v1"

		nonce, err := GetNonce(server)
		if err != nil {
			log.Println(err)
		}
		log.Println("nonce is ", nonce)

		log.Println(resp)
		//actually register the user
		mac, err := ConstructMac(&NewUser{
			Username: username,
			Password: password,
			Admin:    true,
		}, nonce, conf.Auth.SharedSecret)
		if err != nil {
			log.Println(err)
			panic(nil)
		}

		req := &gomatrix.ReqLegacyRegister{
			Username: username,
			Password: password,
			Type:     "org.matrix.login.shared_secret",
			Mac:      mac,
			Admin:    true,
			Nonce:    nonce,
		}

		log.Println("what the hell is usenrmae?", username)

		resp, ui, err := matrix.LegacyRegister(req)

		log.Println("is resp nil?", resp == nil)
		log.Println("is ui nil?", ui == nil)
		log.Println("is err nil?", err == nil)

		if err != nil || resp == nil {
			log.Println(err)
		}

		if resp != nil {
			log.Println(resp)
		}

		if ui != nil {
			log.Println(ui)
		}

		defUser.UserID = resp.UserID
		defUser.AccessToken = resp.AccessToken
		matrix.SetCredentials(resp.UserID, resp.AccessToken)

		matrix.Prefix = "/_matrix/client/r0"

	}

	sess := NewSession(conf.Client.SecureCookie)
	sess.Options.Domain = fmt.Sprintf(`.%s`, conf.Client.Domain)

	hub := syncmatrix.NewHub()

	cron := cron.New()

	c := &Client{
		DB:       db,
		MatrixDB: matrixDB,
		Config:   conf,
		Matrix:   matrix,
		HTTP: &http.Server{
			ReadTimeout:  21 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  120 * time.Second,
			Addr:         fmt.Sprintf(`:%d`, conf.Client.Port),
			Handler:      router,
		},
		Router:      router,
		Templates:   tmpl,
		Sessions:    sess,
		Store:       redis,
		DefaultUser: defUser,
		Cron:        cron,
		SyncHub:     hub,
	}

	c.Middleware()
	c.Routes()

	go c.SyncHub.Run()

	c.Build()

	c.Setup()

	//go c.Cron.AddFunc("*/15 * * * *", c.RefreshCache)
	//go c.Cron.Start()

	c.Activate()
}

func (c *Client) RefreshApp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if c.Config.Mode != "production" {
			log.Println("PRODUCTION MODE")
			w.Write([]byte(RandomString(64)))
			return
		}
		w.Write([]byte(RandomString(64)))
	}
}
