package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/streadway/amqp"
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

type ResidentRegisterRequest struct {
	NIK        string `json:"nik" valid:"required"`
	Birthplace string `json:"birth_place" valid:"required"`
	DoB        string `json:"birth_date" valid:"required"`
	FirstName  string `json:"first_name" valid:"required"`
	LastName   string `json:"last_name" valid:"required"`
}

type Resident struct {
	NIK        string
	Birthplace string
	DoB        string
	FirstName  string
	LastName   string
}

type ResidentRegisterVaccination struct {
	ResidentData Resident
	VaccinationData UserVaccinated
}

type UserVaccinated struct {
	Attempt int
	DateOfVaccinated time.Time
	Status bool
	Reason string

}

var amqConnection *amqp.Connection
var channelRabbitMQ *amqp.Channel

// check user
func CheckUser(w http.ResponseWriter, r *http.Request) {
	nik := chi.URLParam(r, "nik")
	w.Write([]byte(fmt.Sprintf("NIK %s Exists!\n", nik)))
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var resident ResidentRegisterRequest
	err := json.NewDecoder(r.Body).Decode(&resident)
	if err != nil {
		http.Error(w, "Invalid User Data", 500)
	}

	residentData := &ResidentRegisterVaccination{
		ResidentData: Resident(resident),
		VaccinationData: UserVaccinated{
			Attempt:          1,
			DateOfVaccinated: time.Now(),
			Status:           true,
			Reason:           "",
		},
	}

	byteResult , err := json.Marshal(residentData)

	// Create a message to publish.
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        byteResult,
	}

	// Attempt to publish a message to the queue.
	err = channelRabbitMQ.Publish(
		"VaccinationExchange",              // exchange
		"NewResidentVaccination", // queue name
		false,           // mandatory
		false,           // immediate
		message,         // message to publish
	)

	if err != nil {
		log.Fatal("error publish")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resident)
	w.Write(j)
}

func Update(w http.ResponseWriter, r *http.Request) {
	var userVaccinated UserVaccinated
	err := json.NewDecoder(r.Body).Decode(&userVaccinated)
	if err != nil {
		http.Error(w, "Invalid Order Data", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(userVaccinated)
	w.Write(j)
}

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Get( "/check/resident/{nik}", CheckUser)
	router.Post( "/vaccine/{nik}", RegisterHandler)
	router.Put( "/vaccine/{nik}", Update)
	router.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	return router
}

func main() {

	// Define RabbitMQ server URL.
	amqpServerURL := "amqp://guest:guest@localhost:5672"

	// Create a new RabbitMQ connection.
	var err error
	amqConnection, err = amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer amqConnection.Close()

	// Let's start by opening a channel to our RabbitMQ
	// instance over the connection we have already
	// established.
	channelRabbitMQ, err = amqConnection.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()

	err = channelRabbitMQ.ExchangeDeclare(
		"VaccinationExchange",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		panic(err)
	}

	router := Routes()
	fs := http.FileServer(http.Dir("static"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))
	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println( fmt.Sprintf("[API] HTTP serve at %s\n", srv.Addr))

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Println(fmt.Sprintf("[API] Fail to start listen and server: %v", err))
	}
}
