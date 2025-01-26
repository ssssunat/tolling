package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/ssssunat/tolling/types"
)


func main() {
	recv, err := NewDataReceiver()
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgchan chan types.OBUData
	conn    *websocket.Conn
	prod    DataProducer
}

func NewDataReceiver() (*DataReceiver, error) {
	var (
		p DataProducer
		err error
		kafkTopic = "obudata"
	)
	p, err = NewKafkaProducer(kafkTopic)
	if err != nil {
		return nil, err
	}

	p = NewLogMiddleware(p)
	return &DataReceiver{
		msgchan: make(chan types.OBUData, 128),
		prod:    p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU connected client connected !")

	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read error:", err)
			continue
		}
		fmt.Printf("received OBU data from [%d]  :: <lat %.2f, long %.2f>\n", data.OBUID, data.Lat, data.Long)
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka produce error:", err)
		}
	}

}
