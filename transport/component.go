package transport

import (
	"Timos-API/Vuement/persistence"
	"Timos-API/Vuement/service"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Timos-API/authenticator"
	"github.com/gorilla/mux"
)

type ComponentTransporter struct {
	s *service.ComponentService
}

var (
	ErrMissingId = "Missing param: id"
)

func NewComponentTransporter(s *service.ComponentService) *ComponentTransporter {
	return &ComponentTransporter{s}
}

func (c *ComponentTransporter) RegisterComponentRoutes(router *mux.Router) {
	router.HandleFunc("/vuement/component", c.getComponents).Methods("GET")
	router.HandleFunc("/vuement/component/{id}", c.getComponent).Methods("GET")
	router.HandleFunc("/vuement/component", authenticator.Middleware(c.createComponent, authenticator.Guard().G("admin").P("vuement.create"))).Methods("POST")
	router.HandleFunc("/vuement/component/{id}", authenticator.Middleware(c.updateComponent, authenticator.Guard().G("admin").P("vuement.update"))).Methods("PATCH")
	router.HandleFunc("/vuement/component/{id}", authenticator.Middleware(c.deleteComponent, authenticator.Guard().G("admin").P("vuement.delete"))).Methods("DELETE")

	fmt.Println("Component routes registered")
}

func (t *ComponentTransporter) getComponents(w http.ResponseWriter, req *http.Request) {
	components, err := t.s.GetAll(req.Context())

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(components)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *ComponentTransporter) getComponent(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]

	if !ok {
		http.Error(w, ErrMissingId, http.StatusBadRequest)
		return
	}

	component, err := t.s.GetById(req.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(component)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *ComponentTransporter) createComponent(w http.ResponseWriter, req *http.Request) {
	var body persistence.Component
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	component, err := t.s.Create(req.Context(), body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(component)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *ComponentTransporter) updateComponent(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]

	if !ok {
		http.Error(w, ErrMissingId, http.StatusBadRequest)
		return
	}

	var body persistence.Component
	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	component, err := t.s.Update(req.Context(), id, body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(component)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (t *ComponentTransporter) deleteComponent(w http.ResponseWriter, req *http.Request) {
	id, ok := mux.Vars(req)["id"]

	if !ok {
		http.Error(w, ErrMissingId, http.StatusBadRequest)
		return
	}

	success, err := t.s.Delete(req.Context(), id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !success {
		http.Error(w, "Couldn't delete message", http.StatusInternalServerError)
	}
}
