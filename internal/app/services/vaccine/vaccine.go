package vaccine

import (
	"encoding/json"
	residentDomain "github.com/sswastioyono18/vaccination-demo/internal/app/domain/resident"
	"github.com/sswastioyono18/vaccination-demo/internal/app/infra"
	"github.com/sswastioyono18/vaccination-demo/internal/app/middleware"
	"net/http"
)

var logger =  middleware.Logger

// OrderConfiguration is an alias for a function that will take in a pointer to an VaccineService and modify it
type VaccineConfiguration func(os *VaccineService) error

// VaccineService is a implementation of the VaccineService
type VaccineService struct {
	residents residentDomain.ResidentRepository
	queue     infra.MessageBroker
}

func NewVaccineService(cfgs ...VaccineConfiguration) (rs *VaccineService, err error) {
	rs = &VaccineService{}
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

func WithRabbitMQExchange(rq infra.MessageBroker) VaccineConfiguration {
	return func(os *VaccineService) error {
		os.queue = rq
		return nil
	}
}

// this function is to send nik user to vaccination queue. Consumer of vaccination queue will determine whether to update or insert
func (rs *VaccineService) Vaccinate(w http.ResponseWriter, r *http.Request) {
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

func (rs *VaccineService) Update(w http.ResponseWriter, r *http.Request) {
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