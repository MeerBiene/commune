package client

import (
	"encoding/json"
	"log"
	"net/http"
)

func (c *Client) PurgeRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			RoomID string `json:"room_id"`
		}

		var pay payload
		if r.Body == nil {
			log.Println(err)
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&pay)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		log.Println("recieved payload ", pay)

		type Response struct {
			Deleted bool `json:"deleted"`
		}

		matrix, err := c.TempMatrixClient(c.DefaultUser.UserID, c.DefaultUser.AccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}
		matrix.Prefix = "/_synapse/admin/v1"

		ff := Response{
			Deleted: false,
		}

		jr, err := matrix.DeleteRoom(pay.RoomID)
		if err != nil {
			log.Println(err)
		} else if jr != nil {
			log.Println(jr)
			ff.Deleted = true
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
