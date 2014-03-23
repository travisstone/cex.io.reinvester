package stats

import (
	"fmt"
)

func StatBar (high float64, stat float64) string {
	statusBarValue := stat / high
	statBarRounded := fmt.Sprintf("%.2f", statusBarValue)
	if statBarRounded == "1.00" {
	StatBarGraph := fmt.Sprintf("[==========]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.99" && statBarRounded > "0.9" {
	StatBarGraph := fmt.Sprintf("[========= ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.89" && statBarRounded > "0.8" {
	StatBarGraph := fmt.Sprintf("[========  ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.79" && statBarRounded > "0.7"{
	StatBarGraph := fmt.Sprintf("[=======   ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.69" && statBarRounded > "0.6" {
	StatBarGraph := fmt.Sprintf("[======    ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.59" && statBarRounded > "0.5" {
	StatBarGraph := fmt.Sprintf("[=====     ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.49" && statBarRounded > "0.4" {
	StatBarGraph := fmt.Sprintf("[====      ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.39" && statBarRounded > "0.3" {
	StatBarGraph := fmt.Sprintf("[===       ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.29" && statBarRounded > "0.2" {
	StatBarGraph := fmt.Sprintf("[==        ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.19" && statBarRounded > "0.1" {
	StatBarGraph := fmt.Sprintf("[=         ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else if statBarRounded <= "0.09" && statBarRounded >= "0.0" {
	StatBarGraph := fmt.Sprintf("[          ]")
	return fmt.Sprintf("%s", StatBarGraph)
	} else {
	StatBarGraph := fmt.Sprintf("[  broken  ]")
	return fmt.Sprintf("%s", StatBarGraph)
	}
}
