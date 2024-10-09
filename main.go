package main

import (
	"ecommerce-yt/middleware"
	"ecommerce-yt/routes"
	"fmt"
	"net/http"
	"os"
)

func main() {
	router := http.NewServeMux()

	routes.RegisterRoutes(router)

	v1 := http.NewServeMux()
	v1.Handle("/v1/", http.StripPrefix("/v1", router)) // so that prefixes using v1 works

	stack := middleware.CreateStack(
		middleware.Logging,
	)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := http.Server{
		Addr:    ":" + port, // Use the field name Addr for the address
		Handler: stack(v1),  // Use the field name Handler for the router
	}

	fmt.Println("Server listening at port: " + port)

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
