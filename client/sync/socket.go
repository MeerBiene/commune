package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	gomatrix "commune/gomatrix"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

func ServeWs(user *User, hub *Hub, w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		conn: conn,
		send: make(chan []byte),
		user: user,
	}
	user.Drop = make(chan struct{})
	log.Println("user is", user)
	log.Println("client is", client)
	hub.register <- client

	log.Println("did we register ??")

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump(hub)
	go client.readPump(user.Token, hub)
}

type User struct {
	UserID            string `json:"user_id"`
	Token             string `json:"token"`
	MatrixAccessToken string `json:"matrix_access_token"`
	MatrixServer      string `json:"matrix_server"`
	Matrix            *gomatrix.Client
	Drop              chan struct{}
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	conn *websocket.Conn
	id   string
	send chan []byte
	user *User
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump(thing string, hub *Hub) {
	defer func() {
		hub.unregister <- thing
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		log.Println("\n")
		log.Println("reading message", string(message), "oink")
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump(hub *Hub) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			log.Println("lol: ", string(message))

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type Message struct {
	id   string
	data []byte
}

type Hub struct {
	// Registered clients.
	clients map[string]*Client

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan string

	send chan Message

	NewEvent chan *EventConstruct
}

type EventConstruct struct {
	Token string `json:"token"`
	ID    string `json:"id"`
	Event []byte `json:"event"`
}

func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan string),
		clients:    make(map[string]*Client),
		NewEvent:   make(chan *EventConstruct),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client.user.Token] = client

			log.Println("\n")
			log.Println("registering..", client.user)

			matrix, err := gomatrix.NewClient(client.user.MatrixServer, client.user.UserID, client.user.MatrixAccessToken)
			if err != nil {
				log.Println(err)
			}

			syncer := matrix.Syncer.(*gomatrix.DefaultSyncer)

			eventTypes := []string{
				"m.room.message",
				"m.reaction",
			}

			for _, x := range eventTypes {
				syncer.OnEventType(x, func(ev *gomatrix.Event) {
					fmt.Println("Message: ", ev)

					js, err := json.Marshal(ev)
					if err != nil {
						log.Println(err)
					}

					if err == nil {
						h.NewEvent <- &EventConstruct{
							Token: client.user.Token,
							ID:    ev.ID,
							Event: js,
						}
					}

				})
			}

			client.user.Matrix = matrix

			go func() {
				for {
					// Optional: Wait a period of time before trying to sync again.
					select {
					case <-client.user.Drop:
						return
					default:
						if err := client.user.Matrix.Sync(); err != nil {
							fmt.Println("Sync() returned ", err)
						}
					}

				}
			}()

			//h.users[client.id] = append(h.users[client.id], client.user)

		case thing := <-h.unregister:
			if client, ok := h.clients[thing]; ok {
				log.Println("\n")
				client.user.Matrix.StopSync()
				client.user.Drop <- struct{}{}
				delete(h.clients, thing)
				log.Println("total users is", len(h.clients))
			}
		case nc := <-h.NewEvent:
			log.Println("we go teven thhehe")
			if _, ok := h.clients[nc.Token]; ok {
				h.clients[nc.Token].send <- nc.Event
			}
		}
	}
}
