package middleware

import "github.com/go-gorp/gorp"

// CustomerRepository is a interface that defines the rules around what a customer repository
// Has to be able to perform
type ResidentRepository interface {
	Get(id uint64)
}

type ResidentRepo struct {
	db gorp.DbMap
}

func (r ResidentRepo) Get(id uint64) {
	panic("implement me")
}

func NewResidentRepository(dbMap gorp.DbMap) *ResidentRepo {
	return &ResidentRepo{
		db: dbMap,
	}
}
