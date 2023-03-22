package orderbook

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Error("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	l := NewLimit(10000)
	buyOrderA := NewOrder(true, 15)
	buyOrderB := NewOrder(true, 15)
	buyOrderC := NewOrder(true, 50)
	l.AddOrder(buyOrderA)
	l.AddOrder(buyOrderB)
	l.AddOrder(buyOrderC)
	l.DeleteOrder(buyOrderB)
	fmt.Println(l)
}

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderBook()
	sellOrderA := NewOrder(true, 10)
	sellOrderB := NewOrder(true, 15)
	ob.PlaceLimitOrder(10000, sellOrderA)
	ob.PlaceLimitOrder(15000, sellOrderB)

	assert(t, len(ob.asks), 2)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderBook()

	sellOrder := NewOrder(false, 20)
	ob.PlaceLimitOrder(10000, sellOrder)

	buyOrder := NewOrder(true, 100)
	matches := ob.PlaceMarketOrder(buyOrder)

	fmt.Printf("%+v", matches)
}
