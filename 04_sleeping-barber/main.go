package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fatih/color"
)

var seatingCapacity = 10
var arravialRate = 100
var cutDuration = 1000 * time.Millisecond
var timeOpen = 10 * time.Second

func main() {
	rand.Seed(time.Now().UnixNano())
	// print out welcome message
	color.Yellow("The Sleeping Barber Problem")
	color.Yellow("---------------------------")

	// create channels
	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	// create the barbershop
	shop := BarberShop{
		ShopyCapacity:   seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		ClientsChan:     clientChan,
		BarbersDoneChan: doneChan,
		Open:            true,
	}

	// add barbers
	color.Green("The shop is opened for the day!")
	shop.addBarber("Frank")
	shop.addBarber("John")
	shop.addBarber("Bob")
	shop.addBarber("Grag")
	shop.addBarber("Corey")

	// start the barbershop as gorutine
	shopClosing := make(chan bool)
	closed := make(chan bool)
	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		shop.closeShopForDay()
		closed <-true
	}()

	// add clients
	i:=1
	
	go func ()  {
		for{
			// get a random number with average arravial range
			randomMilliseconds := rand.Int() % (2 * arravialRate)
			select {
			case <-shopClosing:
				return
			case <-time.After(time.Millisecond * time.Duration(randomMilliseconds)):
				shop.addClient(fmt.Sprintf("Client #%d", i))
				i++
			}
		}
	}()
	// block till the barbershop is closed
	<-closed
}
