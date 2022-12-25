package db

var (
	Repo Repository
)

type Repository interface {
	Todo() error
}