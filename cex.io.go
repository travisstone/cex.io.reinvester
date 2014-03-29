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
	"stats"
	"settings"
	)

/*	CEX.IO keys */
	var Username = ""
	var Apikey = ""
	var Apisecret = ""

/*	record highs */
	var BTChigh float64 = 0
	var LTChigh float64 = 0
	var NMChigh float64 = 0
	var GHShigh float64 = 0
	var AssetHigh float64 = 0
	var CEXAPICallsMadeHigh float64 = 0
	var CEXAPICallsMade float64 = 0
	
	var BTCToHashTrade = ""
	var LTCToBTCTrade = ""
	var NMCToHashTrade = ""
	var NMCToHashTradeFancy = ""
	var BTCHashResell = "Y"
	var resellAmount = ""
	var resellPrice = ""
	
	var BTCthres = "0.0000001"
	var LTCthres = "0.0000001"
	var NMCthres = "0.0000001"
	var NMCSoldForBTC bool
	var ltcExchange float64 = 0.005
	var ResellMarkup = "1.01"
	
	var nonce = ""
	
func signatureCalc () string {
	message := nonce + Username + Apikey
	key := []byte(Apisecret)
	h := hmac.New(sha256.New, key)
    h.Write([]byte(message))
	signature := strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
	return signature
	
}
func getBalance () (string, string, string, string, string) {
	nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
	sig := signatureCalc ()

	v:= url.Values {}
	v.Set("key", Apikey)
	v.Add("signature", sig)
	v.Add("nonce", nonce)
	
//	fmt.Printf("URL Values :\n%q\n", v)
	
	BTCBalance, err := http.PostForm("https://cex.io/api/balance/", v )
   	if err != nil {log.Fatal(err)}
   	BalanceData, err := ioutil.ReadAll(BTCBalance.Body)
   	BTCBalance.Body.Close()
   	if err != nil {
		fmt.Printf("Balance : %s\n", BalanceData)
	}
	CEXAPICallsMade = CEXAPICallsMade + 1	
	var balance interface {}
//	fmt.Printf("Balance : %s\n", BalanceData)
	
	balanceerr := json.Unmarshal(BalanceData, &balance)
	if balanceerr != nil {
		fmt.Printf("Balance Retrieval Error : %v\n", balanceerr)
		}
	var balanceLTCblob map[string] interface{}
	balanceBlob := balance.(map[string]interface{})
	
	balanceBTCblob := balanceBlob["BTC"].(map[string]interface{})
	if balanceBlob["LTC"] != nil {
		balanceLTCblob = balanceBlob["LTC"].(map[string]interface{})
		} else {
		balanceLTCblob := map[string]interface{}{
		"available":0,
		"orders":0,
		}
		fmt.Printf("Balance Blob LTC Trouble : %q\n", balanceLTCblob)
	}
	
	balanceNMCblob := balanceBlob["NMC"].(map[string]interface{})
	balanceGHSblob := balanceBlob["GHS"].(map[string]interface{})
	
	btcCalcAval, _ := strconv.ParseFloat(fmt.Sprintf("%s", balanceBTCblob["available"]), 64)
	ltcCalcAval, _ := strconv.ParseFloat(fmt.Sprintf("%s", balanceLTCblob["available"]), 64)
	nmcCalcAval, _ := strconv.ParseFloat(fmt.Sprintf("%s", balanceNMCblob["available"]), 64)
	ghsCalcAval, _ := strconv.ParseFloat(fmt.Sprintf("%s", balanceGHSblob["available"]), 64)
	btcCalcOrd, _ := strconv.ParseFloat(fmt.Sprintf("%s", balanceBTCblob["orders"]), 64)
	ltcCalcOrd, _ := strconv.ParseFloat(fmt.Sprintf("%s", balanceLTCblob["orders"]), 64)
	nmcCalcOrd, _ := strconv.ParseFloat(fmt.Sprintf("%s", balanceNMCblob["orders"]), 64)
	ghsCalcOrd, _ := strconv.ParseFloat(fmt.Sprintf("%s", balanceGHSblob["orders"]), 64)
	
	btcCalc := btcCalcAval + btcCalcOrd
	ltcCalc := ltcCalcAval + ltcCalcOrd
	nmcCalc := nmcCalcAval + nmcCalcOrd
	ghsCalc := ghsCalcAval + ghsCalcOrd

	if btcCalc > BTChigh {
		BTChigh = btcCalc
		}
	if ltcCalc > LTChigh {
		LTChigh = ltcCalc
		}
	if nmcCalc > NMChigh {
		NMChigh = nmcCalc
		}
	if ghsCalc > GHShigh {
		GHShigh = ghsCalc
		}

	return fmt.Sprint(balanceBTCblob["available"]), fmt.Sprint(balanceLTCblob["available"]), fmt.Sprint(balanceNMCblob["available"]), fmt.Sprint(ghsCalc), fmt.Sprint(balanceGHSblob["available"])
	 
 }
 
func BTCHashBuy (btc string, btcAsk string) {
	if btc > BTCthres {
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
		buyValues.Set("key", Apikey)
		buyValues.Add("signature", sig)
		buyValues.Add("nonce", nonce)
		buyValues.Add("type", "buy")
		buyValues.Add("amount", buyAmount)
		buyValues.Add("price", btcAsk)
		

		BTCbuy, err := http.PostForm("https://cex.io/api/place_order/GHS/BTC", buyValues )
    	if err != nil {log.Fatal(err)}
    	BTCbuydata, err := ioutil.ReadAll(BTCbuy.Body)
    	BTCbuy.Body.Close()
    	if err != nil {fmt.Printf("BTCbuydata : %s\n", BTCbuydata)}
		CEXAPICallsMade = CEXAPICallsMade + 1

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
		if (BTCHashResell == "Y") {
			go resell (buyAmount, btcAsk)
			}
		}
	
}

func LTCHashBuy (ltc string) {
	if ltc > LTCthres {
		fmt.Printf("Starting LTC Module: %v\n", ltc)
		time.Sleep (1 * 1e9)

		nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
//		sig := signatureCalc ()

		LTCticker, err := http.PostForm("https://cex.io/api/ticker/LTC/BTC", nil )
    	if err != nil {log.Fatal(err)}
    	LTCtickerdata, err := ioutil.ReadAll(LTCticker.Body)
    	LTCticker.Body.Close()
    	if err != nil {fmt.Printf("LTCtickerdata Err: %s\n", LTCtickerdata)}

		CEXAPICallsMade = CEXAPICallsMade + 1
		
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
			sellValues.Set("key", Apikey)
			sellValues.Add("signature", sig)
			sellValues.Add("nonce", nonce)
			sellValues.Add("type", "sell")
			sellValues.Add("amount", sellAmount)
			sellValues.Add("price", ltcbtcask)
		

			LTCSell, err := http.PostForm("https://cex.io/api/place_order/LTC/BTC", sellValues )
			if err != nil {log.Fatal(err)}
			LTCSelldata, err := ioutil.ReadAll(LTCSell.Body)
			LTCSell.Body.Close()
			if err != nil {fmt.Printf("LTCSelldata Err: %s\n", LTCSelldata)}

			CEXAPICallsMade = CEXAPICallsMade + 1
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

func NMCHashBuy (nmc string, btcAsk string) bool {
	if nmc > NMCthres {
		fmt.Printf("Starting NMC Module : %v\n", nmc)
		fmt.Printf("Checking Trade Math\n")
		time.Sleep (10 * 1e9)
		nonce = strconv.FormatInt(time.Now().UnixNano(), 10)

		NMCticker, err := http.PostForm("https://cex.io/api/ticker/GHS/NMC", nil )
    	if err != nil {
			fmt.Printf("NMC Ticker Post Err : %s\n", err)
			return NMCSoldForBTC
		}
    	NMCtickerdata, err := ioutil.ReadAll(NMCticker.Body)
    	NMCticker.Body.Close()
    	if err != nil {fmt.Printf("NMCtickerdata Err : %s\n", NMCtickerdata)}
		
		CEXAPICallsMade = CEXAPICallsMade + 1
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
		
		CEXAPICallsMade = CEXAPICallsMade + 1
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
		
		if nmcToghs > nmcTobtcTotal ||  NMCToHashTradeFancy == "N" {
			fmt.Printf("Buying Hash with NMC\n")
			nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
			sig := signatureCalc ()
			
			buyAmount := fmt.Sprintf("%.8f", nmcToghs  - 0.00000001)
			
			nmcbuyValues := url.Values {}
			nmcbuyValues.Set("key", Apikey)
			nmcbuyValues.Add("signature", sig)
			nmcbuyValues.Add("nonce", nonce)
			nmcbuyValues.Add("type", "buy")
			nmcbuyValues.Add("amount", buyAmount)
			nmcbuyValues.Add("price", nmcghsask)
	
//			fmt.Printf("NMC to GHS : %q\n",nmcbuyValues)
			
			NMCbuy, err := http.PostForm("https://cex.io/api/place_order/GHS/NMC", nmcbuyValues )
			if err != nil {log.Fatal(err)}
			NMCbuydata, err := ioutil.ReadAll(NMCbuy.Body)
			NMCbuy.Body.Close()
			if err != nil {fmt.Printf("NMCbuydata Err : %s", NMCbuydata)}

//			fmt.Printf("NMC to GHS : %v\n", NMCbuydata)
			CEXAPICallsMade = CEXAPICallsMade + 1


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
			NMCSoldForBTC = false
			return NMCSoldForBTC
			} else {
			fmt.Printf("Converting NMC to BTC\n")
			nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
			sig := signatureCalc ()
			
			buyAmount := fmt.Sprintf("%.8f", nmcFloat - 0.00000001)
			
			nmcbuyValues := url.Values {}
			nmcbuyValues.Set("key", Apikey)
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
			if err != nil {fmt.Printf("NMCselldata Err : %s", NMCselldata)}

			CEXAPICallsMade = CEXAPICallsMade + 1
			
			var NMCSellJSON interface {}
		
			NMCSellbloberr := json.Unmarshal(NMCselldata, &NMCSellJSON)
			if NMCSellbloberr != nil {
				fmt.Printf("NMC Sell JSON Error : %v\n", NMCSellbloberr)
				}	

			fmt.Printf("Placed NMC/BTC Sell order for %v\n", buyAmount)
			NMCSoldForBTC = true
			return NMCSoldForBTC
				
			}
		}
	return NMCSoldForBTC
}

func resell (resellAmount string, resellPrice string) {
	for {
		time.Sleep (10 * 1e9)
		_, _, _, _, ghsBal := getBalance ()
		if ghsBal >= resellAmount {
			nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
			sig := signatureCalc ()

			resellPriceFloat, _ := strconv.ParseFloat(resellPrice, 64)
			ResellMarkupFloat, _ := strconv.ParseFloat(ResellMarkup, 64)	

			resellPriceCalc := resellPriceFloat * ResellMarkupFloat
			resellMarkupPriceCalc := fmt.Sprintf("%.8f", resellPriceCalc)
			
			fmt.Printf("Selling %s GHS at %s \n\n", resellAmount, resellMarkupPriceCalc)
		
			resellValues := url.Values {}
			resellValues.Set("key", Apikey)
			resellValues.Add("signature", sig)
			resellValues.Add("nonce", nonce)
			resellValues.Add("type", "sell")
			resellValues.Add("amount", resellAmount)
			resellValues.Add("price", resellMarkupPriceCalc)
		

			BTCbuy, err := http.PostForm("https://cex.io/api/place_order/GHS/BTC", resellValues )
			if err != nil {log.Fatal(err)}
			BTCbuydata, err := ioutil.ReadAll(BTCbuy.Body)
			BTCbuy.Body.Close()
			if err != nil {fmt.Printf("BTCbuydata : %s\n", BTCbuydata)}
			CEXAPICallsMade = CEXAPICallsMade + 1

			var BTCBuyJSON interface {}
		
			BTCBuybloberr := json.Unmarshal(BTCbuydata, &BTCBuyJSON)
			if BTCBuybloberr != nil {
				fmt.Printf("BTC Buy JSON Error : %v\n", BTCBuybloberr)
			}
			break
		}
	}
}
	
func main () {

	StartTime := time.Now()
	Username, Apikey, Apisecret, BTCToHashTrade, LTCToBTCTrade, NMCToHashTrade, NMCToHashTradeFancy = settings.ReadSettings()

	for {
		RunTime := time.Since(StartTime)
		fmt.Printf("CEX API calls Made : %.0f | Runtime : %v\n", CEXAPICallsMade, RunTime)
		nonce = strconv.FormatInt(time.Now().UnixNano(), 10)
	
		BTCticker, err := http.PostForm("https://cex.io/api/ticker/GHS/BTC", nil )
    	if err != nil {
			fmt.Printf("BTC Get Ticker Err : %v\n", err)
			continue
		}
    	BTCtickerdata, err := ioutil.ReadAll(BTCticker.Body)
    	BTCticker.Body.Close()
    	if err != nil {log.Fatal(err)}
		
		CEXAPICallsMade = CEXAPICallsMade + 1

		var m interface {}

		tickerr := json.Unmarshal(BTCtickerdata, &m)
		if tickerr != nil {
			fmt.Printf("BTC Ticker Error : %v\n", tickerr)
			return
			}

		
		jsonBlob := m.(map[string]interface{})
/*		for key, value := range jsonBlob {
			fmt.Println("Key:", key, "Value:", value)
			}
*/		


		CexBTCHigh := fmt.Sprint(jsonBlob["high"])
		CexBTCBid := fmt.Sprint(jsonBlob["bid"])
		
		statBarCexBTCHigh, _ := strconv.ParseFloat(CexBTCHigh, 64)
		statBarCexBTCBid, _ := strconv.ParseFloat(CexBTCBid, 64)

		btcAsk := fmt.Sprint(jsonBlob["ask"])
		btc, ltc, nmc, ghs, _ := getBalance ()

		
		statBarBTC, _ := strconv.ParseFloat(btc, 64)
		statBarGHS, _ := strconv.ParseFloat(ghs, 64)
		statBarLTC, _ := strconv.ParseFloat(ltc, 64)
		statBarNMC, _ := strconv.ParseFloat(nmc, 64)
		AssetCurrent := statBarGHS *statBarCexBTCBid  
		if AssetCurrent > AssetHigh {
			AssetHigh = AssetCurrent
		}

		fmt.Printf("Cex High        : %.8f | Cex Bid        : %.8f | %v\n", statBarCexBTCHigh, statBarCexBTCBid, stats.StatBar (statBarCexBTCHigh, statBarCexBTCBid))
		fmt.Printf("Cex Asset Value : %.8f | Cex Asset High : %.8v | %v\n", AssetCurrent, AssetHigh, stats.StatBar(AssetHigh, AssetCurrent))
		
		fmt.Printf("Balances:\n")
		fmt.Printf("Current BTC     : %s\t| Highest Balance  : %s\t| %v\n", fmt.Sprintf("%.8f",statBarBTC), fmt.Sprintf("%.8f",BTChigh), stats.StatBar(BTChigh, statBarBTC))
		fmt.Printf("Current GHS     : %s\t| Highest Balance  : %s\t| %v\n", fmt.Sprintf("%.8f",statBarGHS), fmt.Sprintf("%.8f",GHShigh), stats.StatBar(GHShigh, statBarGHS))
		fmt.Printf("Current LTC     : %s\t| Highest Balance  : %s\t| %v\n", fmt.Sprintf("%.8f",statBarLTC), fmt.Sprintf("%.8f",LTChigh), stats.StatBar(LTChigh, statBarLTC))
		fmt.Printf("Curremt NMC     : %s\t| Highest Balance  : %s\t| %v\n", fmt.Sprintf("%.8f",statBarNMC), fmt.Sprintf("%.8f",NMChigh) ,stats.StatBar(NMChigh, statBarNMC))

		if (NMCToHashTrade == "Y") {
			NMCSoldForBTC = NMCHashBuy (nmc, btcAsk)
		}

		if NMCSoldForBTC {
			fmt.Printf("Regetting Balance\n",)
			btc, ltc, nmc, ghs, _ = getBalance ()
			NMCSoldForBTC = false
		}
		
		if (BTCToHashTrade == "Y") {
			BTCHashBuy (btc, btcAsk)
		}
			
		if (LTCToBTCTrade == "Y") {
			LTCHashBuy (ltc)
		}

    	time.Sleep (60 * 1e9)
    	}
}
