package resident

// CustomerRepository is a interface that defines the rules around what a customer repository
// Has to be able to perform
type ResidentRepository interface {
	Get(id uint64)
}