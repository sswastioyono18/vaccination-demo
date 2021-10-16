package services

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	residentDomain "github.com/sswastioyono18/vaccination-demo/internal/app/domain/resident"
	"github.com/sswastioyono18/vaccination-demo/internal/app/infra"
	"github.com/sswastioyono18/vaccination-demo/internal/app/middleware"
	"net/http"
)

var logger =  middleware.Logger

// OrderConfiguration is an alias for a function that will take in a pointer to an ResidentService and modify it
type ResidentConfiguration func(os *ResidentService) error

// ResidentService is a implementation of the ResidentService
type ResidentService struct {
	residents residentDomain.ResidentRepository
	queue     infra.MessageBroker
}

func NewResidentService(cfgs ...ResidentConfiguration) (rs *ResidentService, err error) {
	rs = &ResidentService{}
	// Apply all Configurations passed in
	for _, cfg := range cfgs {
		// Pass the service into the configuration function
		err = cfg(rs)
		if err != nil {
			return nil, err
		}
	}

	return
}

func WithResidentRepository(cr residentDomain.ResidentRepository) ResidentConfiguration {
	return func(os *ResidentService) error {
		os.residents = cr
		return nil
	}
}

func WithRabbitMQExchange(rq infra.MessageBroker) ResidentConfiguration {
	return func(os *ResidentService) error {
		os.queue = rq
		return nil
	}
}


// check user
func (rs *ResidentService) GetUser(w http.ResponseWriter, r *http.Request) {
	nik := chi.URLParam(r, "nik")
	w.Write([]byte(fmt.Sprintf("NIK %s Exists!\n", nik)))
}

func (rs *ResidentService) Register(w http.ResponseWriter, r *http.Request) {


	var resident residentDomain.RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&resident)
	if err != nil {
		http.Error(w, "Invalid Resident Data", 500)
	}

	residentData := &residentDomain.Resident{
		NIK:        resident.NIK,
		Birthplace: resident.Birthplace,
		DoB:        resident.DoB,
		FirstName:  resident.FirstName,
		LastName:   resident.LastName,
	}


	residentByte, err := json.Marshal(residentData)
	if err != nil {
		return
	}

	err = rs.queue.Publish("resident_registration", residentByte)
	if err != nil {
		logger.Error("Failed to publish")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(resident)
	w.Write(j)
}

// this function is to send nik user to vaccination queue. Consumer of vaccination queue will determine whether to update or insert
func (rs *ResidentService) Vaccinate(w http.ResponseWriter, r *http.Request) {
	var resident residentDomain.VaccinationRequest
	err := json.NewDecoder(r.Body).Decode(&resident)
	if err != nil {
		http.Error(w, "Invalid Resident Data", 500)
	}

	vaccinationByte, err := json.Marshal(resident.NIK)
	if err != nil {
		return
	}

	err = rs.queue.Publish("resident_vaccination", vaccinationByte)
	if err != nil {
		return 
	}

}

func (rs *ResidentService) Update(w http.ResponseWriter, r *http.Request) {
	var userVaccinated residentDomain.VaccinatedInfo
	err := json.NewDecoder(r.Body).Decode(&userVaccinated)
	if err != nil {
		http.Error(w, "Invalid Resident Data", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	j, _ := json.Marshal(userVaccinated)
	w.Write(j)
}