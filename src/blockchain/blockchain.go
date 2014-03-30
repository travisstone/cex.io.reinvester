package blockchain

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"time"
//	"strings"
)
	
func TimeTillDiffChange () string {

	getblockcount, err := http.Get("https://blockchain.info/q/getblockcount")
   	if err != nil {
		fmt.Printf("Get Block Count Error : %s\n", err)
		return fmt.Sprintf("Get Block Count Error : %s\n", err)
	}
   	getblockcountdata, err := ioutil.ReadAll(getblockcount.Body)
    getblockcount.Body.Close()
    if err != nil {
		fmt.Printf("Get Block Count Error : %s\n", err)
		return fmt.Sprintf("Get Block Count Error : %s\n", err)
	}

	currentBlock, _ := strconv.ParseFloat(fmt.Sprintf("%s",getblockcountdata), 64)

	progress := currentBlock / 2016
	testprogress := math.Trunc(progress) + 1
	nextChange := testprogress * 2016
	blocksRemaining := nextChange - currentBlock
	
	timeEstimate := blocksRemaining * 360
	timeEstimateDuraion, _ := time.ParseDuration(fmt.Sprintf("%.0fs", timeEstimate))
	return fmt.Sprintf("BTC Retarget Guess : %v\n", time.Now().Add(timeEstimateDuraion))
}
