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

type Orderbook struct {
	InstrumentName string
	Type string	
	Price float32
	Quantity int32
	NoOfOrders int16
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
	fullPacket.Levels5 = nil

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
		bid.Type = "Bid"
		bid.Quantity = fullDataFeed.Levels5[i].BidQuantity 
		bid.Price = fullDataFeed.Levels5[i].BidPrice
		bid.NoOfOrders = fullDataFeed.Levels5[i].NoOfBidOrders 

		bids = append(bids, bid)

		err = insertOrderbook(database, bid)
		if err != nil {
			fmt.Println("Error while saving", err)
		}


		ask.InstrumentName = fullDataFeed.InstrumentName
		ask.Type = "Ask"
		ask.Quantity = fullDataFeed.Levels5[i].AskQuantity 
		ask.Price = fullDataFeed.Levels5[i].AskPrice
		ask.NoOfOrders = fullDataFeed.Levels5[i].NoOfAskOrders  

		err = insertOrderbook(database, ask)
		if err != nil {
			fmt.Println("Error while saving", err)
		}

		
		asks = append(asks, ask)
	}
	
	// fmt.Println("ASK:", asks)

	sequances(bids, asks)
}

type InstrumentOrders struct {
	bids []Orderbook
	asks []Orderbook
}

var sequance map[string]InstrumentOrders

func sequances(bids []Orderbook, asks []Orderbook) {
	
	var tempask []Orderbook

	// Falling
	var i int = 0
	for i < len(asks) {
		var j int = i + 1
		var countI int = 0
		for j <= len(asks) {
			countI += 1
			if j == len(asks) {
				if asks[i].Price >= tempask[len(tempask) - 1].Price {
					tempask = append(tempask, asks[i])
					break
				}
			}
			if asks[j].Price >= asks[i].Price {
				if countI == 1 {
					tempask = append(tempask, asks[i])	
					break
				}
			}
			j += 1
		}
		i += 1
	}

	fmt.Println("ASK Sequance", tempask)
	
	var tempbid []Orderbook

	// Rising
	var k int = 0
	for k < len(bids) {
		var l int = k + 1
		var countK int = 0
		for l <= len(bids) {
			countK += 1
			if l == len(bids) {
				if bids[k].Price <= tempbid[len(tempbid) - 1].Price {
					tempbid = append(tempbid, bids[k])	
					break
				}
			}
			if bids[l].Price <= bids[k].Price {
				if countK == 1 {
					tempbid = append(tempbid, bids[k])	
					break
				}	
			}
			l += 1
		}
		k += 1
	}

	fmt.Println("BID Sequance", tempbid)
}
