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

/*// DepthPacket represents a single depth packet (16 bytes)
type DepthPacket struct {
	Price      float64
	Quantity   uint32
	NumOrders  uint32
}

// MarketDepthPacket represents the full packet (header + 20 depth packets)
type MarketDepthPacket struct {
	MessageLength   uint16
	FeedResponseCode byte
	ExchangeSegment byte
	SecurityID      int32
	MessageSequence uint32
	Depth           []DepthPacket
}

func parseDepthPacket(data []byte, isLittleEndian bool) ([]MarketDepthPacket, error) {
	var packets []MarketDepthPacket
	offset := 0

	for offset+12 <= len(data) {
		// Parse Response Header (12 bytes)
		var packet MarketDepthPacket

		// Message Length (int16, little-endian)
		packet.MessageLength = binary.LittleEndian.Uint16(data[offset : offset+2])
		if int(packet.MessageLength) > len(data)-offset {
			return nil, fmt.Errorf("invalid message length %d at offset %d", packet.MessageLength, offset)
		}

		// Feed Response Code (byte)
		packet.FeedResponseCode = data[offset+2]
		if packet.FeedResponseCode != 41 && packet.FeedResponseCode != 51 {
			return nil, fmt.Errorf("invalid feed response code %d at offset %d", packet.FeedResponseCode, offset)
		}

		// Exchange Segment (byte)
		packet.ExchangeSegment = data[offset+3]

		// Security ID (int32, big-endian)
		packet.SecurityID = int32(binary.BigEndian.Uint32(data[offset+4 : offset+8]))

		// Message Sequence (uint32, big-endian, ignored)
		packet.MessageSequence = binary.BigEndian.Uint32(data[offset+8 : offset+12])

		// Parse 20 Depth Packets (320 bytes)
		packet.Depth = make([]DepthPacket, 20)
		depthOffset := offset + 12

		for i := 0; i < 20; i++ {
			if depthOffset+16 > len(data) {
				return nil, fmt.Errorf("incomplete depth data at offset %d", depthOffset)
			}

			// Price (float64)
			var priceBits uint64
			if isLittleEndian {
				priceBits = binary.LittleEndian.Uint64(data[depthOffset : depthOffset+8])
			} else {
				priceBits = binary.BigEndian.Uint64(data[depthOffset : depthOffset+8])
			}
			packet.Depth[i].Price = math.Float64frombits(priceBits)

			// Quantity (uint32)
			packet.Depth[i].Quantity = binary.BigEndian.Uint32(data[depthOffset+8 : depthOffset+12])

			// No. of Orders (uint32)
			packet.Depth[i].NumOrders = binary.BigEndian.Uint32(data[depthOffset+12 : depthOffset+16])

			depthOffset += 16
		}

		packets = append(packets, packet)
		offset += int(packet.MessageLength)
	}

	return packets, nil
}*/


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

		marketDepthData := data[62:162]

		bidQuantity := binary.LittleEndian.Uint32(marketDepthData[0:4])
		askQuantity := binary.LittleEndian.Uint32(marketDepthData[4:8])

		NoOfBidOrder := binary.LittleEndian.Uint16(marketDepthData[8:10])
		NoOfAskOrder := binary.LittleEndian.Uint16(marketDepthData[10:12])

		bidP := binary.LittleEndian.Uint32(marketDepthData[12:16])
    	bidPrice := math.Float32frombits(bidP)

		askP := binary.LittleEndian.Uint32(marketDepthData[16:20])
    	askPrice := math.Float32frombits(askP)
		

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

		fmt.Println("BidQuantity:", bidQuantity)
		fmt.Println("AskQuantity:", askQuantity)
		fmt.Println("NoOfBidOrder:", NoOfBidOrder)
		fmt.Println("NoOfAskOrder:", NoOfAskOrder)
		fmt.Println("BidPrice:", bidPrice)
		fmt.Println("AskPrice:", askPrice)

		fmt.Println("-----------------------------")

		

		/*packets, err := parseDepthPacket(data, true) // Try little-endian for float64
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		for i, packet := range packets {
			fmt.Printf("Packet %d:\n", i+1)
			fmt.Printf("  Message Length: %d bytes\n", packet.MessageLength)
			fmt.Printf("  Feed Response Code: %d (%s)\n", packet.FeedResponseCode, mapFeedCode(packet.FeedResponseCode))
			fmt.Printf("  Exchange Segment: %d\n", packet.ExchangeSegment)
			fmt.Printf("  Security ID: %d\n", packet.SecurityID)
			fmt.Printf("  Message Sequence: %d (ignored)\n", packet.MessageSequence)
			fmt.Printf("  Depth Data:\n")
			for j, depth := range packet.Depth {
				fmt.Printf("    Level %d:\n", j+1)
				fmt.Printf("      Price: %.2f\n", depth.Price)
				fmt.Printf("      Quantity: %d\n", depth.Quantity)
				fmt.Printf("      No. of Orders: %d\n", depth.NumOrders)
			}
			fmt.Println()
		}

		/*feedResponseCode := data[0]
		messageLength := binary.BigEndian.Uint16(data[1:3])
		exchangeSegment := data[3]
		securityID := binary.BigEndian.Uint32(data[4:8])

		fmt.Printf("Feed Response Code: %d\n", feedResponseCode)
		fmt.Printf("Message Length: %d bytes\n", messageLength)
		fmt.Printf("Exchange Segment: %d\n", exchangeSegment)
		fmt.Printf("Security ID: %d\n", securityID)

		bits := binary.LittleEndian.Uint32(data[8:12])
		lastTradedPrice := math.Float32frombits(bits)
		fmt.Printf("Last Traded Price: %.2f\n", lastTradedPrice)

		lastTradeTime := binary.LittleEndian.Uint32(data[12:16])
		timestamp := time.Unix(int64(lastTradeTime), 0).UTC()
		fmt.Printf("Last Trade Time: %s (Epoch: %d)\n", timestamp.Format(time.RFC3339), lastTradeTime)*/
	}
}

/*func mapFeedCode(code byte) string {
	switch code {
	case 41:
		return "Bid"
	case 51:
		return "Ask"
	default:
		return "Unknown"
	}
}*/
