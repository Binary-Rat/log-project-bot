package array

import "log-proj/pkg/models"

type db struct {
	Cars []models.Car
}

func New() *db {
	return &db{
		Cars: make([]models.Car, 0),
	}
}

func (d *db) GetCars(loadV, loadW float64) (cars []models.Car, err error) {
	for _, car := range d.Cars {
		if car.LoadV >= loadV && car.LoadW >= loadW {
			cars = append(cars, car)
		}
	}
	return cars, nil
}
