package client

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (c *Client) Middleware() {

	if c.Config.Mode == "development" {
		c.Router.Use(c.reloadtemplates)
	}
}

func LoggedInUser(r *http.Request) *User {
	us, ok := r.Context().Value("user").(*User)
	if !ok {
		return nil
	}

	return us
}

//Checks for logged in user on routes that use it
func (c *Client) GetLoggedInUser(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			h.ServeHTTP(w, r)
			return
		}

		token, ok := s.Values["access_token"].(string)
		if ok {
			userid, err := c.Store.Get(token).Result()
			if err != nil {
				log.Println(err)
				h.ServeHTTP(w, r)
				return
			}

			user, err := c.Store.Get(userid).Result()
			if err != nil {
				log.Println(err)
				h.ServeHTTP(w, r)
				return
			}

			var us User
			err = json.Unmarshal([]byte(user), &us)
			if err != nil {
				log.Println(err)
				h.ServeHTTP(w, r)
				return
			}

			if len(us.AvatarURL) > 0 {
				us.AvatarURL = c.BuildAvatar(us.AvatarURL)
			}

			ctx := context.WithValue(r.Context(), "user", &us)
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		h.ServeHTTP(w, r)
	})
}

//makes sure this route is autehnticated
func (c *Client) RequireAuthentication(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		token, ok := s.Values["access_token"].(string)

		if !ok {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		userid, err := c.Store.Get(token).Result()
		if err != nil {
			log.Println(err)
			h.ServeHTTP(w, r)
			return
		}

		user, err := c.Store.Get(userid).Result()
		if err != nil {
			log.Println(err)
			h.ServeHTTP(w, r)
			return
		}

		var us User
		err = json.Unmarshal([]byte(user), &us)
		if err != nil || us.UserID == "" {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		h.ServeHTTP(w, r)
	})
}

//makes sure this route is autehnticated
func (c *Client) GuestsOnly(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		s, err := GetSession(r, c)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		token, ok := s.Values["access_token"].(string)

		if ok && len(token) > 0 {
			userid, err := c.Store.Get(token).Result()
			if err == nil && userid != "" {
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
		}

		h.ServeHTTP(w, r)
	})
}

func (c *Client) reloadtemplates(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.ReloadTemplates()

		h.ServeHTTP(w, r)
	})
}

//Checks for logged in user on routes that use it
func (c *Client) WASMmimeType(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		isWASM := strings.Contains(r.URL.Path, ".wasm")

		if isWASM {
			w.Header().Set("Content-Type", "application/wasm")
		}

		log.Println("is itWASM?", isWASM)

		h.ServeHTTP(w, r)
	})
}
