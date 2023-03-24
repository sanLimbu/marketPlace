package orderbook

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10000)
	buyOrderA := NewOrder(true, 35, 0)
	buyOrderB := NewOrder(true, 15, 0)
	buyOrderC := NewOrder(true, 50, 0)
	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)
	l.DeleteOrder(buyOrderB)
	fmt.Println(l)
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderBook()

	sellOrderA := NewOrder(false, 10, 0)
	sellOrderB := NewOrder(false, 15, 0)
	ob.PlaceLimitOrder(10000, sellOrderA)
	ob.PlaceLimitOrder(9000, sellOrderB)

	assert(t, len(ob.Orders), 2)
	assert(t, ob.Orders[sellOrderA.ID], sellOrderA)
	assert(t, ob.Orders[sellOrderB.ID], sellOrderB)
	assert(t, len(ob.asks), 2)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderBook()

	sellOrder := NewOrder(false, 200, 0)
	ob.PlaceLimitOrder(10000, sellOrder)

	buyOrder := NewOrder(true, 100, 0)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskTotalVolume(), 100.0)
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 100.0)
	assert(t, matches[0].Price, 10000.0)
	assert(t, buyOrder.IsFilled(), true)

	fmt.Printf("%+v", matches)
}

func TestPlaceMarketOrderMultiFill(t *testing.T) {
	ob := NewOrderBook()

	buyOrderA := NewOrder(true, 10, 0)
	buyOrderB := NewOrder(true, 20, 0)
	buyOrderC := NewOrder(true, 30, 0)
	buyOrderD := NewOrder(true, 10, 0)

	ob.PlaceLimitOrder(8000, buyOrderC)
	ob.PlaceLimitOrder(8000, buyOrderD)
	ob.PlaceLimitOrder(9000, buyOrderB)
	ob.PlaceLimitOrder(10000, buyOrderA)

	sellOrder := NewOrder(false, 40, 0)
	matches := ob.PlaceMarketOrder(sellOrder)
	assert(t, ob.BidTotalVolume(), 30.0)
	assert(t, len(matches), 3)
	assert(t, len(ob.bids), 1)
	fmt.Printf("%+v", matches)
}

func TestCancelOrder(t *testing.T) {
	ob := NewOrderBook()
	buyOrder := NewOrder(true, 4)
	ob.PlaceLimitOrder(10000.0, buyOrder)

	assert(t, ob.BidTotalVolume(), 4.0)

	ob.CancelOrder(buyOrder)

	assert(t, ob.BidTotalVolume(), 0.0)

	_, ok := ob.Orders[buyOrder.ID]
	assert(t, ok, false)
}
