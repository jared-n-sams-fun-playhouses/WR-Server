package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gopkg.in/ini.v1"
)

var (
	ip string
	port string
	root string
	upgrader = websocket.Upgrader {
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func main() {
	/* Config
	****************************************************************/
	cfg, err := ini.Load("./etc/cfg.ini")
	if err != nil {
		panic("\nConfig: " + err.Error())
		return
	}

	server := cfg.Section("Server")

	ip   = server.Key("ip").String()
	port = server.Key("port").String()
	root = server.Key("rootdir").String()

	fmt.Println("Server Root: " + root)

	/* Routes
	****************************************************************/
	wrs := mux.NewRouter()

	wrs.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		indexPage := &Page{}

		if err := indexPage.CreateView(root + "/test.html"); err != nil {
			fmt.Println("Index: " + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(indexPage.ServePage());
	}).
	Methods("GET")

	wrs.HandleFunc("/dj", func (w http.ResponseWriter, r *http.Request) {
		path := root + "/dj"

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
		path := root + "/audience"

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

			fmt.Printf("Received b:");

			for n := 0;n < len(msg);n++ {
				fmt.Printf("%d,",msg[n]);
			}

			fmt.Printf("\n");

			err = conn.WriteMessage(messageType, msg);
			if  err != nil {
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
	Handler(http.FileServer(http.Dir(root))).
	Methods("GET")

	/* Serve
	****************************************************************/
	http.Handle("/", &MiddleRouter{wrs})
	err = http.ListenAndServe(ip + ":" + port, nil)

	if err != nil {
		panic("Error: " + err.Error())
	}
}
