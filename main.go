package main

import (
	"flag"
	"fmt"

	"github.com/ss7labs/amts26/config"
	"github.com/ss7labs/amts26/internal/app"
	"github.com/ss7labs/amts26/pkg/common"
)

var (
	printcfg = flag.Bool("printcfg", false, "print config")
	conffile = flag.String("c", "", "config yaml file")
)

func main() {
	flag.Parse()
	_config := config.LoadConfig(*conffile)

	if *printcfg {
		fmt.Printf("%+v\n", common.ToJson(_config))
		return
	}

	application := app.NewApplication(_config)
	application.Init()
}
