package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"gopkg.in/ini.v1"
)

var (
	root string
	upgrader = websocket.Upgrader {
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

type MiddleRouter struct {
	mux *mux.Router
}

func (s *MiddleRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if origin := r.Header.Get("Origin"); origin != "" {
        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
        w.Header().Set("Access-Control-Allow-Headers",
            "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
    }
    // Stop here if its Preflighted OPTIONS request
    if r.Method == "OPTIONS" {
        return
    }

    w.Header().Set("Master", "JRod")

    s.mux.ServeHTTP(w, r) // Lets Gorilla work
}

func print_binary(s []byte) {
	fmt.Printf("Received b:");

	for n := 0;n < len(s);n++ {
		fmt.Printf("%d,",s[n]);
	}

	fmt.Printf("\n");
}

func main() {
	/* Config
	****************************************************************/
	cfg, err := ini.Load("./etc/cfg.ini")
	if err != nil {
		panic("\nConfig: " + err.Error())
	}

	server := cfg.Section("Server")

	var (
		ip   = server.Key("ip").String()
		port = server.Key("port").String()
		root = server.Key("rootdir").String()
	)

	fmt.Println("Server Root: " + root)

	/* Routes
	****************************************************************/
	wrs := mux.NewRouter()

	wrs.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		indexPage := &Page{}

		err := indexPage.CreateView(root + "/test.html")
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(indexPage.ServePage());
	}).
	Methods("GET")

	wrs.HandleFunc("/dj", func (w http.ResponseWriter, r *http.Request) {
		path := root + "/dj"

		indexPage := &Page{}
		err := indexPage.CreateView(path + "/index.html")
		if err != nil {
			fmt.Println(err)
			return
		}

		w.Write(indexPage.ServePage())
	}).
	Methods("GET")

	wrs.HandleFunc("/audience", func (w http.ResponseWriter, r *http.Request) {
		path := root + "/audience"

		indexPage := &Page{}
		err := indexPage.CreateView(path + "/index.html")
		if err != nil {
			fmt.Println(err)
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

			print_binary(msg)

			err = conn.WriteMessage(messageType, msg);
			if  err != nil {
				return
			}
		}

	})

	/* File Server
	****************************************************************/
	wrs.
	PathPrefix("/").
	Handler(http.FileServer(http.Dir(root))).
	Methods("GET")

	http.Handle("/", &MiddleRouter{wrs})
	err = http.ListenAndServe(ip + ":" + port, nil)

	if err != nil {
		panic("Error: " + err.Error())
	}
}
