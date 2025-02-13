package source

import (
	"log-proj/pkg/models"

	"github.com/Binary-Rat/atisu"
)

type Interface interface {
	GetCarsWithFilter(filter atisu.Filter) (models.Cars, error)
	GetCityID(cities []string) (citiesWithID *atisu.Cities, err error)
}
