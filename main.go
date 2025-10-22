package main

import (
    "io/ioutil"
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type Payload struct {
    SecurityID string `json:"securityId"`
    ExchangeSegment string `json:"exchangeSegment"`
    Instrument string `json:"instrument"`
    Interval string `json:"interval"`
    Oi bool `json:"oi"`
    FromDate string `json:"fromDate"`
    ToDate string `json:"toDate"`
}

func main() {
    var url string
    url = "https://api.dhan.co/v2/charts/intraday"

    var payload Payload

    payload = Payload{ 
        SecurityID: "1333",
        ExchangeSegment: "NSE_EQ",
        Instrument: "EQUITY",
        Interval: "1",
        Oi: false,
        FromDate: "2024-09-11 09:30:00",
        ToDate: "2024-09-15 13:00:00",
    }

    var jsonData []byte
    var err error

    jsonData, err = json.Marshal(payload)

    if err != nil {
        fmt.Println("Error marshalling JSON:", err)
        return
    }

    var req *http.Request
    req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))

    if err != nil {
        fmt.Println("Error creating HTTP request:", err)
        return
    }
    
    // Set Headers
    req.Header.Set("Accept", "application/json")
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("access-token", "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzUxMiJ9.eyJpc3MiOiJkaGFuIiwicGFydG5lcklkIjoiIiwiZXhwIjoxNzYwODg0NjUyLCJpYXQiOjE3NjA3OTgyNTIsInRva2VuQ29uc3VtZXJUeXBlIjoiU0VMRiIsIndlYmhvb2tVcmwiOiIiLCJkaGFuQ2xpZW50SWQiOiIxMTA4ODcwNTEwIn0.IRPks-Cfbx6ZY1Vp2p7W6hDO6IidollrAav3b3tWSq29wC42u4Wxx10blJ7ZGQ_XDVYgdAL86asJirSkW4UXmg")

    var client *http.Client = &http.Client{}

    // Send the request
    var resp *http.Response
    resp, err = client.Do(req)

    if err != nil {
        fmt.Println("Error sending HTTP request:", err)
        return
    }

    defer resp.Body.Close()

    var body []byte
    body, err = ioutil.ReadAll(resp.Body)

    if err != nil {
        fmt.Println("Error reading response body:", err)
        return
    }

    fmt.Println("Response Status:", resp.Status)
    fmt.Println("Response Body:", string(body))
}
