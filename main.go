package main

import (
	"fmt"
	"net/http"
	//"unicode/utf8"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

var (
	config WRConfig
	upgrader = websocket.Upgrader {
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	/* Config
	****************************************************************/
	if err := config.Populate(); err != nil {
		panic(err)
	}

	fmt.Println("Server Root: " + config.root)

	/* Routes
	****************************************************************/
	wrs := mux.NewRouter()

	wrs.HandleFunc("/dj", func (w http.ResponseWriter, r *http.Request) {
		path := config.root + "/dj"

		indexPage := &Page{}

		if err := indexPage.CreateView(path + "/index.html"); err != nil {
			fmt.Println("DJ: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(indexPage.ServePage())
	}).
	Methods("GET")

	wrs.HandleFunc("/audience", func (w http.ResponseWriter, r *http.Request) {
		path := config.root + "/audience"

		indexPage := &Page{}
		if err := indexPage.CreateView(path + "/index.html"); err != nil {
			fmt.Println("Audience: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(indexPage.ServePage())
	}).
	Methods("GET")

	wrs.HandleFunc("/echo", func (w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		for {
			messageType, msg, err := conn.ReadMessage()
			if err != nil {
				return
			}

			fmt.Printf("\n\n\nReceived b:");

			/*
			for len(msg) > 0 {
				r, size := utf8.DecodeRune(msg)
				fmt.Printf("%c", r)

				msg = msg[size:]
			}
			*/

			for n := 0;n < len(msg);n++ {
				fmt.Printf("%d,",msg[n]);
			}

			fmt.Printf("\n");

			if err = conn.WriteMessage(messageType, msg); err != nil {
				return
			}
		}
	})

	/* File Server
	****************************************************************/
	// Note: Do not put any HandleFunc, MatcherFunc, etc functions below the file
	//       server function, it will cause the Handle functions to not work at all,
	//       could be a bug with Gorilla Mux, but don't really have time to test
	wrs.
	PathPrefix("/").
	Handler(http.FileServer(http.Dir(config.root))).
	Methods("GET")

	/* Serve
	****************************************************************/
	http.Handle("/", &MiddleRouter{wrs})
	if err := http.ListenAndServe(config.ip + ":" + config.port, nil); err != nil {
		panic(err)
	}
}
