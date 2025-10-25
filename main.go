package main

import (
    // "bytes"
    "encoding/binary"
	// "encoding/json"
	"math"
	// "time"
    "fmt"
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

func main() {

	var token string = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJkaGFuIiwicGFydG5lcklkIjoiIiwiZXhwIjoxNzYxMzY3Mjg3LCJpYXQiOjE3NjEyODA4ODcsInRva2VuQ29uc3VtZXJUeXBlIjoiU0VMRiIsIndlYmhvb2tVcmwiOiIiLCJkaGFuQ2xpZW50SWQiOiIxMTA4ODcwNTEwIn0.RdAOJGqfyPnRNIhT0Sbi167PY1Al2SPRSIZVtu0Xy2CXAmewlrLyYn74EO6iABkCBNTs-lppn2GyQuR0zyVNIA"
	
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
		InstrumentCount: 1,
		InstrumentList: []InstrumentList{
			{
				ExchangeSegment: "NSE_EQ",
				SecurityId:      "1333", //   NSE,E,1333,INE040A01034,EQUITY,,HDFCBANK,HDFC BANK LTD,HDFC Bank,ES,EQ,1.0,,,,5.0000,NA,N,N,N,NA,A,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,4.545455,
			},
		},
	}

	
	// jsonData, err := json.Marshal(instrumentList)	

	if err != nil {
		panic(err)
	}

	err = c.WriteJSON(instrumentList)
	
	if err != nil {
		panic(err)
	}

	for {
		/*_, data, err := c.ReadMessage()	
		
		if err != nil {
			panic(err)
			break
		}	
		
		fmt.Println("Data:", data)*/

		data := []byte{8, 162, 0, 1, 53, 5, 0, 0, 154, 89, 120, 68, 35, 0, 20, 144, 251, 104, 72, 49, 122, 68, 58, 220, 211, 0, 0, 78, 17, 0, 152, 4, 10, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 192, 124, 68, 51, 51, 124, 68, 205, 252, 124, 68, 51, 211, 119, 68, 15, 0, 0, 0, 1, 0, 0, 0, 1, 0, 1, 0, 154, 89, 120, 68, 0, 96, 120, 68, 241, 0, 0, 0, 1, 0, 0, 0, 2, 0, 1, 0, 102, 86, 120, 68, 102, 102, 120, 68, 118, 1, 0, 0, 158, 1, 0, 0, 2, 0, 3, 0, 51, 83, 120, 68, 154, 105, 120, 68, 210, 0, 0, 0, 201, 0, 0, 0, 3, 0, 3, 0, 0, 80, 120, 68, 205, 108, 120, 68, 69, 5, 0, 0, 232, 0, 0, 0, 9, 0, 3, 0, 205, 76, 120, 68, 0, 112, 120, 68}
		
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
}

