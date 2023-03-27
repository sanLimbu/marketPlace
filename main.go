package main

import (
	"crypto-exchange/server"
	"fmt"
	"log"
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

	bestBidPrice, err := c.GetBestBid()
	if err != nil {
		panic(err)
	}
	fmt.Println("bestBidPrice :", bestBidPrice)

	bestAskPrice, err := c.GetBestAsk()
	if err != nil {
		panic(err)
	}
	fmt.Println("bestAskPrice: ", bestAskPrice)

	// 		time.Sleep(1 * time.Second)
	// 	}
	// }()

	select {}
}

func marketOrderPlacer(c *client.Client) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		otherMarketSell := &client.PlaceOrderParams{
			UserID: 1,
			Bid:    false,
			Size:   5000,
		}
		orderResp, err := c.PlaceMarketOrder(otherMarketSell)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketSell := &client.PlaceOrderParams{
			UserID: 777,
			Bid:    false,
			Size:   2000,
		}
		orderResp, err = c.PlaceMarketOrder(marketSell)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		marketBuyOrder := &client.PlaceOrderParams{
			UserID: 777,
			Bid:    false,
			Size:   4000,
		}
		orderResp, err = c.PlaceMarketOrder(marketBuyOrder)
		if err != nil {
			log.Println(orderResp.OrderID)
		}

		<-ticker.C
	}
}
