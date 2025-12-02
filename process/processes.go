package process

import (
	"database/sql"
	// "encoding/json"
	"fmt"
	"crypto/rand"
	"math/big"
	"strings"
	"log"
	// "marketdata/main/types"
	"time"
	// "math"
)

/*func margins(price float32) (float64, float64) {
	if price >= 100 && price <= 400 {
		return 0.4, 0.9
	} else if price >= 401 && price <= 700 {
		return 0.25, 0.7
	} else if price >= 701 && price <= 1000 {
		return 0.1, 0.4
	}
	return 0, 0
}*/



func generateCustomID() (string, error) {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var letterPart strings.Builder
	for i := 0; i < 5; i++ {
		letterIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		letterPart.WriteByte(letters[letterIndex.Int64()])
	}

	// Generate the number part (10 random digits)
	var numberPart strings.Builder
	for i := 0; i < 10; i++ {
		// Generate a random digit between 0 and 9
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		numberPart.WriteString(fmt.Sprintf("%d", digit.Int64()))
	}

	// Combine letter part, timestamp, and number part
	timestamp := time.Now().UnixNano() // Use UnixNano for high precision timestamp
	return fmt.Sprintf("%s%d%s", letterPart.String(), timestamp, numberPart.String()), nil
}

func timeconvert(lastTradedTime uint32) string {
	epochTime := int64(lastTradedTime)
    t := time.Unix(epochTime, 0) // This is the UTC time

	
	utcTime := t.UTC()

    // Debugging: Print the raw UTC time
    // fmt.Println("Raw UTC time:", t)

    // Format the time into a readable string
    readableTime := utcTime.Format("2006-01-02 15:04:05")

    return readableTime
}

/*type Order struct {
	InstrumentName string `json:"InstrumentName"`
	Type string           `json:"Type"`
	Quantity uint32       `json:"Quantity"`
	QuotePrice float32    `json:"QuotePrice"`
	Price float32         `json:"Price"`
	Stoploss float32      `json:"Stoploss"`
	Squareoff float32	  `json:"Squareoff"`
	Slippage float32      `json:"Slippage"`
}*/

/*type Orderbook struct {
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
}*/

/*type Instruments struct {
	Bids []Orderbook
	Asks []Orderbook
}*/

type Instruments struct {
	Bids [][]Orderbook
	Asks [][]Orderbook
	Bidquant int32
	Askquant int32
	Bidquantperc float64
	Askquantperc float64
}

func squareoff(connection *sql.DB, toporder Levels1) {

	var orders []Order

	var today string = generateDailyTableName()
	var table string = fmt.Sprintf("orders_%s", today)
	var InstrumentName string = toporder.InstrumentName

	var statement string = fmt.Sprintf("SELECT * FROM %s WHERE instrument = '%s' AND flag = 'ACTIVE'", table, InstrumentName)	
	rows, err := connection.Query(statement)
	
	if err != nil {
		log.Printf("Error fetchin data:", err)
	}

	defer rows.Close()

	for rows.Next() {
		
		var odr Order
		var timestamp string
		var flag string
		err = rows.Scan(&odr.Id, &odr.InstrumentName, &odr.Type, &odr.Quantity, &odr.QuotePrice, &odr.Price, &odr.Stoploss, &odr.Squareoff, &odr.Sp, &odr.Slippage, &flag, &timestamp)

		if err != nil {
			fmt.Println("Error:", err)
		}

		orders = append(orders, odr)
	}

	if len(orders) == 1 {
		var order Order
		order = orders[0]

		if order.Type == "LONG" {
			if order.Stoploss >= toporder.Ltp {
				// delete order, squareoff order

				fmt.Println("Complete Square off TOPORDER:", toporder)
				fmt.Println("Complete Square off ORDER:", order)
				fmt.Println("Complete Square off LTP:", toporder.Ltp)

				var result sql.Result
				var deletestatement string = fmt.Sprintf("UPDATE %s SET flag = 'SETTLED', settledprice = ? WHERE id = %d", table, order.Id)
				result, err = connection.Exec(deletestatement, toporder.Ltp)
			
				if err !=  nil {
					fmt.Println("Error deleting order:", err)
				}

				affected, err := result.RowsAffected()
			
				if err != nil {
					fmt.Println("Error getting affected row count:", affected)
				}

				if affected > 0 {
					fmt.Println("Order settled!")
				} else {
					fmt.Println("Problem occured in the order settlement!")
				}

				// Integrate the DhanHQ api to settle order here
			}

			if order.Squareoff <= toporder.Ltp {

				fmt.Println("Update Square off TOPORDER:", toporder)
				fmt.Println("Update Square off ORDER:", order)

				// Update the stoploss to squareoff - 0.1% and squareoff to +0.3 percnt 
				var stoploss float32 = toporder.Ltp -  ((toporder.Ltp * 0.1) / 100)
				var squareoff float32 = toporder.Ltp + ((toporder.Ltp * 0.3) / 100)

				// Integrate the DhanHQ api to modify the order
			
				var updatestatement string = fmt.Sprintf("UPDATE %s SET stoploss = ?, squareoff = ? WHERE id = %d", table, order.Id)
				result, err := connection.Exec(updatestatement, stoploss, squareoff)

				if err != nil {
					fmt.Errorf("Error executing update: %w", err)
				}

				affected, err := result.RowsAffected()

				if err != nil {
					fmt.Errorf("Error checking affected rows: %w", err)
				}
			
				if affected > 0 {
					fmt.Println("Updated!!!")
				} 
			}
		}

		if order.Type == "SHORT" {
			if toporder.Ltp >= order.Stoploss {
				// delete order, squareoff order

				fmt.Println("Complete Square off TOPORDER:", toporder)
				fmt.Println("Complete Square off ORDER:", order)
				fmt.Println("Complete Square off LTP:", toporder.Ltp)

				var result sql.Result
				var deletestatement string = fmt.Sprintf("UPDATE %s SET flag = 'SETTLED', settledprice = ? WHERE id = %d", table, order.Id)
				result, err = connection.Exec(deletestatement, toporder.Ltp)
			
				if err !=  nil {
					fmt.Println("Error deleting order:", err)
				}

				affected, err := result.RowsAffected()
			
				if err != nil {
					fmt.Println("Error getting affected row count:", affected)
				}

				if affected > 0 {
					fmt.Println("Order settled!")
				} else {
					fmt.Println("Problem occured in the order settlement!")
				}

				// Integrate the DhanHQ api to settle order here
			}

			if toporder.Ltp <= order.Squareoff {
				fmt.Println("Update Square off TOPORDER:", toporder)
				fmt.Println("Update Square off ORDER:", order)

				var stoploss float32 = toporder.Ltp +  ((toporder.Ltp * 0.1) / 100)
				var squareoff float32 = toporder.Ltp - ((toporder.Ltp * 0.3) / 100)

				// Integrate the DhanHQ api to modify the order

				var updatestatement string = fmt.Sprintf("UPDATE %s SET stoploss = ?, squareoff = ? WHERE id = %d", table, order.Id)
				result, err := connection.Exec(updatestatement, stoploss, squareoff)

				if err != nil {
					fmt.Errorf("Error executing update: %w", err)
				}

				affected, err := result.RowsAffected()

				if err != nil {
					fmt.Errorf("Error checking affected rows: %w", err)
				}
			
				if affected > 0 {
					fmt.Println("Updated!!!")
				} 
			}
		}
	}
}


var sequence = make(map[string]Instruments)

var orderno int64 = 0

var bidquant int32 = 0
var askquant int32 = 0

var bidquantperc float64
var askquantperc float64

func existingorder(connection *sql.DB, order Order) []Order {
	
	var placed []Order

	var today string = generateDailyTableName()
	var table string = fmt.Sprintf("orders_%s", today)
	var statement string = fmt.Sprintf("SELECT * FROM %s WHERE instrument = '%s' AND flag = 'ACTIVE'", table, order.InstrumentName)	
	
	rows, err := connection.Query(statement)
	
	if err != nil {
		log.Printf("Error fetchin data:", err)
	}

	defer rows.Close()

	for rows.Next() {
		
		var odr Order
		var timestamp string
		var flag string
		err = rows.Scan(&odr.Id, &odr.InstrumentName, &odr.Type, &odr.Quantity, &odr.QuotePrice, &odr.Price, &odr.Stoploss, &odr.Squareoff, &odr.Sp, &odr.Slippage, &flag, &timestamp)

		if err != nil {
			fmt.Println("Error:", err)
		}

		placed = append(placed, odr)
	}
	
	return placed
}


func placeorder(connection *sql.DB, order Order) {

	var orders []Order
	orders = existingorder(connection, order)
	var today string = generateDailyTableName()
	var table string = fmt.Sprintf("orders_%s", today)

	tradetimeunix := uint32(time.Now().Unix())
	tradetime := timeconvert(tradetimeunix)

	if len(orders) == 0 {
		// API Order
		
	
		// Database Order
		var statement string = fmt.Sprintf("INSERT INTO %s (instrument, type, quantity, quote, price, stoploss, squareoff, settledprice, slippage, flag, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", table)
		op, err := connection.Exec(statement, order.InstrumentName, order.Type, order.Quantity, order.QuotePrice, order.Price, order.Stoploss, order.Squareoff, order.Sp, order.Slippage, "ACTIVE", tradetime)

		if err != nil {
			log.Printf("error inserting into order table: %v", err)
			return
		}
		
		fmt.Println("\nOrder successfully inserted into the table.", op)
	} 
}

func processes(connection *sql.DB, order Levels1) {

	squareoff(connection, order)

	var bidorder Orderbook
	bidorder.OrderNo = orderno + 1
	bidorder.InstrumentName = order.InstrumentName
	bidorder.LastTradedTime = order.LastTradedTime
	bidorder.Ltt = timeconvert(order.LastTradedTime)
	bidorder.Type = "Bid"
	bidorder.Price = order.BidPrice
	bidorder.Quantity = order.BidQuantity
	bidorder.NoOfOrders = order.NoOfBidOrders
	bidorder.Spread = order.AskPrice - order.BidPrice
	bidorder.Ltp = order.Ltp

	var askorder Orderbook
	askorder.OrderNo = orderno + 1
	askorder.InstrumentName = order.InstrumentName
	askorder.LastTradedTime = order.LastTradedTime
	askorder.Ltt = timeconvert(order.LastTradedTime)
	askorder.Type = "Ask"
	askorder.Price = order.AskPrice
	askorder.Quantity = order.AskQuantity
	askorder.NoOfOrders = order.NoOfAskOrders
	askorder.Spread = order.AskPrice - order.BidPrice
	askorder.Ltp = order.Ltp

	_, ok := sequence[order.InstrumentName];
	if !ok {
		sequence[order.InstrumentName] = Instruments{
			Bids: [][]Orderbook{},
			Asks: [][]Orderbook{},
			Bidquant: 0,
			Askquant: 0,	
			Bidquantperc: 0,
			Askquantperc: 0,
		}
	}

	io := sequence[order.InstrumentName]

	if bidorder.Type == "Bid" {
		if len(io.Bids) > 0 {
			lastsequence := io.Bids[len(io.Bids) - 1]
			lastrecord := lastsequence[len(lastsequence) - 1]

			if lastrecord == bidorder {
				return	
			}

			if lastrecord.Price < bidorder.Price {
				lastsequence = append(lastsequence, bidorder)
				io.Bids[len(io.Bids) - 1] = lastsequence
			} else {

				if len(lastsequence) > 3 {
					/*fmt.Print("---- Last Bid Sequence End---- and a new sequence starts ----")
					fmt.Println(lastsequence)*/

					for _, row := range lastsequence {
						// fmt.Println("ROW", row)
						io.Bidquant += row.Quantity
					}

					if io.Bidquant > io.Askquant {
						if io.Bidquant > 2*io.Askquant {
							var longorder Order
							longorder.InstrumentName = bidorder.InstrumentName
							longorder.Type = "LONG"	
							longorder.Price = bidorder.Price
							// longorder.Quantity = positionsizing(connection, longorder)
							longorder.Quantity = 0
							longorder.QuotePrice = 0.0
							longorder.Stoploss = bidorder.Price - ((bidorder.Price * 0.3) / 100)
							longorder.Squareoff = bidorder.Price + ((bidorder.Price * 0.7) / 100)
							longorder.Sp = bidorder.Ltp
							longorder.Slippage = 0.0

							placeorder(connection, longorder)
						}
					}
				}
								
				io.Bids = append(io.Bids, []Orderbook{bidorder})
			}

		} else {
			io.Bids = append(io.Bids, []Orderbook{bidorder})
		}
	}

	if askorder.Type == "Ask" {
		if len(io.Asks) > 0 {
			lastsequence := io.Asks[len(io.Asks) - 1]
			lastrecord := lastsequence[len(lastsequence) - 1]

			if lastrecord == askorder {
				return
			}
			
			if lastrecord.Price > askorder.Price {
				lastsequence = append(lastsequence, askorder)
				io.Asks[len(io.Asks) - 1] = lastsequence
			} else {
				
				if len(lastsequence) > 3 {
					/*fmt.Print("---- Last Ask Sequence End---- and a new sequence starts ----")
					fmt.Println(lastsequence)*/

					// fmt.Println("asks", lastsequence)
					for _, row := range lastsequence {
						io.Askquant += row.Quantity
					}
					
					if io.Askquant > io.Bidquant {

						if io.Askquant > 2*io.Bidquant { 
							var shortorder Order
							shortorder.InstrumentName = askorder.InstrumentName
							shortorder.Type = "SHORT"
							shortorder.Price = askorder.Price
							shortorder.Quantity = 0
							shortorder.QuotePrice = 0.0
							shortorder.Stoploss = askorder.Price + ((askorder.Price * 0.3) / 100)
							shortorder.Squareoff = askorder.Price - ((askorder.Price * 0.7) / 100)
							shortorder.Sp = askorder.Ltp
							shortorder.Slippage = 0.0

							placeorder(connection, shortorder)
						}
					}
				}

				io.Asks = append(io.Asks, []Orderbook{askorder})
			}
			
		} else {
			io.Asks = append(io.Asks, []Orderbook{askorder})
		}
	}

	sequence[order.InstrumentName] = io
}
