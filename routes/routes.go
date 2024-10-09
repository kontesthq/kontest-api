package routes

import (
	"fmt"
	"kontest-api/controllers"
	"net/http"
)

func HelloGETHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! GET")
}

func HelloPOSTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! POST")
}

func HelloPUTHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! PUT")
}

func HelloDELETEHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World! DELETE")
}

func RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("POST /users/register", controllers.SignUp)
	router.HandleFunc("POST /users/login", controllers.SignIn)
	router.HandleFunc("GET /kontests", controllers.GetAllKontests)
	router.HandleFunc("DELETE /purge", controllers.PurgeMetadata)

	registerHelloRoutes(router)
}

func registerHelloRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /hello", HelloGETHandler)
	router.HandleFunc("POST /hello", HelloPOSTHandler)
	router.HandleFunc("DELETE /hello", HelloDELETEHandler)
	router.HandleFunc("PUT /hello", HelloPUTHandler)
}
