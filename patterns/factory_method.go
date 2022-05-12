package main

import "fmt"

type ITransport interface {
	setName(string)
	getName() string
	setType(string)
	getType() string
	Deliver()
}

type Transport struct {
	Name string
	Type string
}

func (t *Transport) setName(name string) {
	t.Name = name
}

func (t Transport) getName() string {
	return t.Name
}

func (t *Transport) setType(transportType string) {
	t.Type = transportType
}

func (t Transport) getType() string {
	return t.Type
}

func (t Transport) Deliver() {
	fmt.Println("deliver from", t.Name, t.Type)
}

type Truck struct {
	Transport
	FuelConsumption float64
}

func NewTruck(name string, transportType string) ITransport {
	return &Truck{
		Transport: Transport{
			Name: name,
			Type: transportType,
		},
		FuelConsumption: 12.00,
	}
}

type Ship struct {
	Transport
	Capacity float64
}

func NewShip(name string, transportType string) ITransport {
	return &Ship{
		Capacity: 250,
		Transport: Transport{
			Name: name,
			Type: transportType,
		},
	}
}

func main() {
	t := NewTruck("man", "land")
	fmt.Println(t.getName())
	t.Deliver()

	s := NewShip("ship", "shipping")
	fmt.Println(s.getType())
	s.Deliver()
}
