package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/sswastioyono18/vaccination-demo/config"
	"github.com/sswastioyono18/vaccination-demo/internal/app/infra"
	zlog "github.com/sswastioyono18/vaccination-demo/internal/app/middleware"
	"github.com/sswastioyono18/vaccination-demo/internal/app/services"
	"log"
	"net/http"
	"time"
)

// spaHandler implements the http.Handler interface, so we can use it
// to respond to HTTP requests. The path to the static directory and
// path to the index file within that static directory are used to
// serve the SPA in the given static directory.
type spaHandler struct {
	staticPath string
	indexPath  string
}

func MustParams(h http.Handler, params string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		h.ServeHTTP(w, r) // all params present, proceed
	}
}

func ResidentRoutes(residentService *services.ResidentService) *chi.Mux {
	router := chi.NewRouter()
	router.Get( "/check/resident/{nik}", residentService.CheckUser)

	router.Post( "/vaccine/{nik}", MustParams(http.HandlerFunc(residentService.RegisterHandler), "ster"))
	router.Put( "/vaccine/{nik}", residentService.Update)
	router.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	return router
}

func main() {
	zlog.NewLogger("PROD")
	zlogger := zlog.Logger
	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal("error during config init")
	}

	//TODO need to be moved later to separate block of init
	//TODO defer connection
	messageQueueUri := fmt.Sprintf("amqp://%s:%s@%s:%d",  appConfig.MQ.User,  appConfig.MQ.Pass,  appConfig.MQ.Host,  appConfig.MQ.Port)
	residentExchange,err  := infra.NewBrokerExchange(appConfig.MQ.Exchanges.ResidentVaccination, appConfig.MQ.Queues.NewVaccineRegistration, messageQueueUri)
	if err != nil {
		log.Fatal("error during init mq", err)
	}
	defer residentExchange.Channel.Close()

	newResidentService, err := services.NewResidentService(services.WithRabbitMQExchange(residentExchange))
	if err != nil {
		log.Fatal("error new resident service", err)
	}

	router := ResidentRoutes(newResidentService)
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	zlogger.Info(fmt.Sprintf("[API] HTTP serve at %s\n", srv.Addr))

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		zlogger.Info(fmt.Sprintf("[API] Fail to start listen and server: %v", err))
	}
}
