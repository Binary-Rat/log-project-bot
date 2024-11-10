package ati

import (
	"fmt"
	"log-proj/pkg/models"

	"github.com/Binary-Rat/atisu"
)

type Atisu struct {
	page         int
	itemsPerPage int
	client       *atisu.Client
}

func New(token string, isDemo bool) (*Atisu, error) {
	client, err := atisu.NewClient(token, isDemo)
	if err != nil {
		return nil, err
	}
	return &Atisu{
		client: client,
	}, nil
}

// In progress
func (a *Atisu) GetCarsWithFilter(filter interface{}) (cars models.Cars, err error) {
	// atisuFilter, ok := filter.(atisu.Filter)
	// if !ok {
	// 	return nil, errors.New("wrong filter format")
	// }
	// body, err := a.client.GetCarsWithFilter(a.page, a.itemsPerPage, atisuFilter)
	// if err != nil {
	// 	return nil, e.Warp("can`t get cars with filter", err)
	// }

	// err = json.Unmarshal(body, &cars)
	// if err != nil {
	// 	return nil, e.Warp("can`t unmarshal cars", err)
	// }
	return cars, nil
}

func (a *Atisu) GetCityID(cities []string) (citiesWithID *atisu.Cities, err error) {
	citiesWithID, err = a.client.GetCityID(cities)
	if err != nil {
		return nil, err
	}
	for name, city := range *citiesWithID {
		if !city.IsSuccess {
			return nil, fmt.Errorf("city %s not found", name)
		}
	}

	return citiesWithID, nil
}
