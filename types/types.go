package types

type Levels5 struct {
	BidQuantity int32
	AskQuantity int32
	NoOfBidOrders int16
	NoOfAskOrders int16
	BidPrice float32
	AskPrice float32
}

type FullPacket struct {
	InstrumentName string
    FeedResponseCode  byte
    MessageLength     uint16
    ExchangeSegment   byte
    SecurityID        uint32
    LastTradedPrice   float32
    LastTradedQuantity uint16
    LastTradedTime    uint32
    AverageTradePrice float32
    Volume            uint32
    TotalSellQuantity uint32
    TotalBuyQuantity  uint32
    OI                uint32
    HOI               uint32
    LOI               uint32
    DayOpenValue      float32
    DayCloseValue     float32
    DayHighValue      float32
    DayLowValue       float32
    Levels5            []Levels5
}

type Levels20 struct {
	Price float64
	Quantity uint32
	NoOfOrders uint32
}

type Orderbook struct {
	InstrumentName string
    LastTradedTime uint32
	PacketType string
	MessageLength uint16
	FeedResponseCode byte
	ExchangeSegment byte
	SecurityID uint32
	MessageSequance uint32
	Levels []Levels20
}

