package main

import (
    // "bytes"
    "encoding/binary"
	// "encoding/json"
	"math"
	// "time"
    "fmt"

	// Custom Packge Imports
	// "marketdata/main/process"
	"marketdata/main/types"
	
	"github.com/gorilla/websocket"
)


func Marketdepth() {

	var token string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJkaGFuIiwicGFydG5lcklkIjoiIiwiZXhwIjoxNzYxOTY5NjkzLCJpYXQiOjE3NjE4ODMyOTMsInRva2VuQ29uc3VtZXJUeXBlIjoiU0VMRiIsIndlYmhvb2tVcmwiOiIiLCJkaGFuQ2xpZW50SWQiOiIxMTA4ODcwNTEwIn0.ldtrKVlUu755WjecWwchB9mWzBcPPUcNnOjmLNxdVf7m63UKH42lYcCvqhpZVTGfTRQl2lIAvh_ssXN0LRC7iA"

	var clientId string = "1108870510"
    var url string
	
	// 20 Level Depth URL
    url = fmt.Sprintf("wss://depth-api-feed.dhan.co/twentydepth?token=%s&clientId=%s&authType=2", token, clientId)

	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	
	if err != nil {
		fmt.Println("Error:", err)
	}

	defer c.Close()	

	instrumentList := Instrument{
		RequestCode: 23,
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
		fmt.Println("Error:", err)
	}

	// Workerpool	
	var messages chan []byte = make(chan []byte, ChannelSize)
	
	var i int
	for i = 0; i < Workers; i++ {
		go workerMD(i, messages)
	}

	for {
		_, data, err := c.ReadMessage()	
		
		if err != nil {
			fmt.Println("Error:", err)
		}

		// fmt.Println("Example Packets in stream:", data)
		
		select {
			case messages <- data:	
			default:
				fmt.Println("Dropping Messages (channel full)")
		}
	}
}

func workerMD(id int, ch <- chan []byte) {
	for data := range ch {
		parsingMarketDepthData(data)
	}
}

func parsingMarketDepthData(data []byte) {
	fmt.Println("Example Packet:", data)
	var orderbook types.Orderbook
	
	// var messagelength int16
	orderbook.MessageLength = binary.LittleEndian.Uint16(data[0:2])
	orderbook.FeedResponseCode = data[2]
	orderbook.ExchangeSegment = data[3]
	orderbook.SecurityID = binary.LittleEndian.Uint32(data[4:8])
	orderbook.MessageSequance = binary.LittleEndian.Uint32(data[8:12])

	var instrumentSecurityID uint32
	instrumentSecurityID = binary.LittleEndian.Uint32(data[4:8])

	if orderbook.FeedResponseCode == 41 {
		orderbook.PacketType = "Bid"	
	} else if orderbook.FeedResponseCode == 51 {
		orderbook.PacketType = "Ask"	
	} else {
		orderbook.PacketType = "_"
	}

	switch instrumentSecurityID {
		case 163:
			orderbook.InstrumentName = "APOLLO TYRE"	
		case 1333:
			orderbook.InstrumentName = "HDFC BANK"
		case 11723:
			orderbook.InstrumentName = "JSWSTEEL"
		case 19020:
			orderbook.InstrumentName = "JSWINFRA"
		default:
			orderbook.InstrumentName = "None"
	}

	var levels []types.Levels20
	depthData := data[12:orderbook.MessageLength]
	var i int = 0

	for i < len(depthData) {
		var level types.Levels20

		price  := binary.LittleEndian.Uint64(depthData[i : i+8])
    	Price := math.Float64frombits(price)
		level.Price = Price
		i += 8

		level.Quantity = binary.LittleEndian.Uint32(depthData[i : i+4])
		i += 4

		level.NoOfOrders = binary.LittleEndian.Uint32(depthData[i : i+4])
		i += 4
			
		levels = append(levels, level)
	}

	orderbook.Levels = levels
	fmt.Println("Orderbook:", orderbook)
}
