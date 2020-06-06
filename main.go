package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/kaseat/pManager/docs"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/api"
	"github.com/kaseat/pManager/portfolio"
	httpSwagger "github.com/swaggo/http-swagger"
)

//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6ImFkbWluIiwiZXhwIjoxNTkxMDI5MTgyfQ.QdY08oVd-e-qb1sLlCpuq0GcdjN1UEy4vPcaO9DZe7A
// @title Portfolio manager API
// @version 1.0
// @license.name MIT
// @license.url https://github.com/kaseat/pManager/blob/master/LICENSE
// @host localhost
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	cfg := portfolio.Config{
		MongoURL: "mongodb://localhost:27017",
		DbName:   "tcs2",
	}
	portfolio.Init(cfg)

	fmt.Println("Started!")

	router := mux.NewRouter()

	router.PathPrefix("/api/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("https://totallink.ru/api/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	tokens := router.PathPrefix("/api/token").Subrouter().StrictSlash(true)
	tokens.Use(api.VerifyTokenMiddleware)
	tokens.HandleFunc("/validate", api.ValidateToken).Methods("GET")
	portfolios := router.PathPrefix("/api/portfolios").Subrouter().StrictSlash(true)
	portfolios.Use(api.VerifyTokenMiddleware)
	portfolios.HandleFunc("", api.CreateSinglePortfolio).Methods("POST")
	portfolios.HandleFunc("", api.ReadAllPortfolios).Methods("GET")
	portfolios.HandleFunc("", api.DeleteAllPortfolios).Methods("DELETE")
	portfolios.HandleFunc("/{id}", api.ReadSinglePortfolio).Methods("GET")
	portfolios.HandleFunc("/{id}", api.UptateSinglePortfolio).Methods("PUT")
	portfolios.HandleFunc("/{id}", api.DeleteSinglePortfolio).Methods("DELETE")
	portfolios.HandleFunc("/{id}/operations", api.ReadOperations).Methods("GET")
	portfolios.HandleFunc("/{id}/operations", api.CreateSingleOperation).Methods("POST")
	portfolios.HandleFunc("/{id}/operations", api.DeleteAllOperations).Methods("DELETE")
	portfolios.HandleFunc("/{id}/average", api.GetAveragePrice).Methods("GET")
	portfolios.HandleFunc("/{id}/balance", api.GetBalance).Methods("GET")
	router.HandleFunc("/api/auth/login", api.GetToken).Methods("POST")
	router.HandleFunc("/api/auth/signup", api.CreateOwner).Methods("POST")
	log.Fatal(http.ListenAndServe(":8081", router))
}
