package source

import (
	"log-proj/pkg/models"

	"github.com/Binary-Rat/atisu"
)

type Interface interface {
	GetCarsWithFilter(filter interface{}) (models.Cars, error)
	GetCityID(cities []string) (citiesWithID *atisu.Cities, err error)
}
