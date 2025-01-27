package main

import (
	"log"
)

//	type DistanceCalculator struct {
//		consumer DataConsumer
//	}
const kafkaTopic = "obudata"

// transport (http, grpc, kafka) --> attachh business Logic to this transport
func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
