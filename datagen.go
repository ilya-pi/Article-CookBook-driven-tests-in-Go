package main

import "math/rand"

type DataGen struct {
}

func (d *DataGen) PositiveNumber() int {
	return rand.Intn(100)
}

func (d *DataGen) Increment(limit int) int {
	return rand.Intn(limit)
}

func (d *DataGen) NegativeNumber() int {
	return -rand.Intn(100)
}
