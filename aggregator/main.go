package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ssssunat/tolling/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpAddr", ":4000", "the listen address of HTTP server")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "the listen address of GRPC server")
	flag.Parse()
	var (
		store = NewMemoryStore()
		svc   = NewInvoiceAggregator(store)
	)
	svc = NewMetricsMiddleware(svc)
	svc = NewLogMiddleware(svc)
	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenAddr, svc))
	}()
	makeHTTPTransport(*httpListenAddr, svc)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("GRPC transport runnig on port", listenAddr)
	// make  a TCP listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	// make a GRPC native server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register OUR grpc  server implementation to th GRPC package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport runnig on port", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(listenAddr, nil))
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]
		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obu ID"})
			return
		}
		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}
		inv, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, inv)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": "method not supported",
			})
			return		
		}
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{
				"error": err.Error(),
			})
			return
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
