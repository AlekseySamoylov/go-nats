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
	log.Println("Dialer started")

	time.Sleep(5 * time.Second)
	natsConnection, _ := nats.Connect(nats.DefaultURL)
	deliveryCount := uint32(0)
	subscribeForDelivery(natsConnection, &deliveryCount)
	for true {
		orderManyCars(natsConnection)
		time.Sleep(2 * time.Second)
		log.Println("Number of delivered Mustang Shelby GT350:")
		log.Println(deliveryCount)
	}

	defer natsConnection.Close()
}

func orderManyCars(natsConnection *nats.Conn) {
	order := new(car.Order)
	order.Id = uuid.New().String()
	order.Amount = 500_000
	orderData, _ := proto.Marshal(order)
	msg, err := natsConnection.Request("order.service", orderData, 1000*time.Millisecond)
	if err != nil {
		log.Fatal("Cannot request")
	}
	orderAccepted := car.OrderAccepted{}
	_ = proto.Unmarshal(msg.Data, &orderAccepted)
	if orderAccepted.OrderId == order.Id {
		log.Println("Order accepted")
	}
}

func subscribeForDelivery(natsConnection *nats.Conn, counter *uint32) {
	_, subError := natsConnection.Subscribe("delivery.service", func(m *nats.Msg) {
		carDelivery := car.Delivery{}
		_ = proto.Unmarshal(m.Data, &carDelivery)
		if carDelivery.Model == "Ford Mustang Shelby GT350" {
			atomic.AddUint32(counter, 1)
		}
	})
	if subError != nil {
		log.Fatal("Cannot subscribe")
	}
}
