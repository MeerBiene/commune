package client

import (
	"commune/gomatrix"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/unrolled/secure"
)

func (c *Client) CreateRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		us := LoggedInUser(r)

		type page struct {
			BasePage
			UserExists bool
		}

		nonce := secure.CSPNonce(r.Context())

		t := &page{}

		t.Nonce = nonce
		t.LoggedInUser = us

		c.Templates.ExecuteTemplate(w, "create", t)
	}
}

func (c *Client) SpaceUsernameAvailable() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Username string `json:"username"`
			Room     bool   `json:"room"`
			SpaceID  string `json:"space_id"`
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
			Available bool `json:"available"`
		}
		ff := Response{Available: false}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		username := strings.ToLower(pay.Username)

		/*
			if pay.Room {
				state, err := matrix.RoomState(pay.SpaceID)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), 400)
					return
				}

				alias := c.CanonicalAliasFromState(state)

				alias = alias[1:]
				s := strings.Split(alias, ":")
				localPart := s[0]
				username = fmt.Sprintf(`%s_%s`, localPart, strings.ToLower(username))
			}
		*/

		canon := fmt.Sprintf(`#%s:%s`, username, c.Config.Matrix.FederationServer)

		if user.Federated {
			canon = fmt.Sprintf(`#%s:%s`, username, user.HomeServer)
		}

		av, err := matrix.ResolveAlias(canon)

		if av == nil {
			ff.Available = true
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

func (c *Client) ValidateRoomCreation() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		type payload struct {
			Username     string        `json:"username"`
			Title        string        `json:"title"`
			About        string        `json:"about"`
			Type         string        `json:"type"`
			Private      bool          `json:"private"`
			NSFW         bool          `json:"nsfw"`
			Room         bool          `json:"room"`
			DM           bool          `json:"dm"`
			DMUsers      []string      `json:"dm_users"`
			ServerID     string        `json:"server_id"`
			RoomID       string        `json:"room_id"`
			Page         bool          `json:"page"`
			Avatar       string        `json:"avatar"`
			Thread       bool          `json:"thread"`
			ThreadEvents []interface{} `json:"thread_events"`
			ExpireThread string        `json:"expire_thread"`
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

		type NewRoom struct {
			RoomID          string            `json:"room_id"`
			Alias           string            `json:"alias"`
			ShortAlias      string            `json:"short_alias"`
			DefaultBoard    *NewRoom          `json:"default_board"`
			DefaultChatroom *NewRoom          `json:"default_chatroom"`
			Streams         map[string]string `json:"streams"`
		}

		type Response struct {
			Created bool    `json:"created"`
			Room    NewRoom `json:"room,omitempty"`
		}

		ff := Response{
			Created: false,
		}

		matrix, err := c.TempMatrixClient(user.UserID, user.MatrixAccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		fl := genSonyflake()
		username := strconv.FormatUint(fl, 10)

		/*
			sp := strings.Split(pay.Title, " ")
			jp := strings.Join(sp, "-")
			lp := strings.ToLower(jp)
			reg := regexp.MustCompile("[^abcdefghijklmnopqrstuvwxyz0123456789-]+")
			username := reg.ReplaceAllString(lp, "")
		*/

		//canon := fmt.Sprintf(`#%s:%s`, username, c.Config.Client.Domain)

		/*
			var state []*gomatrix.Event
			if pay.Room {
				state, err = matrix.RoomState(pay.ServerID)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), 400)
					return
				}

				alias = c.CanonicalAliasFromState(state)
				s := strings.Split(alias, ":")
				localPart := s[0]
				username := fmt.Sprintf(`%s_%s`, localPart, username)
				canon = fmt.Sprintf(`#%s:%s`, username, c.Config.Client.Domain)
			}
		*/

		//create the room

		users := map[string]interface{}{
			user.UserID: 100,
		}

		if pay.DM && len(pay.DMUsers) > 0 {
			for _, id := range pay.DMUsers {
				users[id] = 100
			}
		}

		pl := gomatrix.Event{
			Type: "m.room.power_levels",
			Content: map[string]interface{}{
				"ban": 60,
				"events": map[string]interface{}{
					"m.room.name":         60,
					"m.room.power_levels": 100,
					"m.room.create":       10,
					"m.space.child":       10,
					"m.space.parent":      10,
				},
				"events_default": 10,
				"invite":         10,
				"kick":           60,
				"notifications": map[string]interface{}{
					"room": 20,
				},
				"redact":        10,
				"state_default": 10,
				"users":         users,
				"users_default": 10,
			},
		}

		historyVisibility := "joined"

		if pay.Room {
			historyVisibility = "world_readable"
		}

		initState := []gomatrix.Event{
			gomatrix.Event{
				Type: "m.room.history_visibility",
				Content: map[string]interface{}{
					"history_visibility": historyVisibility,
				},
			}, gomatrix.Event{
				Type: "m.room.guest_access",
				Content: map[string]interface{}{
					"guest_access": "can_join",
				},
			}, gomatrix.Event{
				Type: "commune.room",
				Content: map[string]interface{}{
					"room_type": "space",
				},
			}, gomatrix.Event{
				Type: "m.room.type",
				Content: map[string]interface{}{
					"type":  "m.space",
					"alias": username,
				},
			},
			pl,
		}

		if pay.DM && len(pay.DMUsers) > 0 {

			initState = []gomatrix.Event{
				gomatrix.Event{
					Type: "m.room.history_visibility",
					Content: map[string]interface{}{
						"history_visibility": "world_readable",
					},
				}, gomatrix.Event{
					Type: "m.room.guest_access",
					Content: map[string]interface{}{
						"guest_access": "can_join",
					},
				}, gomatrix.Event{
					Type: "commune.room",
					Content: map[string]interface{}{
						"room_type": "dm",
					},
				}, gomatrix.Event{
					Type: "m.room.type",
					Content: map[string]interface{}{
						"type":  "m.space",
						"alias": username,
					},
				},
				pl,
			}
		}

		if len(pay.Avatar) > 0 {
			avev := gomatrix.Event{
				Type: "m.room.avatar",
				Content: map[string]interface{}{
					"url": pay.Avatar,
				},
			}
			initState = append(initState, avev)
		}

		if pay.Room {

			roomType := "commune.room"

			if pay.Thread {
				roomType = "commune.room.thread"
			}

			ev := gomatrix.Event{
				Type: roomType,
				Content: map[string]interface{}{
					"room_type": pay.Type,
				},
			}

			if pay.Thread {
				ev.Content["thread_in_room_id"] = pay.RoomID
			}

			initState = append(initState, ev)

			if pay.ThreadEvents != nil && len(pay.ThreadEvents) > 0 {
				for _, event := range pay.ThreadEvents {
					x := RandomString(32)
					ev := gomatrix.Event{
						Type:     "commune.room.thread.initial",
						StateKey: &x,
						Content: map[string]interface{}{
							"event": event,
						},
					}
					initState = append(initState, ev)
				}
			}

			/*
				alias = alias[1:]
				s := strings.Split(alias, ":")
				localPart := s[0]
				username = fmt.Sprintf(`%s_%s`, localPart, username)

				x := strings.Split(localPart, "_")
				j := strings.Join(x, "/")

				path = fmt.Sprintf(`%s/%s`, j, username)
			*/

			content := map[string]interface{}{
				"canonical_alias": fmt.Sprintf(`#%s:%s`, username, user.HomeServer),
				"local_part":      username,
				"via":             []string{c.Config.Client.Domain},
			}

			if pay.Page {
				content["page"] = true
			}

			initState = append(initState, gomatrix.Event{
				Type:     `m.space.parent`,
				StateKey: &pay.ServerID,
				Content:  content,
			})
		}

		if pay.NSFW {
			initState = append(initState, gomatrix.Event{
				Type: "commune.room.nsfw",
				Content: map[string]interface{}{
					"nsfw": true,
				},
			})
		}

		req := &gomatrix.ReqCreateRoom{
			Preset:     "public_chat",
			Visibility: "public",
			CreationContent: map[string]interface{}{
				"m.federate": true,
			},
			InitialState: initState,
		}

		if pay.DM && len(pay.DMUsers) > 0 {
			req.IsDirect = true
			req.Preset = "private_chat"
			req.Visibility = "private"
			for _, id := range pay.DMUsers {
				req.Invite = append(req.Invite, id)
			}
		} else {
			req.RoomAliasName = username
			req.Name = pay.Title
			req.Topic = pay.About
		}

		/*
			if pay.Room {
				req.RoomAliasName = username
				req.Name = strings.ToLower(pay.Title)

				members, err := matrix.JoinedMembers(pay.ServerID)

				if err != nil {
					log.Println(err)
				}

				if members != nil {
					for userID, _ := range members.Joined {
						if userID != user.UserID {
							req.Invite = append(req.Invite, userID)
						}
					}
				}

			}
		*/

		if pay.Private {
			req.Preset = "private_chat"
			req.Visibility = "private"
		}

		crr, err := matrix.CreateRoom(req)

		if err != nil || crr == nil {
			log.Println(err)
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		shortAlias, err := c.AddRoomShortAlias(crr.RoomID, matrix)
		if err != nil {
			log.Println(err)
		}

		ff.Created = true

		gcidg := ""
		gcag := ""

		//create #general channel

		type createStreamsRequest struct {
			RoomType string
			ParentID string
		}

		createStream := func(r createStreamsRequest) (string, error) {

			initState := []gomatrix.Event{
				gomatrix.Event{
					Type: "m.room.history_visibility",
					Content: map[string]interface{}{
						"history_visibility": "world_readable",
					},
				}, gomatrix.Event{
					Type: "m.room.guest_access",
					Content: map[string]interface{}{
						"guest_access": "can_join",
					},
				}, gomatrix.Event{
					Type: "commune.room",
					Content: map[string]interface{}{
						"room_type": r.RoomType,
					},
				},
				pl,
			}

			f := genSonyflake()
			u := strconv.FormatUint(f, 10)

			content := map[string]interface{}{
				"canonical_alias": fmt.Sprintf(`#%s:%s`, u, user.HomeServer),
				"local_part":      u,
				"via":             []string{c.Config.Client.Domain},
			}

			initState = append(initState, gomatrix.Event{
				Type:     `m.space.parent`,
				StateKey: &r.ParentID,
				Content:  content,
			})

			req := &gomatrix.ReqCreateRoom{
				Preset:        "public_chat",
				Visibility:    "public",
				RoomAliasName: u,
				CreationContent: map[string]interface{}{
					"m.federate": true,
				},
				InitialState: initState,
			}

			cr, err := matrix.CreateRoom(req)

			if err != nil || cr == nil {
				log.Println(err)
				return "", err
			}

			{
				content["room_id"] = cr.RoomID
				content["room_type"] = r.RoomType
				content["default"] = true
				_, err := matrix.SendStateEvent(r.ParentID, "m.space.child", cr.RoomID, content)
				if err != nil {
					log.Println(err)
					return "", err
				}
			}

			return cr.RoomID, nil

		}
		rid1 := ""
		rid2 := ""
		if !pay.Room && !pay.DM {

			/*
			 */

			//create te general chatroom

			newState := []gomatrix.Event{
				gomatrix.Event{
					Type: "m.room.history_visibility",
					Content: map[string]interface{}{
						"history_visibility": "world_readable",
					},
				}, gomatrix.Event{
					Type: "m.room.guest_access",
					Content: map[string]interface{}{
						"guest_access": "can_join",
					},
				}, gomatrix.Event{
					Type: "commune.room",
					Content: map[string]interface{}{
						"room_type": "chat",
					},
				}, gomatrix.Event{
					Type: "m.room.type",
					Content: map[string]interface{}{
						"type":  "m.space",
						"alias": username,
					},
				},
				pl,
			}

			s := genSonyflake()
			v := strconv.FormatUint(s, 10)

			ncontent := map[string]interface{}{
				"canonical_alias": fmt.Sprintf(`#%s:%s`, v, user.HomeServer),
				"local_part":      v,
				"via":             []string{c.Config.Client.Domain},
			}
			newState = append(newState, gomatrix.Event{
				Type:     `m.space.parent`,
				StateKey: &crr.RoomID,
				Content:  ncontent,
			})

			requ := &gomatrix.ReqCreateRoom{
				Preset:        "public_chat",
				Visibility:    "public",
				Name:          `general`,
				RoomAliasName: v,
				CreationContent: map[string]interface{}{
					"m.federate": true,
				},
				InitialState: newState,
			}

			ccr, err := matrix.CreateRoom(requ)

			if err != nil || ccr == nil {
				log.Println(err)
				log.Println(err)
			}

			rid1, err = createStream(createStreamsRequest{
				RoomType: "chat",
				ParentID: ccr.RoomID,
			})
			if err != nil {
				log.Println(err)
				log.Println("COULD NOT CREATE STREAM")
			}
			rid2, err = createStream(createStreamsRequest{
				RoomType: "topics",
				ParentID: ccr.RoomID,
			})
			if err != nil {
				log.Println(err)
				log.Println("COULD NOT CREATE STREAM")
			}

			gcidg = ccr.RoomID
			gcag = v

			go func() {
				ncontent["room_id"] = ccr.RoomID
				ncontent["name"] = `general`
				ncontent["default_stream"] = `chat`
				ncontent["streams"] = map[string]string{
					"chat":   rid1,
					"topics": rid2,
				}
				_, err := matrix.SendStateEvent(crr.RoomID, "m.space.child", ccr.RoomID, ncontent)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), 400)
					return
				}
				co := map[string]interface{}{
					"short_alias": shortAlias,
				}
				_, err = matrix.SendStateEvent(crr.RoomID, "commune.room.short_alias", ccr.RoomID, co)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), 400)
					return
				}
			}()

		}

		if pay.Room {

			if !pay.Thread {

				rid1, err = createStream(createStreamsRequest{
					RoomType: "chat",
					ParentID: crr.RoomID,
				})
				if err != nil {
					log.Println(err)
					log.Println("COULD NOT CREATE STREAM")
				}
				rid2, err = createStream(createStreamsRequest{
					RoomType: "topics",
					ParentID: crr.RoomID,
				})
				if err != nil {
					log.Println(err)
					log.Println("COULD NOT CREATE STREAM")
				}
			}

			go func() {
				content := map[string]interface{}{
					"canonical_alias": fmt.Sprintf(`#%s:%s`, username, user.HomeServer),
					"local_part":      username,
					"name":            strings.ToLower(pay.Title),
					"via":             []string{c.Config.Client.Domain},
					"default_stream":  pay.Type,
				}
				if !pay.Thread {
					content["streams"] = map[string]string{
						"chat":   rid1,
						"topics": rid2,
					}
				}

				if pay.Thread {
					content["thread"] = map[string]interface{}{
						"thread_in_room_id": pay.RoomID,
						"expire_thread":     pay.ExpireThread,
					}
				}

				_, err := matrix.SendStateEvent(pay.ServerID, "m.space.child", crr.RoomID, content)
				if err != nil {
					log.Println(err)
					http.Error(w, err.Error(), 400)
					return
				}
			}()

		}

		//canon := fmt.Sprintf(`#%s:%s`, username, user.HomeServer)

		ff.Room = NewRoom{
			RoomID:     crr.RoomID,
			Alias:      username,
			ShortAlias: shortAlias,
			DefaultChatroom: &NewRoom{
				RoomID: gcidg,
				Alias:  gcag,
			},
		}

		if pay.Room {
			ff.Room.Streams = map[string]string{
				"chat":   rid1,
				"topics": rid2,
			}
		} else {
			ff.Room.DefaultChatroom.Streams = map[string]string{
				"chat":   rid1,
				"topics": rid2,
			}
		}

		ctx := context.Background()
		ctx, _ = context.WithTimeout(ctx, 3*time.Second)

		/*

			err = c.UpdateJoinedRooms(matrix, r)
			if err != nil {
				log.Println(err)
			}
				if !pay.Page {
					alias := fmt.Sprintf(`#%s:%s`, username, user.HomeServer)
					path := username
					if strings.Contains(path, "_") {
						s := strings.Split(path, "_")
						path = strings.Join(s, "/")
					}

					c.Cache.Rooms.Set(crr.RoomID, alias, 1)

				}

				if crr != nil {
					c.OperatorJoinRoom(crr.RoomID)

					if pay.Type != "gallery" && pay.Type == "page" {
						text, html := InitialMessage()
						npe := gomatrix.CreatePostEvent{
							RoomID:        crr.RoomID,
							Text:          text,
							FormattedText: html,
						}

						_, err := matrix.CreatePost(&npe)
						if err != nil {
							log.Println(err)
							log.Println(err)
						}
					}
					err = c.RefreshRoomsCache()
					if err != nil {
						log.Println(err)
					}

				}
		*/

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}

func (c *Client) CreateTopicRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		user, err := c.GetTokenUser(token)
		if err != nil || user == nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		matrix, err := c.TempMatrixClient(c.DefaultUser.UserID, c.DefaultUser.AccessToken)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		type Response struct {
			Created bool   `json:"created"`
			RoomID  string `json:"room_id"`
		}

		ff := Response{
			Created: false,
		}

		//create the room

		users := map[string]interface{}{
			c.DefaultUser.UserID: 100,
		}

		pl := gomatrix.Event{
			Type: "m.room.power_levels",
			Content: map[string]interface{}{
				"ban": 60,
				"events": map[string]interface{}{
					"m.room.name":         60,
					"m.room.power_levels": 100,
					"m.room.create":       10,
					"m.space.child":       10,
					"m.space.parent":      10,
				},
				"events_default": 10,
				"invite":         10,
				"kick":           60,
				"notifications": map[string]interface{}{
					"room": 20,
				},
				"redact":        10,
				"state_default": 10,
				"users":         users,
				"users_default": 10,
			},
		}

		initState := []gomatrix.Event{
			gomatrix.Event{
				Type: "m.room.history_visibility",
				Content: map[string]interface{}{
					"history_visibility": "world_readable",
				},
			}, gomatrix.Event{
				Type: "m.room.guest_access",
				Content: map[string]interface{}{
					"guest_access": "can_join",
				},
			}, gomatrix.Event{
				Type: "commune.room",
				Content: map[string]interface{}{
					"room_type": "space",
				},
			},
			pl,
		}

		req := &gomatrix.ReqCreateRoom{
			Preset:     "public_chat",
			Visibility: "public",
			CreationContent: map[string]interface{}{
				"m.federate": true,
			},
			InitialState: initState,
		}

		crr, err := matrix.CreateRoom(req)

		if err != nil || crr == nil {
			log.Println(err)
			log.Println(err)
			http.Error(w, err.Error(), 400)
			return
		}

		ff.Created = true
		ff.RoomID = crr.RoomID

		js, err := json.Marshal(ff)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	}
}
