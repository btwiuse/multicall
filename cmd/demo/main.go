package main

import (
	"log"
	"os"

	"github.com/btwiuse/multicall"
	"k0s.io/pkg/cli/agent"
	"k0s.io/pkg/cli/chassis"
	"k0s.io/pkg/cli/client"
	"k0s.io/pkg/cli/hub"
	"k0s.io/pkg/cli/miniclient"
	"k0s.io/pkg/cli/mnt"
	"k0s.io/pkg/cli/upgrade"
)

var cmdRun multicall.RunnerFuncMap = map[string]multicall.RunnerFunc{
	"mnt":        mnt.Run,
	"chassis":    chassis.Run,
	"client":     client.Run,
	"miniclient": miniclient.Run,
	"hub":        hub.Run,
	"hub2":       hub.Run2,
	"agent":      agent.Run,
	"upgrade":    upgrade.Run,
}

func main() {
	log.Fatalln(cmdRun.Run(os.Args))
}
