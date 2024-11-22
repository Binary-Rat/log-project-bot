package array

import "log-proj/pkg/models"

type db struct {
	Cars models.Cars
}

func New() *db {
	return &db{
		Cars: cars,
	}
}

// Need to think about neccessity of Cars struct couse d.Cars.Cars is terriable
func (d *db) GetCars(loadV, loadW float64) (cars models.Cars, err error) {
	for _, car := range d.Cars.Cars {
		if cars.Len() > 2 {
			break
		}
		if car.LoadV >= loadV && car.LoadW >= loadW {
			cars.Add(car)
		}
	}
	return cars, nil
}
