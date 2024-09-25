package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Handle the root route "/"
    http.HandleFunc("/", homeHandler)
    // Handle the "/hello" route
    http.HandleFunc("/hello", helloHandler)

    fmt.Println("Server is listening on port 8080...")
    // Start the server on port 8080
    if err := http.ListenAndServe(":8080", nil); err != nil {
        fmt.Println("Failed to start server:", err)
    }
}

// Handler for the root route "/"
func homeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the Home Page!")
}

// Handler for the "/hello" route
func helloHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, World!")
}

