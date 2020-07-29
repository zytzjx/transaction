package main

import (
	"os"

	Log "github.com/zytzjx/anthenacmc/loggersys"
	"github.com/zytzjx/anthenacmc/reportcmc"
)

func main() {
	Log.NewLogger("reportcmc")
	if reportcmc.ReportCMC() != nil {
		os.Exit(1)
	}
	os.Exit(0)
}
