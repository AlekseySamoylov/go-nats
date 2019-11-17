package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/nats-io/go-nats"
	"go-nats/car"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	log.Println("Dealer started")
	time.Sleep(5 * time.Second)

	natsConnection, _ := nats.Connect(nats.DefaultURL)
	deliveryCount := uint32(0)
	waitGroup.Add(0)
	subscribeForDelivery(natsConnection, &deliveryCount)

	for true {
		waitDelivery(&deliveryCount)
		waitGroup.Add(carAmount)
		time.Sleep(5 * time.Second)
		orderManyCars(natsConnection)
	}
	defer natsConnection.Close()
}

func orderManyCars(natsConnection *nats.Conn) {
	order := new(car.Order)
	order.Id = uuid.New().String()
	order.Amount = carAmount
	order.Subject = deliverySubject
	orderData, _ := proto.Marshal(order)
	msg, err := natsConnection.Request("order.service", orderData, 1000*time.Millisecond)
	if err != nil {
		log.Println("Cannot request")
	} else {
		orderAccepted := car.OrderAccepted{}
		_ = proto.Unmarshal(msg.Data, &orderAccepted)
		if orderAccepted.OrderId == order.Id {
			log.Println("Order accepted")
		}
	}
}

func subscribeForDelivery(natsConnection *nats.Conn, counter *uint32) {
	_, subError := natsConnection.Subscribe(deliverySubject, func(m *nats.Msg) {
		carDelivery := car.Delivery{}
		_ = proto.Unmarshal(m.Data, &carDelivery)
		if carDelivery.Model == "Ford Mustang Shelby GT350" {
			atomic.AddUint32(counter, 1)
			countDown()
		}
	})
	if subError != nil {
		log.Fatal("Cannot subscribe")
	}
}

func countDown() {
	waitGroup.Done()
}

func waitDelivery(deliveryCount *uint32) {
	if waitTimeout(&waitGroup, waitDeliverySeconds*time.Second) {
		log.Println("Timed out waiting for car delivery")
	} else {
		log.Println("Delivery finished successfully")
	}
	log.Printf("Number of delivered Mustang Shelby GT350: %d \n", *deliveryCount)
}

// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}

var waitGroup = sync.WaitGroup{}

const deliverySubject = "delivery.service"
const carAmount = 500_000
const waitDeliverySeconds = 5
