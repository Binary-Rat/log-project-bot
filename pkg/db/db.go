package db

import "log-proj/pkg/models"

type Interface interface {
	GetCars(loadV, loadW float64) ([]models.Car, error)
}
