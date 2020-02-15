package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/api"
	"github.com/kaseat/pManager/portfolio"
)

func main() {
	cfg := portfolio.Config{
		MongoURL: "mongodb://localhost:27017",
		DbName:   "tcs2",
	}
	portfolio.Init(cfg)

	fmt.Println("Started!")

	router := mux.NewRouter()
	portfolios := router.PathPrefix("/portfolios").Subrouter().StrictSlash(true)
	portfolios.Use(api.VerifyTokenMiddleware)
	portfolios.HandleFunc("/", api.CreateSinglePortfolio).Methods("POST")
	portfolios.HandleFunc("/", api.ReadAllPortfolios).Methods("GET")
	portfolios.HandleFunc("/", api.DeleteAllPortfolios).Methods("DELETE")
	portfolios.HandleFunc("/{id}", api.ReadSinglePortfolio).Methods("GET")
	portfolios.HandleFunc("/{id}", api.UptateSinglePortfolio).Methods("PUT")
	portfolios.HandleFunc("/{id}", api.DeleteSinglePortfolio).Methods("DELETE")
	portfolios.HandleFunc("/{id}/operations", api.ReadOperations).Methods("GET")
	portfolios.HandleFunc("/{id}/operations", api.CreateSingleOperation).Methods("POST")
	portfolios.HandleFunc("/{id}/operations", api.DeleteAllOperations).Methods("DELETE")
	portfolios.HandleFunc("/{id}/average", api.GetAveragePrice).Methods("GET")
	portfolios.HandleFunc("/{id}/balance", api.GetBalance).Methods("GET")
	router.HandleFunc("/auth/getjwt", api.GetToken).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", router))
}
