package main

import (
	"flag"
	"github.com/zangarmarsh/inplainsight/ui"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/list"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/login"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/newsecret"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/register"
	"io"
	"log"
)

var verbosity bool

func main() {
	{
		// handle flags in here
		setLoggingLevel()
	}

	ui.Bootstrap()
	pages.Init()

	err := ui.InPlainSight.App.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

func setLoggingLevel() {
	flag.BoolVar(&verbosity, "v", false, "allow verbosity")
	flag.Parse()

	if !verbosity {
		log.SetOutput(io.Discard)
	}
}

