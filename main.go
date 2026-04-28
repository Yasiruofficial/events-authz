package main

import (
	"events-authz/configs"
	"events-authz/internal/cache"
	eventsHttp "events-authz/internal/http"
	"events-authz/internal/service"
	"events-authz/internal/spicedb"
	"log"
)

func main() {
	cfg := configs.Load()

	spice := spicedb.NewClient(cfg.SpiceAddr)
	c := cache.New()

	svc := service.NewAuthzService(spice, c)
	handler := eventsHttp.NewHandler(svc)

	router := eventsHttp.NewRouter(handler)

	log.Println("AuthZ service running on :8080")
	log.Fatal(router.Run(":8080"))
}
