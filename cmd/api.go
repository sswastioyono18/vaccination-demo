package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-gorp/gorp"
	"github.com/sswastioyono18/vaccination-demo/config"
	"github.com/sswastioyono18/vaccination-demo/internal/app/infra"
	mdlware "github.com/sswastioyono18/vaccination-demo/internal/app/middleware"
	"github.com/sswastioyono18/vaccination-demo/internal/app/services/resident"
	"github.com/sswastioyono18/vaccination-demo/internal/app/services/vaccine"
	"go.uber.org/zap"
	"log"
	"net/http"
	"time"
)

func residentRoutes(router *chi.Mux, residentService *resident.ResidentService) *chi.Mux {
	router.Get( "/resident/{nik}", residentService.GetUser)
	router.Post( "/resident/{nik}", residentService.Register)
	return router
}

func vaccineRoutes(router *chi.Mux, vaccineService *vaccine.VaccineService) *chi.Mux {
	router.Post( "/vaccinate/{nik}", vaccineService.Vaccinate)
	return router
}

func healthRoutes(router *chi.Mux) *chi.Mux {
	router.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	return router
}

func main() {
	mdlware.NewLogger("PROD")
	zlogger := mdlware.Logger
	appConfig, err := config.NewConfig()
	if err != nil {
		log.Fatal("error during config init")
	}

	//TODO need to be moved later to separate block of init
	messageQueueUri := fmt.Sprintf("amqp://%s:%s@%s",  appConfig.MQ.User,  appConfig.MQ.Pass,  appConfig.MQ.Host)
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

	var dbPostgre *gorp.DbMap
	dbPostgre, err = infra.NewPostgreDatabase(appConfig)
	if err != nil {
		return
	}

	newResidentRepo := mdlware.NewResidentRepository(*dbPostgre)
	newResidentService, err := resident.NewResidentService(resident.WithRabbitMQExchange(residentExchangeRegistrationQueue), resident.WithResidentRepository(newResidentRepo))
	if err != nil {
		zlogger.Fatal("error new resident service", zap.Error(err))
	}

	newVaccineService, err := vaccine.NewVaccineService(vaccine.WithRabbitMQExchange(residentExchangeVaccinationQueue))
	if err != nil {
		zlogger.Fatal("error new vaccine service", zap.Error(err))
	}

	router := chi.NewRouter()
	router = residentRoutes(router, newResidentService)
	router = vaccineRoutes(router, newVaccineService)
	router = healthRoutes(router)

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
