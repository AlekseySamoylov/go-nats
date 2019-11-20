package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/nats-io/go-nats"
	"go-nats/car"
	"log"
	"sync/atomic"
	"time"
)

func main() {
	log.Println("Dealer started")
	time.Sleep(5 * time.Second)

	natsConnection, _ := nats.Connect(nats.DefaultURL)
	subscribeForDelivery(natsConnection)
	carNumberCorrection := 0
	for true {
		orderManyCars(natsConnection, carNumberCorrection)
		carNumberCorrection = waitDelivery()
	}
	defer natsConnection.Close()
}

func orderManyCars(natsConnection *nats.Conn, carNumberCorrection int) {
	order := new(car.Order)
	order.Id = uuid.New().String()
	//log.Printf("Car number correction %d \n", carNumberCorrection)
	order.Amount = int32(carAmount + carNumberCorrection)
	//log.Printf("Car number correction final %d \n", order.Amount)
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

func subscribeForDelivery(natsConnection *nats.Conn) {
	_, subError := natsConnection.Subscribe(deliverySubject, func(m *nats.Msg) {
		carDelivery := car.Delivery{}
		_ = proto.Unmarshal(m.Data, &carDelivery)
		if carDelivery.Model == "Ford Mustang Shelby GT350" {
			atomic.AddUint32(&deliverySum, 1)
		}
	})
	if subError != nil {
		log.Fatal("Cannot subscribe")
	}
}

func waitDelivery() int {
	previousDeliveryAmount := uint32(0)
	for true {
		if (time.Now().After(deliveryTimeout) || deliverySum%carAmount == 0) && previousDeliveryAmount == deliverySum {
			//log.Printf("Previous amount and current %d , %d \n", previousDeliveryAmount, deliveryCount)
			break
		}
		previousDeliveryAmount = deliverySum
		time.Sleep(500 * time.Microsecond)
	}

	log.Printf("Number of delivered Mustang Shelby GT350: %d \n", deliverySum)
	deliveryTimeout = time.Now().Add(time.Second * 5)
	deliveredForOrderOrZero := deliverySum % carAmount
	if deliveredForOrderOrZero == 0 {
		return 0
	} else {
		return int(carAmount - deliveredForOrderOrZero)
	}
}

var deliveryTimeout = time.Now().Add(time.Second * 5)
var deliverySum = uint32(0)

const deliverySubject = "delivery.service"
const carAmount = 500_000
