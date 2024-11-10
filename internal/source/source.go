package source

import "log-proj/pkg/models"

type Interface interface {
	GetCars(filter interface{}) (models.Cars, error)
}
