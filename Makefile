obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./data_receiver
	@./bin/receiver

calculator:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator 

agg:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator 

gate:
	@go build -o bin/gate gateway/main.go 
	@./bin/gate

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

.PHONY: obu agg proto gate