package router

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Fa7C0n/mdm-be/handlers/employees"
	"github.com/Fa7C0n/mdm-be/middleware"
)

var middlewareHandler = middleware.CreateStack(
	middleware.AllowCors,
)

func NewServer() *http.Server {
	host := os.Getenv("HOST_ADDRESS")
	port := os.Getenv("PORT")

	if host == "" {
		host = "localhost"
	}

	if port == "" {
		port = "8080"
		log.Printf("defaulting to port :%s", port)
	}

	return &http.Server{
		Addr:    fmt.Sprintf("%s:%s", host, port),
		Handler: middlewareHandler(newHandler()),
	}
}

func newHandler() http.Handler {
	employeesHandler := employees.NewHandler()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /employees", employeesHandler.ListEmployees)
	mux.HandleFunc("GET /employee/{id}", employeesHandler.GetEmployeeById)
	mux.HandleFunc("POST /employee", employeesHandler.CreateEmployee)
	mux.HandleFunc("PATCH /employee/{id}", employeesHandler.UpdateEmployee)
	mux.HandleFunc("DELETE /employee/{id}", employeesHandler.DeleteEmployee)

	return mux
}
