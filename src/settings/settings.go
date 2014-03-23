package settings

import (
	"os"
	"encoding/gob"
	"fmt"
	"bufio"
	"bytes"
	"io/ioutil"
	"strings"
)

type Config struct {
	Username, Apikey, Apisecret, BTCToHashTrade, LTCToBTCTrade, NMCToHashTrade, NMCToHashTradeFancy string
}

/*	CEX.IO keys */
	var cexUserName string 
	var cexAPIKey string
	var cexAPISecret string
	var Settings Config

func ReadSettings () (string, string, string, string, string, string, string){

	inputFile, inputError := ioutil.ReadFile("settings.cfg")
	if inputError != nil {
		fmt.Printf("Error opening input file.\n")
		inputReader := bufio.NewReader(os.Stdin)
		fmt.Printf("Please enter your CEX.IO user name :\n")
		cexUserName, _ := inputReader.ReadString('\n')
		fmt.Printf("Please enter your CEX.IO API key :\n")
		cexAPIKey, _ := inputReader.ReadString('\n')
		fmt.Printf("Please enter your CEX.IO Secret key :\n")
		cexAPISecret, _ := inputReader.ReadString('\n')
		fmt.Printf("Trade BTC to GHS (Y/N)?\n")
		bth, _ := inputReader.ReadString('\n')
		fmt.Printf("Trade LTC to BTC (Y/N)?\n")
		lth, _ := inputReader.ReadString('\n')
		fmt.Printf("Trade NMC to GHS (Y/N)?\n")
		ntg, _ := inputReader.ReadString('\n')
		fmt.Printf("Perform NMC to GHS (via BTC) (Y/N)?\n")
		fnmc, _ := inputReader.ReadString('\n')
		
		var Settings = Config{cexUserName, cexAPIKey, cexAPISecret, bth, lth, ntg, fnmc}
	//	fmt.Printf("GOB Values : %+v\n", Settings)
		b := new(bytes.Buffer)
		
		enc := gob.NewEncoder(b)
		err := enc.Encode(Settings)
		if err != nil {
                fmt.Printf("Encode : %v\n", err)
        }
		eopen := ioutil.WriteFile("settings.cfg", b.Bytes(), 0666)
		if eopen != nil {
			fmt.Printf("File Open Error : %v\n", err)
        }
		return strings.TrimSuffix(Settings.Username, "\r\n"), strings.TrimSuffix(Settings.Apikey, "\r\n"), strings.TrimSuffix(Settings.Apisecret, "\r\n"), strings.TrimSuffix(Settings.BTCToHashTrade, "\r\n"), strings.TrimSuffix(Settings.LTCToBTCTrade, "\r\n"), strings.TrimSuffix(Settings.NMCToHashTrade, "\r\n"), strings.TrimSuffix(Settings.NMCToHashTradeFancy, "\r\n")

	} else {
		b := bytes.NewBuffer(inputFile)

		dec :=gob.NewDecoder(b)
		err := dec.Decode(&Settings)
		if err != nil {
		fmt.Printf("GOB : %v\n", err)
		fmt.Printf("GOB Values : %v\n", Settings)
		}
//		fmt.Printf("GOB File Load Values : %q\n", Settings)
		}

//		fmt.Printf("U? %s", cexUserName)
//		fmt.Printf("WTF GOB Values : %+v\n", Settings)

	return strings.TrimSuffix(Settings.Username, "\r\n"), strings.TrimSuffix(Settings.Apikey, "\r\n"), strings.TrimSuffix(Settings.Apisecret, "\r\n"), strings.TrimSuffix(Settings.BTCToHashTrade, "\r\n"), strings.TrimSuffix(Settings.LTCToBTCTrade, "\r\n"), strings.TrimSuffix(Settings.NMCToHashTrade, "\r\n"), strings.TrimSuffix(Settings.NMCToHashTradeFancy, "\r\n")
}
