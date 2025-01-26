package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

func SetupV1Routes() http.Handler {
	router := httprouter.New()
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	})
	SetupBaseRoutes(router)
	return corsConfig.Handler(router)
}

func Listen(listenAddress string) {
	router := SetupV1Routes()
	err := http.ListenAndServe(listenAddress, router)
	if err != nil {
		logrus.Fatal(err)
	}
}
