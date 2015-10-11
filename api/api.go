package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/matobet/verdi/api/dto"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, "Hello World")
		},
	},
	Route{
		"VmsIndex",
		"GET",
		"/vms",
		func(w http.ResponseWriter, r *http.Request) {
			vms := dto.Vms{
				dto.Vm{ID: "4832480324-sufdo48304", Name: "my-vm-1"},
				dto.Vm{ID: "1232328309-kfjds32903", Name: "my-vm-2"},
			}

			if err := json.NewEncoder(w).Encode(vms); err != nil {
				log.Println(err)
			}
		},
	},
	Route{
		"VmCreate",
		"POST",
		"/vms",
		func(w http.ResponseWriter, r *http.Request) {
			vm := dto.Vm{}
			if err := json.NewDecoder(r.Body).Decode(&vm); err != nil {
				w.WriteHeader(422)
				log.Println(err)
			}

			// backend.CreateVM(vm)
		},
	},
	Route{
		"VmShow",
		"GET",
		"/vms/{vmId}",
		func(w http.ResponseWriter, r *http.Request) {
			//json.Marshal(v)
		},
	},
}

func Serve() error {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	log.Println("Listening on port 8080")
	return http.ListenAndServe(":8080", router)
}
