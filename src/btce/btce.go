package btce

import (
	"net/http"
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
//	"strings"
)


func Balance() string {

		btce, err := http.Get("https://btc-e.com/api/2/ltc_btc/ticker")
    	if err != nil {log.Fatal(err)}
    	btcedata, err := ioutil.ReadAll(btce.Body)
    	btce.Body.Close()
    	if err != nil {log.Fatal(err)}
		
		var btceTickerJSON interface {}

		btceTickerbloberr := json.Unmarshal(btcedata, &btceTickerJSON)
		if btceTickerbloberr != nil {
			fmt.Printf("BTCE Ticker JSON Error : %v\n", btceTickerbloberr)
			}
			
		TickerJSON := btceTickerJSON.(map[string]interface{})
		btceticker := TickerJSON["ticker"].(map[string]interface{})


//		fmt.Printf("BTCE ask price: %v\n", btceticker["buy"])
		return fmt.Sprint(btceticker["buy"])
}
