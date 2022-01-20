package gem

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

func init() {
	log.SetPrefix("[gem-spaas] ")
	log.SetFlags(log.LstdFlags)
}

type Gem struct {
	conns []*websocket.Conn
	ch    chan interface{}
}

func New() *Gem {
	g := Gem{
		ch: make(chan interface{}),
	}
	go g.dispatch()
	return &g
}

func (g *Gem) Register(c *websocket.Conn) {
	g.conns = append(g.conns, c)
}

func (g *Gem) Dispatch(in interface{}, items []Item) {
	c := struct {
		Input interface{} `json:"input"`
		Items []Item      `json:"plan"`
	}{
		Input: in,
		Items: items,
	}
	select {
	case g.ch <- c:
		// do nothing, message pass to the channel
	case <-time.After(time.Millisecond * 100):
		// skip this message
	}
}

func (g *Gem) dispatch() {
	for msg := range g.ch {
		for i := 0; i < len(g.conns); {
			err := g.conns[i].WriteJSON(msg)
			if err == nil {
				i++
				continue
			}
			// any errors, close the connection and remove the conn from the Gem
			g.conns[i].Close()
			g.conns = append(g.conns[:i], g.conns[i+1:]...)
		}
	}
}

func GetPlan(g *Gem) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		in := struct {
			Load   float64      `json:"load"`
			Fuels  Fuels        `json:"fuels"`
			Plants []PowerPlant `json:"powerplants"`
		}{}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, err.Error())
			return
		}
		plan, err := Compute(in.Load, in.Fuels, in.Plants)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return
		}
		go g.Dispatch(in, plan)

		w.Header().Set("Content-Type", "application/json")
		enc := json.NewEncoder(w)
		enc.SetIndent("", "\t")
		if err = enc.Encode(plan); err != nil {
			log.Printf("fail to encode response: %s", err)
		}
	}
	return http.HandlerFunc(fn)
}

func GetWebsocket(g *Gem) http.Handler {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	fn := func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("opening websocket fail: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, err.Error())
			return
		}
		g.Register(conn)
		for {
			// busy loop to maintain connection open - probably there is a better way
			// to do it
			_, _, err := conn.NextReader()
			if err != nil {
				log.Printf("unexpected error from client connection: %s", err)
				break
			}
		}
	}
	return http.HandlerFunc(fn)
}
