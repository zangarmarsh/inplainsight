package main

import (
	"github.com/zangarmarsh/inplainsight/ui"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/list"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/login"
	_ "github.com/zangarmarsh/inplainsight/ui/pages/register"
	"log"
)

func main() {
	ui.Bootstrap( )
	pages.Init()

	err := ui.InPlainSight.App.Run()
	if err != nil {
		log.Fatalln(err)
	}
}

