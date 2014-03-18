package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"log"
	"time"
	"strconv"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"btce"
	)

/*	CEX.IO keys */
	var username = ""
	var apikey = ""
	var apisecret = ""

	var btcToHashTrade = "1"
	var ltcToBTCTrade = "1"
	var nmcToHashTrade = "1"
	var btcthres = "0.0000001"
	var ltcthres = "0.0000001"
	var nmcthres = "0.0000001"

	var ltcExchange float64 = 0.005
	
	var nonce = ""
	
			
	
	
func signatureCalc () string {
	message := nonce + username + apikey
//	fmt.Printf("Message: %v\n", message)
	key := []byte(apisecret)
	h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
	signature := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
//	fmt.Printf("Signature Calc: %v\n", signature )
	return signature
	
}
func getBalance () (string, string, string) {
	sig := signatureCalc ()
	
//	fmt.Printf("Sig: %v\n", sig)

	v:= url.Values {}
	v.Set("key", apikey)
	v.Add("signature", sig)
	v.Add("nonce", nonce)

	BTCBalance, err := http.PostForm("https://cex.io/api/balance/", v )
   	if err != nil {log.Fatal(err)}
   	BalanceData, err := ioutil.ReadAll(BTCBalance.Body)
   	BTCBalance.Body.Close()
   	if err != nil {log.Fatal(err)}
//	fmt.Printf("Balance : %s\n", BalanceData)
 
	var balance interface {}
	
	balanceerr := json.Unmarshal(BalanceData, &balance)
	if balanceerr != nil {
		fmt.Printf("Balance Retrieval Error : %v\n", balanceerr)
		}
	
	balanceBlob := balance.(map[string]interface{})
	balanceBTCblob := balanceBlob["BTC"].(map[string]interface{})
	balanceLTCblob := balanceBlob["LTC"].(map[string]interface{})
	balanceNMCblob := balanceBlob["NMC"].(map[string]interface{})

	return fmt.Sprint(balanceBTCblob["available"]), fmt.Sprint(balanceLTCblob["available"]), fmt.Sprint(balanceNMCblob["available"])
	 
 }
 
func BTCHashBuy (btc string, btcAsk string) {
	if btc > btcthres {
		fmt.Printf("Starting BTC Module : %v\n", btc)
		fmt.Printf("Buying Hash with BTC\n")
		time.Sleep (1 * 1e9)

		nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
		sig := signatureCalc ()
	
//		fmt.Printf("Sig: %v\n", sig)
		
		btcFloat, _ := strconv.ParseFloat(btc, 64)
		btcAskFloat, _ := strconv.ParseFloat(btcAsk, 64)
		buyAmountFloat := btcFloat/btcAskFloat
		buyAmount := fmt.Sprintf("%.8f", buyAmountFloat - 0.00000001)
		fmt.Printf("Buying : %s\n\n", buyAmount)
		
		buyValues := url.Values {}
		buyValues.Set("key", apikey)
		buyValues.Add("signature", sig)
		buyValues.Add("nonce", nonce)
		buyValues.Add("type", "buy")
		buyValues.Add("amount", buyAmount)
		buyValues.Add("price", btcAsk)
		

		BTCbuy, err := http.PostForm("https://cex.io/api/place_order/GHS/BTC", buyValues )
    	if err != nil {log.Fatal(err)}
    	BTCbuydata, err := ioutil.ReadAll(BTCbuy.Body)
    	BTCbuy.Body.Close()
    	if err != nil {log.Fatal(err)}

		var BTCBuyJSON interface {}
		
		BTCBuybloberr := json.Unmarshal(BTCbuydata, &BTCBuyJSON)
		if BTCBuybloberr != nil {
			fmt.Printf("BTC Buy JSON Error : %v\n", BTCBuybloberr)
			}
/*
		BTCBuyblob := BTCBuyJSON.(map[string]interface{})
		for key, value := range BTCBuyblob {
			fmt.Println("Key:", key, "Value:", value)
			}
*/
		}
	
}

func LTCHashBuy (ltc string) {
	if ltc > ltcthres {
		fmt.Printf("Starting LTC Module: %v\n", ltc)
		time.Sleep (1 * 1e9)

		nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
//		sig := signatureCalc ()

		LTCticker, err := http.PostForm("https://cex.io/api/ticker/LTC/BTC", nil )
    	if err != nil {log.Fatal(err)}
    	LTCtickerdata, err := ioutil.ReadAll(LTCticker.Body)
    	LTCticker.Body.Close()
    	if err != nil {log.Fatal(err)}
		
		var ltcTickerJSON interface {}

		ltcTickerbloberr := json.Unmarshal(LTCtickerdata, &ltcTickerJSON)
		if ltcTickerbloberr != nil {
			fmt.Printf("NMC/GHS Ticker JSON Error : %v\n", ltcTickerbloberr)
			}

		ltcTickerBlob := ltcTickerJSON.(map[string]interface{})
/*		for key, value := range ltcTickerBlob {
			fmt.Println("LTC/BTC Key:", key, "Value:", value)
			}
*/
		ltcbtcask := fmt.Sprint(ltcTickerBlob["bid"])

		fmt.Printf("Cex ask price : %v\n", ltcbtcask)
		btceltcAsk := btce.Balance()
		fmt.Printf("BTCE ask price: %v\n", btceltcAsk)
		btceltcAskFloat, _ := strconv.ParseFloat(btceltcAsk, 64)
		ltcbtcaskFloat, _ := strconv.ParseFloat(ltcbtcask, 64)

		ltcExcVar := btceltcAskFloat / ltcbtcaskFloat
		
		ltcExchigh := 1.0 + ltcExchange
		
		if ltcExcVar > ltcExchigh {
			fmt.Printf("Exchange Var : %v\n", ltcExcVar)
			fmt.Printf("Price at BTCE exceeds Threshold\n")
			fmt.Print("\x07")
			time.Sleep (1 * 1e9)
			fmt.Print("\x07")
			time.Sleep (1 * 1e9)
			fmt.Print("\x07")
			}
		if ltcExcVar <= ltcExchigh{
			fmt.Printf("Exchange Var : %v\n", ltcExcVar)
			fmt.Printf("Cex Price is Acceptable. Initiating Sale\n")

			nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
			sig := signatureCalc ()

			ltcFloat, _ := strconv.ParseFloat(ltc, 64)
			sellAmount := fmt.Sprintf("%.8f", ltcFloat - 0.00000001)
			
			sellValues := url.Values {}
			sellValues.Set("key", apikey)
			sellValues.Add("signature", sig)
			sellValues.Add("nonce", nonce)
			sellValues.Add("type", "sell")
			sellValues.Add("amount", sellAmount)
			sellValues.Add("price", ltcbtcask)
		

			LTCSell, err := http.PostForm("https://cex.io/api/place_order/LTC/BTC", sellValues )
			if err != nil {log.Fatal(err)}
			LTCSelldata, err := ioutil.ReadAll(LTCSell.Body)
			LTCSell.Body.Close()
			if err != nil {log.Fatal(err)}

			var LTCSellJSON interface {}

			LTCSellbloberr := json.Unmarshal(LTCSelldata, &LTCSellJSON)
			if LTCSellbloberr != nil {
				fmt.Printf("LTC Sell JSON Error : %v\n", LTCSellbloberr)
				}	

/*			NMCSellblob := NMCSellJSON.(map[string]interface{})
			for key, value := range NMCSellblob {
				fmt.Println("Key:", key, "Value:", value)
				}
*/
			fmt.Printf("Placed LTC/BTC Sell order for %v\n", sellAmount)
			}
		}
}

func NMCHashBuy (nmc string, btcAsk string) {
	if nmc > nmcthres {
		fmt.Printf("Starting NMC Module : %v\n", nmc)
		fmt.Printf("Checking Trade Math\n")
		time.Sleep (10 * 1e9)
		nonce = strconv.FormatInt(time.Now().UnixNano(), 10)

		NMCticker, err := http.PostForm("https://cex.io/api/ticker/GHS/NMC", nil )
    	if err != nil {log.Fatal(err)}
    	NMCtickerdata, err := ioutil.ReadAll(NMCticker.Body)
    	NMCticker.Body.Close()
    	if err != nil {log.Fatal(err)}
		
		var nmcTickerJSON interface {}

		nmcTickerbloberr := json.Unmarshal(NMCtickerdata, &nmcTickerJSON)
		if nmcTickerbloberr != nil {
			fmt.Printf("NMC/GHS Ticker JSON Error : %v\n", nmcTickerbloberr)
			}

		
		nmcTickerBlob := nmcTickerJSON.(map[string]interface{})
/*		for key, value := range nmcTickerBlob {
			fmt.Println("NMC/GHS Key:", key, "Value:", value)
			}
*/
		nmcghsask := fmt.Sprint(nmcTickerBlob["ask"])
		fmt.Printf("NMCToGHS Ask : %s\n", nmcghsask)
		time.Sleep (5 * 1e9)
		
		NMCBTCticker, err := http.PostForm("https://cex.io/api/ticker/NMC/BTC", nil )
    	if err != nil {log.Fatal(err)}
    	NMCBTCtickerdata, err := ioutil.ReadAll(NMCBTCticker.Body)
    	NMCBTCticker.Body.Close()
    	if err != nil {log.Fatal(err)}
		
		var nmcBTCTickerJSON interface {}

		nmcBTCTickerbloberr := json.Unmarshal(NMCBTCtickerdata, &nmcBTCTickerJSON)
		if nmcBTCTickerbloberr != nil {
			fmt.Printf("NMC/BTC Ticker JSON Error : %v\n", nmcBTCTickerbloberr)
			}

		
		nmcBTCTickerBlob := nmcBTCTickerJSON.(map[string]interface{})
/*		for key, value := range nmcBTCTickerBlob {
			fmt.Println("NMC/BTC Key:", key, "Value:", value)
			}
*/		
		nmcbtcask := fmt.Sprint(nmcBTCTickerBlob["ask"])
		nmcbtcbid := fmt.Sprint(nmcBTCTickerBlob["bid"])
		fmt.Printf("NMCToBTC Ask : %s\n", nmcbtcask)
		fmt.Printf("NMCToBTC Bid : %s\n", nmcbtcbid)

		
		nmcFloat, _ := strconv.ParseFloat(nmc, 64)	
		nmcbtcbidFloat, _ := strconv.ParseFloat(nmcbtcbid, 64)			
		btcAskFloat, _ := strconv.ParseFloat(btcAsk, 64)	
		nmcghsaskFloat, _ := strconv.ParseFloat(nmcghsask, 64)	
		
		nmcFloatStr := fmt.Sprintf("%.8f", nmcFloat)
		nmcbtcbidFloatStr := fmt.Sprintf("%.8f", nmcbtcbidFloat)
		nmcghsaskFloatStr := fmt.Sprintf("%.8f", nmcghsaskFloat)
		
		nmcFloat, _ = strconv.ParseFloat(nmcFloatStr, 64)	
		nmcbtcbidFloat, _ = strconv.ParseFloat(nmcbtcbidFloatStr, 64)			
		nmcghsaskFloat, _ = strconv.ParseFloat(nmcghsaskFloatStr, 64)	

		
		nmcTobtc := nmcFloat * nmcbtcbidFloat
		nmcTobtcTotal := nmcTobtc / btcAskFloat
		nmcToghs := nmcFloat/nmcghsaskFloat

		fmt.Printf("nmcTobtc? : %f\n", nmcTobtc)
		fmt.Printf("GHS if converted? : %v\n", nmcTobtcTotal)
		fmt.Printf("GHS if bought via NMC? : %v\n", nmcToghs)
		
		if nmcToghs > nmcTobtcTotal {
			fmt.Printf("Buying Hash with NMC\n")
			nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
			sig := signatureCalc ()
			
			buyAmount := fmt.Sprintf("%.8f", nmcToghs  - 0.00000001)
			
			nmcbuyValues := url.Values {}
			nmcbuyValues.Set("key", apikey)
			nmcbuyValues.Add("signature", sig)
			nmcbuyValues.Add("nonce", nonce)
			nmcbuyValues.Add("type", "buy")
			nmcbuyValues.Add("amount", buyAmount)
			nmcbuyValues.Add("price", nmcghsask)
	
			fmt.Printf("NMC to GHS : %v\n",nmcbuyValues)
			
			NMCbuy, err := http.PostForm("https://cex.io/api/place_order/GHS/NMC", nmcbuyValues )
			if err != nil {log.Fatal(err)}
			NMCbuydata, err := ioutil.ReadAll(NMCbuy.Body)
			NMCbuy.Body.Close()
			if err != nil {log.Fatal(err)}

			fmt.Printf("%s", NMCbuydata)

			var NMCBuyJSON interface {}
		
			NMCBuybloberr := json.Unmarshal(NMCbuydata, &NMCBuyJSON)
			if NMCBuybloberr != nil {
				fmt.Printf("NMC Buy JSON Error : %v\n", NMCBuybloberr)
				}	
				
/*			NMCBuyblob := NMCBuyJSON.(map[string]interface{})
			for key, value := range NMCBuyblob {
				fmt.Println("Key:", key, "Value:", value)
				}
*/
				fmt.Sprintf("Placed NMC/GHS Buy order for %v\n", buyAmount)
				
			} else {
			fmt.Printf("Converting NMC to BTC\n")
			nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
			sig := signatureCalc ()
			
			buyAmount := fmt.Sprintf("%.8f", nmcFloat - 0.00000001)
			
			nmcbuyValues := url.Values {}
			nmcbuyValues.Set("key", apikey)
			nmcbuyValues.Add("signature", sig)
			nmcbuyValues.Add("nonce", nonce)
			nmcbuyValues.Add("type", "sell")
			nmcbuyValues.Add("amount", buyAmount)
			nmcbuyValues.Add("price", nmcbtcbid)
	
			fmt.Printf("NMC to BTC : %v\n",nmcbuyValues)
			
			NMCsell, err := http.PostForm("https://cex.io/api/place_order/NMC/BTC", nmcbuyValues )
			if err != nil {log.Fatal(err)}
			NMCselldata, err := ioutil.ReadAll(NMCsell.Body)
			NMCsell.Body.Close()
			if err != nil {log.Fatal(err)}

			fmt.Printf("%s", NMCselldata)
			
			var NMCSellJSON interface {}
		
			NMCSellbloberr := json.Unmarshal(NMCselldata, &NMCSellJSON)
			if NMCSellbloberr != nil {
				fmt.Printf("NMC Sell JSON Error : %v\n", NMCSellbloberr)
				}	

/*			NMCSellblob := NMCSellJSON.(map[string]interface{})
			for key, value := range NMCSellblob {
				fmt.Println("Key:", key, "Value:", value)
				}
*/
			fmt.Printf("Placed NMC/BTC Sell order for %v\n", buyAmount)

				
			}
		}
	
}

	
func main () {


	for {
		nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
//		fmt.Printf("Nonce :%s\n", nonce)
	
		BTCticker, err := http.PostForm("https://cex.io/api/ticker/GHS/BTC", nil )
    	if err != nil {log.Fatal(err)}
    	BTCtickerdata, err := ioutil.ReadAll(BTCticker.Body)
    	BTCticker.Body.Close()
    	if err != nil {log.Fatal(err)}
		
		var m interface {}

		tickerr := json.Unmarshal(BTCtickerdata, &m)
		if tickerr != nil {
			fmt.Printf("Ticker Error : %v\n", tickerr)
			break
			}

		
		jsonBlob := m.(map[string]interface{})
/*		for key, value := range jsonBlob {
			fmt.Println("Key:", key, "Value:", value)
			}
*/		

		fmt.Printf("Ask : %v\n", jsonBlob["ask"])
		fmt.Printf("Bid : %v\n", jsonBlob["bid"])
		
		btcAsk := fmt.Sprint(jsonBlob["ask"])
		btc, ltc, nmc := getBalance ()
		fmt.Printf("Balances:\n")
		fmt.Printf("BTC : %v\n", btc)
		fmt.Printf("LTC : %v\n", ltc)
		fmt.Printf("NMC : %v\n", nmc)

		BTCHashBuy (btc, btcAsk)
		NMCHashBuy (nmc, btcAsk)
		LTCHashBuy (ltc)

    	time.Sleep (60 * 1e9)
    	}
}
