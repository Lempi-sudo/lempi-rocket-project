package service

//go:generate ../../.././bin/mockery --case=underscore --all

type PaymentService interface {
	Pay() (string, error)
}
