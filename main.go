package main

import (
	"fmt"
	"net/http"
	"io/ioutil"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
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
	root = "/home/jared/projects/go/src/github.com/jrods/wr_server/"
	fmt.Println("Server Root: ", root)

	wrs := mux.NewRouter()

	wrs.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		indexPage, err := ioutil.ReadFile(root + "test.html")
		if err != nil {
			fmt.Println(err)
			return 
		}

		w.Write(indexPage)
	})

	wrs.HandleFunc("/echo", func (w http.ResponseWriter, r *http.Request) {
		
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}
			
			print_binary(p)
			
			err = conn.WriteMessage(messageType, p);
			if  err != nil {
				return
			}
		}
		
	})

	wrs.
	PathPrefix("/").
	Handler(http.FileServer(http.Dir(root + "/"))).
	Methods("GET")

	http.Handle("/", &MiddleRouter{wrs})
	err := http.ListenAndServe(":10101", nil)
	
	if err != nil {
		panic("Error: " + err.Error())
	}
}
