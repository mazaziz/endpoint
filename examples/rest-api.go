package main

import (
	"fmt"
	"net/http"
	"endpoint"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, "Hello World")
}

func noteList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "note list")
}

func noteCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "note create")
}

func noteRetrieve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "note retrieve")
}

func noteDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "note delete")
}

func underConstruction(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(503)
	fmt.Fprintln(w, "Under Construction. Please try again later.")
}

func main() {
	e := endpoint.New()
	e.Handler404 = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		fmt.Fprint(w, "{\"success\":false,\"code\":404,\"message\":\"Not Found\"}")
	}
	e.Handler405 = func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(405)
		fmt.Fprint(w, "{\"success\":false,\"code\":405,\"message\":\"Method Not Allowed\"}")
	}
	e.Match("id", `^\d+$`)
	e.Match("tag", `^[[:alnum:]]+$`)
	e.Route("GET /", home)
	e.Route("GET /note", noteList)
	e.Route("POST /note", noteCreate)
	e.Route("GET /note/:id", noteRetrieve)
	e.Route("DELETE /note/:id", noteDelete)
	e.Route("GET /note/:id/tag/:tag", underConstruction)
	e.Serve("127.0.0.1:3333")
}
