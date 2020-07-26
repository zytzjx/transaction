package main

import (
	Log "github.com/zytzjx/anthenacmc/loggersys"
	"github.com/zytzjx/anthenacmc/reportcmc"
)

func main() {
	Log.NewLogger("reportcmc")
	reportcmc.ReportCMC()
}
