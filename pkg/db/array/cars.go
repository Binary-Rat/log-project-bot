package array

import "log-proj/pkg/models"

var (
	cars = models.Cars{
		Cars: []models.Car{
			{Name: "Car1", LoadV: 100, LoadW: 100},
			{Name: "Car2", LoadV: 200, LoadW: 200},
			{Name: "Car3", LoadV: 300, LoadW: 300},
		},
	}
)
