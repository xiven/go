package main

import (
	"fmt"
)

type contactInfo struct {
	email   string
	zipCode int
}

type person struct {
	firstName string
	lastName  string
	contactInfo
}

func main() {
	var jonathan person
	jonathan.firstName = "Jonathan"
	jonathan.lastName = "Osterman"
	jonathan.contactInfo.email = "j.m.osterman@gmail.com"
	jonathan.contactInfo.zipCode = 30076

	jonathan.updateName("johnny")
	jonathan.print()
}

func (pointerToPerson *person) updateName(newFirstName string) {
	(*pointerToPerson).firstName = newFirstName
}

func (p person) print() {
	fmt.Printf("%+v", p)
}
