package main

import (
	// "bytes"
	"database/sql"
	"encoding/binary"

	// "encoding/json"
	"math"
	"time"
	"fmt"
	"os"

	// Custom Packge Imports
	"marketdata/main/process"
	"marketdata/main/types"

	"github.com/gorilla/websocket"
	_ "github.com/mattn/go-sqlite3"
)

type InstrumentList struct {
	ExchangeSegment string `json:"ExchangeSegment"`
	SecurityId string `json:"SecurityId"`
}

type Instrument struct {
	RequestCode int32 `json:"RequestCode"`
	InstrumentCount int32 `json:"InstrumentCount"`	
	InstrumentList []InstrumentList `json:"InstrumentList"`
}

const (
	Workers = 14
	ChannelSize = 1024
)


var database *sql.DB

func generateDailyTableName() string {
    now := time.Now() // Type: time.Time
    return now.Format("2006_01_02") // Type: string
}

func createtables(database *sql.DB) error {
	// Generate the table name for today
	var today string
    today = generateDailyTableName() // Type: string
	
	var orderbookTableName = fmt.Sprintf("orderbook_%s", today) // concat with today
	var marketbookTableName = fmt.Sprintf("marketbook_%s", today) // concat with today
    
    orderbookSQL := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            data TEXT,
            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`, orderbookTableName) // Type: string (the full SQL command)

	marketbookSQL := fmt.Sprintf(`
        CREATE TABLE IF NOT EXISTS %s (
            id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
            data TEXT,
            timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );`, marketbookTableName) 

    fmt.Printf("Attempting to verify/create table: %s\n", orderbookTableName)
    fmt.Printf("Attempting to verify/create table: %s\n", marketbookTableName)
    
    if _, err := database.Exec(orderbookSQL); err != nil {
		return fmt.Errorf("failed to create table %s: %w", orderbookTableName, err)
	}

	// Use = for the second assignment of err (FIXED: Go syntax for error handling)
	if _, err := database.Exec(marketbookSQL); err != nil {
		return fmt.Errorf("failed to create table %s: %w", marketbookTableName, err)
	}

	return nil
}

func main() {
	
	/*const databasePath string = "./db/development.db"
	const dir string = "./db"*/

	const databasePath string = "../databases/development.db"
	const dir string = "./databases"

	var err error

	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println("Error creating database directory %s: %v", dir, err)
	}

	database, err = sql.Open("sqlite3", databasePath)

	if err != nil {
		fmt.Println("Error opening database: %v", err)	
	}
	
	var closeErr error

	defer func() {
		if closeErr = database.Close(); closeErr != nil {
			fmt.Println("Error closing database: %v", closeErr)
		}
	}()

	fmt.Println("Successfully connected to SQLite database:", databasePath)

	if err = createtables(database); err != nil {
		fmt.Println("Error creating tables: %v", err)
	}
	if err != nil {
		fmt.Println("Error creating tables %v", err)
	}
	

	var token string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJkaGFuIiwicGFydG5lcklkIjoiIiwiZXhwIjoxNzYxOTY5NjkzLCJpYXQiOjE3NjE4ODMyOTMsInRva2VuQ29uc3VtZXJUeXBlIjoiU0VMRiIsIndlYmhvb2tVcmwiOiIiLCJkaGFuQ2xpZW50SWQiOiIxMTA4ODcwNTEwIn0.ldtrKVlUu755WjecWwchB9mWzBcPPUcNnOjmLNxdVf7m63UKH42lYcCvqhpZVTGfTRQl2lIAvh_ssXN0LRC7iA"
	var clientId string = "1108870510"
    var url string
	
	// Live Market Feed URL
    url = fmt.Sprintf("wss://api-feed.dhan.co?version=2&token=%s&clientId=%s&authType=2", token, clientId)

	// 20 Level Depth URL
    // url = fmt.Sprintf("wss://depth-api-feed.dhan.co/twentydepth?token=%s&clientId=%s&authType=2", token, clientId)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	
	if err != nil {
		panic(err)
	}

	defer c.Close()	
	
	instrumentList := Instrument{
		RequestCode:   21,
		InstrumentCount: 4,
		InstrumentList: []InstrumentList{
			{
				ExchangeSegment: "NSE_EQ",
				SecurityId:      "1333", //   NSE,E,1333,INE040A01034,EQUITY,,HDFCBANK,HDFC BANK LTD,HDFC Bank,ES,EQ,1.0,,,,5.0000,NA,N,N,N,NA,A,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,4.545455,
			},
			{
				ExchangeSegment: "NSE_EQ",
				SecurityId:      "11723", 
			},
			{
				ExchangeSegment: "NSE_EQ",
				SecurityId:      "19020" , 
			},
			{
				ExchangeSegment: "NSE_EQ",
				SecurityId:      "163", 
			},
		},
	}

	// go Marketdepth()
	
	err = c.WriteJSON(instrumentList)
	
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Workerpool	
	var messages chan []byte = make(chan []byte, ChannelSize)
	// fmt.Println("Type of message channel for FullPacketData:", messages)
	
	var i int
	for i = 0; i < Workers; i++ {
		go worker(i, messages)
	}

	for {
		_, data, err := c.ReadMessage()	
		
		if err != nil {
			fmt.Println("Error:", err)
		}	
		
		// fmt.Println("Data:", data)
		// fmt.Println("Data Length:", len(data))

		select {
			case messages <- data:	
			default:
				fmt.Println("Dropping Messages (channel full)")
		}
	}
}

func worker(id int, ch <-chan []byte) {
	for data := range ch {
		parsing(data)
	}	
}

func parsing(data []byte) {

	var fullDataFeed types.FullPacket
	// var instrumentName string
	var instrumentSecurityID uint32

	instrumentSecurityID = binary.LittleEndian.Uint32(data[4:8])

	// fmt.Println("type of %T\n", instrumentSecurityID)

	switch instrumentSecurityID {
		case 163:
			fullDataFeed.InstrumentName = "APOLLO TYRE"	
		case 1333:
			fullDataFeed.InstrumentName = "HDFC BANK"
		case 11723:
			fullDataFeed.InstrumentName = "JSWSTEEL"
		case 19020:
			fullDataFeed.InstrumentName = "JSWINFRA"
		default:
			fullDataFeed.InstrumentName = "None"
	}
		
	// Extracting Response Header
	fullDataFeed.FeedResponseCode = data[0]
	fullDataFeed.MessageLength = binary.LittleEndian.Uint16(data[1:3])
	fullDataFeed.ExchangeSegment = data[3]
	fullDataFeed.SecurityID = instrumentSecurityID

	// PacketData
	bits := binary.LittleEndian.Uint32(data[8:12])
    fullDataFeed.LastTradedPrice = math.Float32frombits(bits)

	fullDataFeed.LastTradedQuantity = binary.LittleEndian.Uint16(data[12:14])
	fullDataFeed.LastTradedTime = binary.LittleEndian.Uint32(data[14:18])

	bitsAP := binary.LittleEndian.Uint32(data[18:22])
    fullDataFeed.AverageTradePrice = math.Float32frombits(bitsAP)

	fullDataFeed.Volume = binary.LittleEndian.Uint32(data[22:26])
	fullDataFeed.TotalSellQuantity = binary.LittleEndian.Uint32(data[26:30])
	fullDataFeed.TotalBuyQuantity = binary.LittleEndian.Uint32(data[30:34])
	fullDataFeed.OI = binary.LittleEndian.Uint32(data[34:38])
	fullDataFeed.HOI = binary.LittleEndian.Uint32(data[38:42])
	fullDataFeed.LOI = binary.LittleEndian.Uint32(data[42:46])

	bitsOV := binary.LittleEndian.Uint32(data[46:50])
    fullDataFeed.DayOpenValue = math.Float32frombits(bitsOV)

	bitsCV := binary.LittleEndian.Uint32(data[50:54])
    fullDataFeed.DayCloseValue = math.Float32frombits(bitsCV)

	bitsHV := binary.LittleEndian.Uint32(data[54:58])
    fullDataFeed.DayHighValue = math.Float32frombits(bitsHV)

	bitsLV := binary.LittleEndian.Uint32(data[58:62])
   	fullDataFeed.DayLowValue = math.Float32frombits(bitsLV)
		
	var levels []types.Levels5
	marketDepthData := data[62:len(data)]
	var i int = 0

	for i < len(marketDepthData) {
		var level types.Levels5

		bidQuantity := binary.LittleEndian.Uint32(marketDepthData[i : i+4])
		level.BidQuantity = int32(bidQuantity)
		i += 4

		askQuantity := binary.LittleEndian.Uint32(marketDepthData[i : i+4])
		level.AskQuantity = int32(askQuantity)
		i += 4

		noOfBidOrder := binary.LittleEndian.Uint16(marketDepthData[i : i+2])
		level.NoOfBidOrders = int16(noOfBidOrder)
		i += 2

		noOfAskOrder := binary.LittleEndian.Uint16(marketDepthData[i : i+2])
		level.NoOfAskOrders = int16(noOfAskOrder)
		i += 2

		bidP := binary.LittleEndian.Uint32(marketDepthData[i : i+4])
    	bidPrice := math.Float32frombits(bidP)
		level.BidPrice = bidPrice
		i += 4

		askP := binary.LittleEndian.Uint32(marketDepthData[i : i+4])
    	askPrice := math.Float32frombits(askP)
		level.AskPrice = askPrice
		i += 4
			
		levels = append(levels, level)
	}
		
	fullDataFeed.Levels5 = levels
	process.Process(fullDataFeed, database)
}

