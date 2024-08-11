package main

import (
	"flag"
	_ "github.com/zangarmarsh/inplainsight/core/events/effects"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/list"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/newsecret"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/register"
	"io"
	"log"
	"os"
)

func main() {
	{
		// handle flags in here
		setLoggingLevel()
	}

	ui.Bootstrap()
	pages.Init()

	err := inplainsight.InPlainSight.App.Run()
	if err != nil {
		log.Fatalln(err)
	}

	inplainsight.InPlainSight.App.
		EnableMouse(true)
}

// ToDo: handle different verbosity levels rather than just one `v`
func setLoggingLevel() {
	var verbosity bool

	flag.BoolVar(&verbosity, "v", false, "allow verbosity")
	flag.Parse()

	if !verbosity {
		log.SetOutput(io.Discard)
	} else {
		log.SetOutput(os.Stdout)
	}
}
