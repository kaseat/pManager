package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/kaseat/pManager/docs"

	"github.com/gorilla/mux"
	"github.com/kaseat/pManager/api"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Portfolio manager API
// @version 1.0
// @license.name MIT
// @license.url https://github.com/kaseat/pManager/blob/master/LICENSE
// @host totallink.ru
// @BasePath /api
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {

	fmt.Println("Started!")

	router := mux.NewRouter()

	router.PathPrefix("/api/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("https://totallink.ru/api/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("#swagger-ui"),
	))

	misc := router.PathPrefix("/api/misc").Subrouter().StrictSlash(true)
	misc.Use(api.VerifyTokenMiddleware)
	misc.HandleFunc("/validate", api.ValidateToken).Methods("GET")
	misc.HandleFunc("/gmail/url", api.AddGoogleAuth).Methods("GET")

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
	portfolios.HandleFunc("/{id}/sync", api.SyncOperations).Methods("GET")
	router.HandleFunc("/api/user/login", api.Login).Methods("POST")
	router.HandleFunc("/api/user/signup", api.SignUp).Methods("POST")
	router.HandleFunc("/api/google/callback", api.AppCallback).Methods("GET")
	log.Fatal(http.ListenAndServe(":8081", router))
}
