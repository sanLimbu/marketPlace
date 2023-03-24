package server

import (
	"context"
	"crypto-exchange/orderbook"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/labstack/echo/v4"
)

const (
	exchangePrivateKey           = "e485d098507f54e7733a205420dfddbe58db035fa577fc294ebd14db90767a52"
	MarketOrder        OrderType = "MARKET"
	LimitOrder         OrderType = "LIMIT"
	MarketETH          Market    = "ETH"
)

type (
	OrderType string
	Market    string
)

type Exchange struct {
	Client *ethclient.Client
	Users  map[int64]*User
	Orders map[int64][]*orderbook.Order

	PrivateKey *ecdsa.PrivateKey
	orderbooks map[Market]*orderbook.OrderBook
}

type PlaceOrderRequest struct {
	UserID int64
	Type   OrderType
	Bid    bool
	Size   float64
	Price  float64
	Market Market
}

type Order struct {
	UserID    int64
	ID        int64
	Price     float64
	Size      float64
	Bid       bool
	Timestamp int64
}

type OrderBookData struct {
	TotalBidVolume float64
	TotalAskVolume float64
	Asks           []*Order
	Bids           []*Order
}

type MatchedOrder struct {
	UserID int64
	Price  float64
	Size   float64
	ID     int64
}

func StartServer() {
	e := echo.New()

	client, err := ethclient.Dial("http://localhost:8545")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("we have a connection")
	ex, err := NewExchange(exchangePrivateKey, client)
	if err != nil {
		log.Fatal(err)
	}

	buyeraddress := "0xFFcf8FDEE72ac11b5c542428B35EEF5769C409f0"
	buyerbalance, err := client.BalanceAt(context.Background(), common.HexToAddress(buyeraddress), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(buyerbalance)

	selleraddress := "0x22d491Bde2303f2f43325b2108D26f1eAbA1e32b"
	sellerbalance, err := client.BalanceAt(context.Background(), common.HexToAddress(selleraddress), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(sellerbalance)

	johnaddress := "0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1"
	johnbalance, err := client.BalanceAt(context.Background(), common.HexToAddress(johnaddress), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(johnbalance)

	pkStr1 := "6cbed15c793ce57650b9877cf6fa156fbef513c4e6134f022a85b1ffdd59b2a1"
	user1 := NewUser(pkStr1, 1)
	ex.Users[user1.ID] = user1

	pkStr2 := "6370fd033278c143179d81c5526140625662b8daa446c22ee2d73db3707e620c"
	user2 := NewUser(pkStr2, 2)
	ex.Users[user2.ID] = user2

	johnPk := "4f3edf983ac636a65a842ce7c78d9aa706d3b113bce9c46f30d7d21715b23b1d"
	john := NewUser(johnPk, 0)
	ex.Users[john.ID] = john

	e.HTTPErrorHandler = httpErrorHandler

	e.GET("/book/:market", ex.handleGetBook)
	e.POST("/order", ex.handlePlaceHolder)
	e.DELETE("/order/:id", ex.cancelOrder)
	e.GET("/book/:market/bid", ex.handleGetBestBid)
	e.GET("/book/:market/ask", ex.handleGetBestAsk)

	e.Start(":3000")
}

type User struct {
	ID         int64
	PrivateKey *ecdsa.PrivateKey
}

func NewUser(privKey string, id int64) *User {
	pk, err := crypto.HexToECDSA(privKey)
	if err != nil {
		panic(err)
	}

	return &User{
		ID:         id,
		PrivateKey: pk,
	}
}

func httpErrorHandler(err error, c echo.Context) {
	fmt.Println(err)
}

func NewExchange(exchangePrivateKey string, client *ethclient.Client) (*Exchange, error) {
	orderbooks := make(map[Market]*orderbook.OrderBook)
	orderbooks[MarketETH] = orderbook.NewOrderBook()

	privateKey, err := crypto.HexToECDSA(exchangePrivateKey)
	if err != nil {
		log.Fatal(err)
	}

	return &Exchange{
		Client:     client,
		Users:      make(map[int64]*User),
		Orders:     make(map[int64][]*orderbook.Order),
		PrivateKey: privateKey,
		orderbooks: orderbooks,
	}, nil
}

func (ex *Exchange) handleGetBook(c echo.Context) error {
	market := Market(c.Param("market"))
	ob, ok := ex.orderbooks[market]
	if !ok {
		return c.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found"})
	}

	orderBookData := OrderBookData{
		TotalBidVolume: ob.BidTotalVolume(),
		TotalAskVolume: ob.AskTotalVolume(),
		Asks:           []*Order{},
		Bids:           []*Order{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			o := Order{
				UserID:    order.UserID,
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderBookData.Asks = append(orderBookData.Asks, &o)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			o := Order{
				UserID:    order.UserID,
				ID:        order.ID,
				Price:     limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderBookData.Bids = append(orderBookData.Bids, &o)
		}
	}
	return c.JSON(http.StatusOK, orderBookData)

}

type PriceResponse struct {
	Price float64
}

func (ex *Exchange) handleGetBestBid(c echo.Context) error {

	market := Market(c.Param("market"))
	ob := ex.orderbooks[market]

	if len(ob.Bids()) == 0 {
		return fmt.Errorf("the bids are empty")
	}

	bestBidPrice := ob.Bids()[0].Price
	pr := PriceResponse{
		Price: bestBidPrice,
	}
	return c.JSON(http.StatusOK, pr)

}

func (ex *Exchange) handleGetBestAsk(c echo.Context) error {

	market := Market(c.Param("market"))
	ob := ex.orderbooks[market]

	if len(ob.Asks()) == 0 {
		return fmt.Errorf("the asks are empty")
	}

	bestAsksPrice := ob.Asks()[0].Price
	pr := PriceResponse{
		Price: bestAsksPrice,
	}
	return c.JSON(http.StatusOK, pr)

}

func (ex *Exchange) cancelOrder(c echo.Context) error {

	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)

	ob := ex.orderbooks[MarketETH]
	order := ob.Orders[int64(id)]
	ob.CancelOrder(order)

	log.Println("order canceled id => ", id)

	return c.JSON(200, map[string]any{"msg": "order deleted"})
}

type PlaceOrderResponse struct {
	OrderID int64
}

func (ex *Exchange) handlePlaceMarketOrder(market Market, order *orderbook.Order) ([]orderbook.Match, []*MatchedOrder) {

	ob := ex.orderbooks[market]
	matches := ob.PlaceMarketOrder(order)
	matchedOrders := make([]*MatchedOrder, len(matches))

	isBid := false
	if order.Bid {
		isBid = true
	}

	for i := 0; i < len(matchedOrders); i++ {
		id := matches[i].Bid.ID
		if isBid {
			id = matches[i].Ask.ID
		}

		matchedOrders[i] = &MatchedOrder{
			ID:    id,
			Size:  matches[i].SizeFilled,
			Price: matches[i].Price,
		}
	}
	return matches, matchedOrders
}

func (ex *Exchange) handlePlaceLimitOrder(market Market, price float64, order *orderbook.Order) error {
	ob := ex.orderbooks[market]
	ob.PlaceLimitOrder(price, order)
	log.Printf("new LIMIT order => type: [%t] | price: [%.2f] | size [%.2f]", order.Bid, order.Limit.Price, order.Size)
	return nil
}

func (ex *Exchange) handlePlaceHolder(c echo.Context) error {

	var placeOrderData PlaceOrderRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&placeOrderData); err != nil {
		return err
	}

	market := Market(placeOrderData.Market)
	order := orderbook.NewOrder(placeOrderData.Bid, placeOrderData.Size, placeOrderData.UserID)

	if placeOrderData.Type == LimitOrder {
		if err := ex.handlePlaceLimitOrder(market, placeOrderData.Price, order); err != nil {
			return err
		}
	}
	if placeOrderData.Type == MarketOrder {
		matches, _ := ex.handlePlaceMarketOrder(market, order)
		if err := ex.handleMatches(matches); err != nil {
			return err
		}

	}

	resp := &PlaceOrderResponse{
		OrderID: order.ID,
	}
	return c.JSON(200, resp)
}

func (ex *Exchange) handleMatches(matches []orderbook.Match) error {
	for _, match := range matches {
		fromUser, ok := ex.Users[match.Ask.UserID]
		if !ok {
			return fmt.Errorf("user not found: %d", match.Ask.UserID)
		}

		toUser, ok := ex.Users[match.Bid.UserID]
		if !ok {
			return fmt.Errorf("user not found: %d", match.Bid.UserID)
		}

		toAddress := crypto.PubkeyToAddress(toUser.PrivateKey.PublicKey)

		// this is only used for the fees
		// exchangePubKey := ex.PrivateKey.Public()
		// publicKeyECDSA, ok := exchangePubKey.(*ecdsa.PublicKey)
		// if !ok {
		// 	return fmt.Errorf("error casting public key to ECDSA")
		// }
		amount := big.NewInt(int64(match.SizeFilled))
		transferETH(ex.Client, fromUser.PrivateKey, toAddress, amount)

	}
	return nil
}

func transferETH(client *ethclient.Client, fromPrivKey *ecdsa.PrivateKey, to common.Address, amount *big.Int) error {
	ctx := context.Background()
	publicKey := fromPrivKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return err
	}

	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal(err)
	}

	tx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)

	chainID := big.NewInt(1337)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), fromPrivKey)
	if err != nil {
		return err
	}

	return client.SendTransaction(ctx, signedTx)
}
