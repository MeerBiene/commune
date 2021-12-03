package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

func (c *Client) Shortlink() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		id := chi.URLParam(r, "id")

		w.Write([]byte(id))
	}
}

func (c *Client) ShortlinkIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		url := fmt.Sprintf("http://%s:%d", c.Config.Client.Domain, c.Config.Client.Port)
		if c.Config.Mode == "production" {
			url = fmt.Sprintf("https://%s", c.Config.Client.Domain)
		}

		http.Redirect(w, r, url, http.StatusSeeOther)
	}
}

func (c *Client) ResolveShortlinkDomain() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		type Response struct {
			Server string `json:"server"`
		}

		ff := Response{
			Server: c.Config.Matrix.FederationServer,
		}

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}
