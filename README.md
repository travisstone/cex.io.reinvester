cex.io.reinvester
=================

Send BTC Donations to : 15D6r1KdA7J3sLZAfX9Zy1aPfQD7V6E7fm

Robot to Buy GHS from BTC/NMC/LTC earnings

Written in Go.
Program will prompt for API and Trade settings on first run.
You will need to provide your Cex.IO username, API key, and API Secret key.
Answer wiht a "Y" or a "N" for the Trade settings.

Trade settings include options to trade BTC for GHS, LTC for BTC, NMC to GHS, and NMC to GHS (via BTC).
LTC for BTC will check BTC-E to make sure that the market is with in tolerance. Adjsut the global varible "ltcExchange" if you want a different tolerance level.
NMC to GHS (via BTC) trading requires that you enable NMC to GHS. It checks the math on in you receive more GHS if you convert to BTC first. You will also need to enable BTC to GHS (If you do not, it won't execut the BTC to GHS buy order if NMC was converted to BTC). 

There are Thresholds for trading. If you want to maintain a balance you will need to adjust the following global values: BTCthres, LTCthres, NMCthres. Threshold Values should be set to >= 0.0000001 to avoid issues with floating point math.

Added resell func. If variable BTCHashResell = "Y", it should post a sell order at the markup you specify with variable ResellMarkup.
