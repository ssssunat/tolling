package main

import (
	"log"

	"github.com/ssssunat/tolling/aggregator/client"
)


const (
	kafkaTopic = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
)

// transport (http, grpc, kafka) --> attachh business Logic to this transport
func main() {
	var (
		err error
		svc CalculatorServicer
	)
	svc = NewCalculatorService()
	svc = NewLogMiddleware(svc)
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client.Newclient(aggregatorEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
