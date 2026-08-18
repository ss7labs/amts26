package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/ss7labs/amts26/config"
	"github.com/ss7labs/amts26/internal/app"
	"github.com/ss7labs/amts26/internal/r11"

	"embed"

	"github.com/ss7labs/amts26/pkg/common"
	"github.com/ss7labs/amts26/pkg/rtprefix"
)

var (
	nolog    = flag.Bool("nolog", false, "logger disable")
	printcfg = flag.Bool("printcfg", false, "print config")
	conffile = flag.String("c", "", "config yaml file")
	period   = flag.String("period", "", "YYYYMM month for report")
)

//go:embed rt-prefixes-2023.csv
var files embed.FS

func main() {
	flag.Parse()
	layout := "200601"

	_, err := time.Parse(layout, *period)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return
	}
	_config := config.LoadConfig(*conffile)

	if *printcfg {
		fmt.Printf("%+v\n", common.ToJson(_config))
		return
	}

	application := app.NewApplication(_config)
	application.Init(*nolog)

	ps := rtprefix.NewPrefixStore(files)
	application.SetPrefixStore(ps)

	r11.Report(*period, application)
}
