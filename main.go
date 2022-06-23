package main

import (
	"bytes"
	"embed"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"strconv"

	"github.com/c-goetz/traitor-card-game/lobby"
)

type Error int

const (
	ErrLobbyCreate Error = iota
	// must be last
	ErrLast
)

func errorFlash(err Error) string {
	switch err {
	case ErrLobbyCreate:
		return "Internal error creating lobby."
	default:
		return ""
	}
}

type TemplateData struct {
	Static Strings
	Flash  string
}

type LobbyTemplateData struct {
	TemplateData
	Lobby   *lobby.Lobby
	Player  *lobby.Player
	Initial bool
	LobbyId string
}

// TODO one dir, serve fs directly
//go:embed static/htmx-1.7.0-min.js
var htmx []byte

//go:embed static/htmx-1.7.0-sse.js
var htmxSse []byte

//go:embed templates
var templates embed.FS

type Strings struct {
	Title,
	Create,
	Join,
	Rules,
	Lobby,
	PlayerName string
}

func main() {
	strings := Strings{
		// TODO i18n
		Title:      "Traitor Card Game",
		Create:     "Create Lobby",
		Join:       "Join Existing Lobby",
		Rules:      "View Rules",
		Lobby:      "Lobby",
		PlayerName: "Name",
	}
	tsFS, err := fs.Sub(templates, "templates")
	if err != nil {
		log.Fatal(err)
	}
	ts := template.Must(template.ParseFS(tsFS, "*.html"))
	mux := http.NewServeMux()
	mux.HandleFunc("/static/htmx.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/javascript")
		reader := bytes.NewReader(htmx)
		io.Copy(w, reader)
	})
	mux.HandleFunc("/static/sse.js", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/javascript")
		reader := bytes.NewReader(htmxSse)
		io.Copy(w, reader)
	})
	// TODO sse
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html")
		err := r.ParseForm()
		if err != nil {
			log.Printf("parse form: %v", err)
			w.WriteHeader(400)
			return
		}
		var flash string
		errorCode := r.Form.Get("err")
		if errorCode != "" {
			parsed, err := strconv.Atoi(errorCode)
			if err != nil {
				log.Printf("parse errorCode: %v", err)
				w.WriteHeader(400)
				return
			}
			if 0 > parsed || Error(parsed) > ErrLast {
				log.Printf("errorCode out of range: %d", parsed)
				w.WriteHeader(400)
				return
			}
			flash = errorFlash(Error(parsed))
		}
		switch r.URL.Path {
		case "/":
			data := TemplateData{strings, flash}
			ts.ExecuteTemplate(w, "index.html", data)
		case "/lobby":
			id := r.Form.Get("id")
			if id == "" {
				player, lobby, err := lobby.CreateLobby("Host")
				if err != nil {
					log.Printf("can't create lobby: %v", err)
					http.Redirect(w, r, fmt.Sprintf("/?err=%d", ErrLobbyCreate), 301)
				}
				bs := make([]byte, 8)
				binary.LittleEndian.PutUint64(bs, lobby.Uuid)
				// TODO register host channel
				data := LobbyTemplateData{
					TemplateData{strings, flash},
					lobby,
					player,
					true,
					base64.RawURLEncoding.EncodeToString(bs),
				}
				ts.ExecuteTemplate(w, "lobby.html", data)
			}
			// TODO get lobby, show
		default:
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("api request: %v\n", r)
	})
	err = http.ListenAndServe(":8000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
