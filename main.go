package main

import (
	"crypto-exchange/server"
	"fmt"
	"time"

	"crypto-exchange/client"
)

func main() {
	go server.StartServer()
	time.Sleep(1 * time.Second)
	c := client.NewClient()

	// go func() {
	// 	for {

	limitOrderParams := &client.PlaceOrderParams{
		UserID: 1,
		Bid:    false,
		Price:  10000,
		Size:   5000,
	}

	resp, err := c.PlaceLimitOrder(limitOrderParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("place limit order_id => ", resp.OrderID)

	otherLimitOrderParams := &client.PlaceOrderParams{
		UserID: 0,
		Bid:    false,
		Price:  9000,
		Size:   500,
	}

	resp, err = c.PlaceLimitOrder(otherLimitOrderParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("OTHER place limit order_id => ", resp.OrderID)

	buyLimitOrderParams := &client.PlaceOrderParams{
		UserID: 0,
		Bid:    true,
		Price:  11000,
		Size:   500,
	}

	resp, err = c.PlaceLimitOrder(buyLimitOrderParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("buy limit place limit order_id => ", resp.OrderID)

	marketOrderParams := &client.PlaceOrderParams{
		UserID: 2,
		Bid:    true,
		Size:   1000,
	}
	resp, err = c.PlaceMarketOrder(marketOrderParams)
	if err != nil {
		panic(err)
	}
	fmt.Println("placed market order_id => ", resp.OrderID)

	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	select {}
}
