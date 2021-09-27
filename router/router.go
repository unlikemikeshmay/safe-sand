package router

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/mike/inv/middleware"
)

func Router(addr string, err error) *mux.Router {

	r := mux.NewRouter().StrictSlash(true)

	//router.HandleFunc("/", middleware.HomePage).Methods("GET")
	//middleware.IsAuthorized takes an endpoint and returns a http handler if jwt is valid
	r.Handle("/get-tools", middleware.IsAuthorized(middleware.GetItems)).Methods("GET")
	r.HandleFunc("/login", middleware.UserLogin).Methods("POST")
	r.Handle("/validate", middleware.IsAuthorized(middleware.Validate)).Methods("GET", "OPTIONS")

	r.Handle("/inventory", middleware.IsAuthorized(middleware.CreateItem)).Methods("POST")
	r.Handle("/inventory/{uid}", middleware.IsAuthorized(middleware.DeleteItem)).Methods("DELETE")
	r.Handle("/inventory", middleware.IsAuthorized(middleware.UpdateItem)).Methods("PUT")
	// user route handlers
	r.HandleFunc("/users", middleware.CreateUser).Methods("POST")
	r.Handle("/users/{userid}", middleware.IsAuthorized(middleware.DeleteUser)).Methods("DELETE")
	r.Handle("/users", middleware.IsAuthorized(middleware.GetUsers)).Methods("GET")
	r.Handle("/users", middleware.IsAuthorized(middleware.UpdateUser)).Methods("PUT")
	// show route handlers
	r.Handle("/shows/{uid}",middleware.IsAuthorized(middleware.DeleteShow)).Methods("DELETE")
	r.Handle("/shows",middleware.IsAuthorized(middleware.AddShow)).Methods("POST")
	r.Handle("/shows",middleware.IsAuthorized(middleware.GetShows)).Methods("GET")
	r.Handle("/updateShow",middleware.IsAuthorized(middleware.UpdateShow)).Methods("PUT")
	if err != nil {
		log.Fatal(err)
	} else {
		// CORS THAT IS WORKING?
		origins := handlers.AllowedOrigins([]string{"*"})
		headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Accept", "Authorization"})
		methods := handlers.AllowedMethods([]string{"POST", "GET", "OPTIONS", "PUT", "DELETE"})
		log.Fatal(http.ListenAndServe(addr, handlers.CORS(headers, methods, origins)(r)))
	}
	return r
}
