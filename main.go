package main

import (
    // "bytes"
    "encoding/binary"
	// "encoding/json"
	"math"
	// "time"
    "fmt"

	// Custom Packge Imports
	"process/process"
	
	"github.com/gorilla/websocket"
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


type Levels struct {
	BidQuantity int32
	AskQuantity int32
	NoOfBidOrders int16
	NoOfAskOrders int16
	BidPrice float32
	AskPrice float32
}

const (
	Workers = 10
	ChannelSize = 1000
)


func main() {

	var token string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJkaGFuIiwicGFydG5lcklkIjoiIiwiZXhwIjoxNzYxNjI0MTM3LCJpYXQiOjE3NjE1Mzc3MzcsInRva2VuQ29uc3VtZXJUeXBlIjoiU0VMRiIsIndlYmhvb2tVcmwiOiIiLCJkaGFuQ2xpZW50SWQiOiIxMTA4ODcwNTEwIn0.bpa7zX1dMdZXEz2CmVoqYv31zOO5MqJxfO9VRXYsens7k05psCnVNwBKQ4wOaMGCiNcOW9Qs8G7VH0uGnK4UfA"
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
	
	err = c.WriteJSON(instrumentList)
	
	if err != nil {
		panic(err)
	}

	// Workerpool	
	var messages chan []byte = make(chan []byte, ChannelSize)
	fmt.Println("Type of message channel:", messages)
	
	var i int
	for i = 0; i < Workers; i++ {
		go worker(i, messages)
	}

	for {
		_, data, err := c.ReadMessage()	
		
		if err != nil {
			fmt.Println("Error:", err)
		}	
		
		fmt.Println("Data:", data)

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

	fmt.Println("LENGTH:", len(data))
	fmt.Println("LENGTH Before Depth:", len(data[0:62]))
		
	// Extracting Response Header
	feedResponseCode := data[0]
	messageLength := binary.LittleEndian.Uint16(data[1:3])
	exchangeSegment := data[3]
	securityID := binary.LittleEndian.Uint32(data[4:8])

	// PacketData
	bits := binary.LittleEndian.Uint32(data[8:12])
    lastTradedPrice := math.Float32frombits(bits)

	lastTradedQuantity := binary.LittleEndian.Uint16(data[12:14])
	lastTradedTime := binary.LittleEndian.Uint32(data[14:18])

	bitsAP := binary.LittleEndian.Uint32(data[18:22])
    averageTradePrice := math.Float32frombits(bitsAP)

	volume := binary.LittleEndian.Uint32(data[22:26])
	totalSellQuantity := binary.LittleEndian.Uint32(data[26:30])
	totalBuyQuantity := binary.LittleEndian.Uint32(data[30:34])
	oi := binary.LittleEndian.Uint32(data[34:38])
	hoi := binary.LittleEndian.Uint32(data[38:42])
	loi := binary.LittleEndian.Uint32(data[42:46])

	bitsOV := binary.LittleEndian.Uint32(data[46:50])
    dayOpenValue := math.Float32frombits(bitsOV)

	bitsCV := binary.LittleEndian.Uint32(data[50:54])
    dayCloseValue := math.Float32frombits(bitsCV)

	bitsHV := binary.LittleEndian.Uint32(data[54:58])
    dayHighValue := math.Float32frombits(bitsHV)

	bitsLV := binary.LittleEndian.Uint32(data[58:62])
   	dayLowValue := math.Float32frombits(bitsLV)
		
	fmt.Println("Data Length:", len(data))	
	var levels []Levels
	marketDepthData := data[62:len(data)]
	fmt.Println("Market Depth Data Length:", len(marketDepthData))
	var i int = 0

	for i < len(marketDepthData) {
		var level Levels

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
		fmt.Println("I values:", i)
	}
		
	fmt.Println("feedResponseCode:", feedResponseCode)
	fmt.Println("MessageLength:", messageLength)

	fmt.Println("ExchangeSegment:", exchangeSegment)
	fmt.Println("securityID:", securityID)

	// PacketData
	// fmt.Println("LastTradedPrice:", lastTradedPrice)
	fmt.Println("LastTradedPrice:", lastTradedPrice)
	fmt.Println("LastTradedQuantity:", lastTradedQuantity)
	fmt.Println("LastTradedTime:", lastTradedTime)
	fmt.Println("AverageTradePrice:", averageTradePrice)
	fmt.Println("Volume:", volume)

	fmt.Println("totalSellQuantity:", totalSellQuantity)
	fmt.Println("totalBuyQuantity:", totalBuyQuantity)
	fmt.Println("Oi:", oi)
	fmt.Println("Hoi:", hoi)
	fmt.Println("Loi:", loi)
	fmt.Println("dayOpenValue:", dayOpenValue)
	fmt.Println("dayCloseValue:", dayCloseValue)
	fmt.Println("dayHighValue:", dayHighValue)
	fmt.Println("dayLowValue:", dayLowValue)

	fmt.Println("Market Depth Data:", marketDepthData)
	fmt.Println("--- Depth ---")
		
	for _, lvl := range levels {
		fmt.Printf("BidQuantity: %d, AskQuantity: %d, NoOfBidOrders: %d, NoOfAskOrders: %d, BidPrice: %.2f, AskPrice: %.2f\n", lvl.BidQuantity, lvl.AskQuantity, lvl.NoOfBidOrders, lvl.NoOfAskOrders, lvl.BidPrice, lvl.AskPrice)
	}

	fmt.Println("-----------------------------")
}

