package client

import (
	"fmt"
	"log"
	"net/http"

	syncmatrix "commune/client/sync"

	"github.com/go-chi/chi"
)

func (c *Client) syncConnector() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		token := chi.URLParam(r, "token")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		serverName := c.URLScheme(c.Config.Matrix.Server) + fmt.Sprintf(`:%d`, c.Config.Matrix.Port)

		us := syncmatrix.User{
			UserID:            user.UserID,
			Token:             RandomString(32),
			MatrixAccessToken: user.MatrixAccessToken,
			MatrixServer:      serverName,
		}

		syncmatrix.ServeWs(&us, c.SyncHub, w, r)
	}
}
