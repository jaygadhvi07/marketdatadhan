package process

import (
	"fmt"
	"marketdata/main/types"
	"time"
	// "os"

	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
)

type Order struct {
	Id int32              `json:"id"`
	InstrumentName string `json:"InstrumentName"`
	Type string           `json:"Type"`
	Quantity uint32       `json:"Quantity"`
	QuotePrice float32    `json:"QuotePrice"`
	Price float32         `json:"Price"`
	Stoploss float32      `json:"Stoploss"`
	Squareoff float32	  `json:"Squareoff"`
	Slippage float32      `json:"Slippage"`
	Sp float32 			  `json:"Sp"`
}


type Orderbook struct {
	OrderNo int64
	InstrumentName string
	LastTradedTime uint32
	Ltt string
	Type string	
	Price float32
	Quantity int32
	NoOfOrders int16
	Spread float32
	Ltp float32
}

type Levels1 struct {
	InstrumentName string
	LastTradedTime uint32
	BidQuantity int32
	NoOfBidOrders int16
	BidPrice float32
	AskQuantity int32
	NoOfAskOrders int16
	AskPrice float32
	Ltp float32
}

func insertTopOrderbook(database *sql.DB, topBid Orderbook, topAsk Orderbook, lastTradedPrice float32) error {
	// Convert Orderbook struct to JSON

	// processes(database, order)

	var levelOne Levels1
	if topBid.InstrumentName == topAsk.InstrumentName {
		levelOne.InstrumentName = topBid.InstrumentName
		levelOne.LastTradedTime = topBid.LastTradedTime
		
		levelOne.BidQuantity = topBid.Quantity
		levelOne.NoOfBidOrders = topBid.NoOfOrders
		levelOne.BidPrice = topBid.Price

		levelOne.AskQuantity = topAsk.Quantity
		levelOne.NoOfAskOrders = topAsk.NoOfOrders
		levelOne.AskPrice = topAsk.Price
		
		levelOne.Ltp = lastTradedPrice
	}

	// fmt.Println("TopLevel:", levelOne)

	processes(database, levelOne)
	
	var today string = generateDailyTableName()

	// Insert JSON into orderbook table
	topOrderbookJSON, err := json.Marshal(levelOne)
	
	if err != nil {
		return fmt.Errorf("error marshaling top orderbook: %v", err)
	}


	var topOrderbookTableName = fmt.Sprintf("orderbook_top_%s", today) // concat with today
	var query string = fmt.Sprintf("INSERT INTO %s (data) VALUES (?)", topOrderbookTableName)
    
	_, err = database.Exec(query, topOrderbookJSON)

	if err != nil {
		return fmt.Errorf("error inserting into top orderbook table: %v", err)
	}

	return nil
}


func generateDailyTableName() string {
    now := time.Now() // Type: time.Time
    return now.Format("2006_01_02") // Type: string
}


func insertOrderbook(database *sql.DB, orderbook Orderbook) error {
	// Convert Orderbook struct to JSON
	orderbookJSON, err := json.Marshal(orderbook)
	if err != nil {
		return fmt.Errorf("error marshaling orderbook: %v", err)
	}

	var today string = generateDailyTableName()

	// Insert JSON into orderbook table
	var orderbookTableName = fmt.Sprintf("orderbook_%s", today) // concat with today
	var query string = fmt.Sprintf("INSERT INTO %s (data) VALUES (?)", orderbookTableName)
    
	_, err = database.Exec(query, orderbookJSON)

	if err != nil {
		return fmt.Errorf("error inserting into orderbook table: %v", err)
	}

	return nil
}

func insertMarketbook(database *sql.DB, fullPacket types.FullPacket) error {
	// Set Levels5 field to nil before serializing

	// Convert FullPacket struct to JSON
	fullPacketJSON, err := json.Marshal(fullPacket)
	if err != nil {
		return fmt.Errorf("error marshaling fullPacket: %v", err)
	}

	
	var today string = generateDailyTableName()

	var marketbookTableName = fmt.Sprintf("marketbook_%s", today) // concat with today
	var query string = fmt.Sprintf("INSERT INTO %s (data) VALUES (?)", marketbookTableName)
    

	// Insert JSON into marketbook table
	_, err = database.Exec(query, fullPacketJSON)
	if err != nil {
		return fmt.Errorf("error inserting into marketbook table: %v", err)
	}

	return nil
}


func Process(fullDataFeed types.FullPacket, database *sql.DB) {

	// fmt.Println("FullPacketEntry:", fullDataFeed)

	var bids []Orderbook
	var asks []Orderbook
	var lastTradedPrice float32 = fullDataFeed.LastTradedPrice

	var err error

	err = insertMarketbook(database, fullDataFeed)
	if err != nil {
		fmt.Println("Error while saving:", err)
	}
	
	var i int
	for i = 0; i < len(fullDataFeed.Levels5); i++ {	
		
		var bid Orderbook
		var ask Orderbook

		bid.InstrumentName = fullDataFeed.InstrumentName
		bid.LastTradedTime = fullDataFeed.LastTradedTime
		bid.Type = "Bid"
		bid.Quantity = fullDataFeed.Levels5[i].BidQuantity 
		bid.Price = fullDataFeed.Levels5[i].BidPrice
		bid.NoOfOrders = fullDataFeed.Levels5[i].NoOfBidOrders 

		bids = append(bids, bid)

		ask.InstrumentName = fullDataFeed.InstrumentName
		ask.LastTradedTime = fullDataFeed.LastTradedTime
		ask.Type = "Ask"
		ask.Quantity = fullDataFeed.Levels5[i].AskQuantity 
		ask.Price = fullDataFeed.Levels5[i].AskPrice
		ask.NoOfOrders = fullDataFeed.Levels5[i].NoOfAskOrders  

		asks = append(asks, ask)
	}

	/*var order Levels1
	order.InstrumentName = bids[0].InstrumentName
	order.LastTradedTime = bids[0].LastTradedTime
	order.BidQuantity = bids[0].BidQuantity
	order.NoOfBidOrders = bids[0].NoOfBidOrders 
	order.BidPrice = bids[0].BidPrice
	order.AskQuantity = asks[0].AskQuantity
	order.NoOfAskOrders = asks[0].NoOfAskOrders
	order.AskPrice = asks[0].AskPrice
	order.Ltp = lastTradedPrice*/


	err =  insertTopOrderbook(database, bids[0], asks[0], lastTradedPrice)	
	if err != nil {
		fmt.Println("Error while saving", err)
	}
}

