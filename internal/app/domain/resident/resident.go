package resident

import (
	"time"
)

type Resident struct {
	NIK        string
	Birthplace string
	DoB        string
	FirstName  string
	LastName   string
}

type RegistrationRequest struct {
	NIK        string `json:"nik" valid:"required"`
	Birthplace string `json:"birth_place" valid:"required"`
	DoB        string `json:"birth_date" valid:"required"`
	FirstName  string `json:"first_name" valid:"required"`
	LastName   string `json:"last_name" valid:"required"`
}

type VaccinationRequest struct {
	NIK string
}

type VaccinatedInfo struct {
	NIK string
	Attempt int
	DateOfVaccinated time.Time
	Status bool
	Reason string

}
