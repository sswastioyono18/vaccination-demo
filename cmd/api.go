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

func ResidentRoutes(residentService *services.ResidentService) *chi.Mux {
	router := chi.NewRouter()
	router.Get( "/resident/{nik}", residentService.GetUser)
	router.Post( "/resident/{nik}", residentService.Register)
	router.Post( "/vaccinate/{nik}", residentService.Vaccinate)
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
	messageQueueUri := fmt.Sprintf("amqp://%s:%s@%s:%d",  appConfig.MQ.User,  appConfig.MQ.Pass,  appConfig.MQ.Host,  appConfig.MQ.Port)
	residentExchangeRegistrationQueue, err  := infra.NewBrokerExchange(appConfig.MQ.Resident.Exchanges.ResidentVaccination, appConfig.MQ.Resident.Queues.Registration, messageQueueUri)
	if err != nil {
		log.Fatal("error during init mq", err)
	}
	defer residentExchangeRegistrationQueue.Channel.Close()

	residentExchangeVaccinationQueue, err  := infra.NewBrokerExchange(appConfig.MQ.Resident.Exchanges.ResidentVaccination, appConfig.MQ.Resident.Queues.Vaccination, messageQueueUri)
	if err != nil {
		log.Fatal("error during init mq", err)
	}
	defer residentExchangeVaccinationQueue.Channel.Close()

	newResidentService, err := services.NewResidentService(services.WithRabbitMQExchange(residentExchangeRegistrationQueue), services.WithRabbitMQExchange(residentExchangeVaccinationQueue))
	if err != nil {
		log.Fatal("error new resident service", err)
	}

	router := ResidentRoutes(newResidentService)
	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	zlogger.Info(fmt.Sprintf("[API] HTTP serve at %s\n", srv.Addr))

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		zlogger.Info(fmt.Sprintf("[API] Fail to start listen and server: %v", err))
	}
}
