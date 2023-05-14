package main

// import (
// 	"crypto-exchange/client"
// 	"crypto-exchange/marketmaker"
// 	"crypto-exchange/server"
// 	"math/rand"
// 	"time"
// )

import (
	k "crypto-exchange/kafka"
	w "crypto-exchange/websocket"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"golang.org/x/net/websocket"
)

func main() {

	//set up kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "producer-1",
	})

	if err != nil {
		log.Fatalf("Failed to create kafka producer: %s\n", err)
	}

	//Set up kafka consumer

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"group.id":          "consumer-group-1",
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %s\n", err)
	}

	//Create a wbesocket server
	server := w.NewServer()

	//Create an OrderPlacer for kafka message production
	orderPlacer := k.NewOrderPlacer(p, "exchange")

	// Handle WebSocket connections and Kafka messages
	go server.HandleConnections()
	go k.ConsumeMessages(c, server.BroadcastMessages())

	// Set up HTTP routes
	http.Handle("/ws", websocket.Handler(server.HandleWebSocket))

	http.HandleFunc("/placeorder", func(w http.ResponseWriter, r *http.Request) {
		orderType := r.FormValue("orderType")
		size := r.FormValue("size")
		numSize, _ := strconv.Atoi(size)

		// Place order using Kafka producer
		err := orderPlacer.PlaceOrder(orderType, numSize)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to place order: %s", err.Error()), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	fs := http.FileServer(http.Dir("web/static"))
	http.Handle("/web/static/", http.StripPrefix("/web/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	// http.Handle("/orderbookfeed", websocket.Handler(func(ws *websocket.Conn) {
	// 	// Start a goroutine to handle outgoing messages
	// 	go handleMessages(ws, server.Messages)

	// }))

	// Start the server
	log.Println("Server listening on :3000")
	log.Fatal(http.ListenAndServe(":3000", nil))

}

func handleMessages(ws *websocket.Conn, messages <-chan string) {
	for message := range messages {
		err := websocket.Message.Send(ws, message)
		if err != nil {
			log.Printf("Error sending message to WebSocket client: %v", err)
		}
	}
}

// func main() {
// 	go server.StartServer()
// 	time.Sleep(1 * time.Second)
// 	c := client.NewClient()

// 	cfg := marketmaker.Config{
// 		UserID:         0,
// 		OrderSize:      10,
// 		MinSpread:      20,
// 		MakeInterval:   1 * time.Second,
// 		SeedOffset:     40,
// 		ExchangeClient: c,
// 		PriceOffset:    10,
// 	}

// 	maker := marketmaker.NewMarketMaker(cfg)
// 	maker.Start()

// 	time.Sleep(2 * time.Second)

// 	go marketOrderPlacer(c)

// 	select {}
// }

// func marketOrderPlacer(c *client.Client) {
// 	ticker := time.NewTicker(500 * time.Millisecond)

// 	for {
// 		randint := rand.Intn(10)
// 		bid := true
// 		if randint < 7 {
// 			bid = false
// 		}

// 		order := client.PlaceOrderParams{
// 			UserID: 0,
// 			Bid:    bid,
// 			Size:   1,
// 		}

// 		_, err := c.PlaceMarketOrder(&order)
// 		if err != nil {
// 			panic(err)
// 		}

// 		<-ticker.C
// 	}
// }
