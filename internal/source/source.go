package source

import "log-proj/pkg/models"

type Interface interface {
	GetCarsWithFilter(filter interface{}) (models.Cars, error)
}
