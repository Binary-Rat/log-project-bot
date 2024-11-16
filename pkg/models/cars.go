package models

type Car struct {
	Name  string
	LoadV float64
	LoadW float64
}

type Cars struct {
	Cars []Car
}

// Add appends a car to the slice of cars.
func (c *Cars) Add(car Car) {
	c.Cars = append(c.Cars, car)
}

// Len returns the length of the slice of cars.
func (c *Cars) Len() int {
	return len(c.Cars)
}

// Names returns a slice of car names.
func (c *Cars) Names() []string {
	var names []string
	for _, car := range c.Cars {
		names = append(names, car.Name)
	}
	return names
}
