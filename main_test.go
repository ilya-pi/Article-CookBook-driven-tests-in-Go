package main

import (
	"fmt"
	"testing"
)

func boot(addr string) (DataGen, CookBook, func()) {
	dg := DataGen{}
	cb := CookBook{addr}
	cleanup := func() {
		fmt.Println("Cleanup")
	}
	return dg, cb, cleanup
}

func TestPositiveNumbers(t *testing.T) {
	dg, cb, cleanup := boot("http://localhost:3000")
	defer cleanup()

	number := dg.PositiveNumber()
	increment := dg.Increment(number)
	total := number + increment

	Run(t, []TestStruct{
		cb.Ac("Clear all"),
		cb.Add("Add some positive number", number),
		cb.Total("Check the total amount", &number),
		cb.Add("Add an incerement", increment),
		cb.Total("Check total amount again", &total),
		cb.Ac("Check we can clear after all operations"),
	})
}

func TestNegativeNumbers(t *testing.T) {
	dg, cb, cleanup := boot("http://localhost:3000")
	defer cleanup()

	number := dg.NegativeNumber()
	number1 := dg.NegativeNumber()
	total := number + number1

	Run(t, []TestStruct{
		cb.Ac("Clear all"),
		cb.Add("Add a negative number", number),
		cb.Total("Check the total amount", &number),
		cb.Add("Add another negative number", number1),
		cb.Total("Check total negative again", &total),
		cb.Ac("Check we can clear after all operations"),
	})
}

func TestRand(t *testing.T) {
	_, cb, cleanup := boot("http://localhost:3000")
	defer cleanup()

	var r int

	Run(t, []TestStruct{
		cb.Rand("Assign a random number", &r),
		cb.Total("Check total is persisted", &r),
	})
}
