package ati

import (
	"encoding/json"
	"fmt"
	"log-proj/pkg/lib/e"
	"log-proj/pkg/models"

	"github.com/Binary-Rat/atisu"
)

type Atisu struct {
	page         int
	itemsPerPage int
	client       *atisu.Client
}

// Add ping
func New(token string, isDemo bool) (*Atisu, error) {
	client, err := atisu.NewClient(token, isDemo)
	if err != nil {
		return nil, err
	}
	return &Atisu{
		page:         1,
		itemsPerPage: 10,
		client:       client,
	}, nil
}

// In progress
func (a *Atisu) GetCarsWithFilter(filter atisu.Filter) (cars models.Cars, err error) {
	body, err := a.client.GetCarsWithFilter(a.page, a.itemsPerPage, filter)
	if err != nil {
		return cars, e.Warp("can`t get cars with filter", err)
	}

	err = json.Unmarshal(body, &cars)
	if err != nil {
		return cars, e.Warp("can`t unmarshal cars", err)
	}
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
