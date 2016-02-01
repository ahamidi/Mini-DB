package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

type appContext struct {
	db *MemDB
}

var context appContext

func retrieveValueAndLockHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	rec, err := context.db.getRecordAndLock(key)
	if err != nil {
		JSONResponse(w, 404, nil, err)
		return
	}
	JSONResponse(w, 200, rec, nil)
}

func updateKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]
	lockID := mux.Vars(r)["lock_id"]
	var release bool
	if r := r.FormValue("release"); r == "true" {
		release = true
	} else {
		release = false
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		JSONResponse(w, 500, nil, err)
		return
	}
	r.Body.Close()

	err = context.db.updateKey(key, lockID, string(body), release)
	if err != nil {
		if err.Error() == "Key Not Found" {
			JSONResponse(w, 404, nil, err)
		} else {
			JSONResponse(w, 401, nil, err)
		}
		return
	}

	JSONResponse(w, 204, nil, nil)

}

func createKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := mux.Vars(r)["key"]

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		JSONResponse(w, 500, nil, err)
		return
	}
	r.Body.Close()

	rec, err := context.db.createKey(key, string(body))
	if err != nil {
		JSONResponse(w, 404, nil, err)
		return
	}

	// According to specs, should only return `lock_id`, so remove value
	rec.Value = nil
	JSONResponse(w, 200, rec, nil)
}

func handlers() *mux.Router {
	// Mux setup
	router := mux.NewRouter()

	// API Routes
	router.HandleFunc("/reservations/{key}", retrieveValueAndLockHandler).Methods("POST")
	router.HandleFunc("/values/{key}/{lock_id}", updateKeyHandler).Methods("POST")
	router.HandleFunc("/values/{key}", createKeyHandler).Methods("PUT")

	return router
}

func main() {
	log.Println("MiniDB")

	// Init MemDB
	db := newMemDB()
	context.db = db

	// Handlers
	router := handlers()

	// Middleware
	n := negroni.Classic()
	n.UseHandler(router)

	// Server
	if os.Getenv("PORT") != "" {
		n.Run(strings.Join([]string{":", os.Getenv("PORT")}, ""))
	} else {
		n.Run(":8080")
	}

}
