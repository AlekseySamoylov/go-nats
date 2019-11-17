package main

import (
	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
	"go-nats/car"
	"log"
	"runtime"
	"time"
)

func main() {
	log.Println("Factory started")
	time.Sleep(3 * time.Second)

	natsConnection1, _ := nats.Connect(nats.DefaultURL)
	natsConnection2, _ := nats.Connect(nats.DefaultURL)
	go startProductionForOrderSubject(natsConnection1, "order.jvm.service")
	go startProductionForOrderSubject(natsConnection2, "order.service")

	// Keep the connection alive
	runtime.Goexit()
	defer natsConnection1.Close()
	defer natsConnection2.Close()
}

func startProductionForOrderSubject(natsConnection *nats.Conn, orderSubject string) {
	_, _ = natsConnection.Subscribe(orderSubject, func(m *nats.Msg) {
		carOrder := car.Order{}
		log.Println("Order recieved")
		_ = proto.Unmarshal(m.Data, &carOrder)
		acceptOrder(natsConnection, m.Reply, carOrder.Id)

		for i := int32(0); i < carOrder.Amount; i++ {
			carDelivery := assembleTheCar(carOrder.Id)
			deliverCar(natsConnection, carDelivery, carOrder.Subject)
		}
		log.Println("All cars sent to delivery")
	})
}

func acceptOrder(natsConnection *nats.Conn, replySubject string, orderId string) {
	orderAccepted := car.OrderAccepted{}
	orderAccepted.OrderId = orderId
	orderAcceptedData, _ := proto.Marshal(&orderAccepted)
	_ = natsConnection.Publish(replySubject, orderAcceptedData)
	log.Println("Order accept published")
}
func assembleTheCar(carOrderId string) car.Delivery {
	carDelivery := car.Delivery{}
	carDelivery.OrderId = carOrderId
	carDelivery.Model = "Ford Mustang Shelby GT350"
	carDelivery.Details = "5.2L Ti-VCT V8"
	return carDelivery
}

func deliverCar(natsConnection *nats.Conn, carDelivery car.Delivery, deliverySubject string) {
	data, err := proto.Marshal(&carDelivery)
	if err == nil {
		_ = natsConnection.Publish(deliverySubject, data)
	}
}
