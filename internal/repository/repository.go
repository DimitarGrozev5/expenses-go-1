package repository

type DatabaseRepo interface {
	Close() error
	AllUsers() bool
}
