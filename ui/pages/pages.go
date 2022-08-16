package pages

import (
	"errors"
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui"
	"log"
)


type PageFactoryDictionary map[string]PageFactoryInterface
var PageFactories = make(PageFactoryDictionary)

type PageFactoryInterface interface {
	Create() PageInterface
}

type PageInterface interface {
	GetName() string
	GetPrimitive() tview.Primitive
}

type GridPage struct {
	grid *tview.Grid
	name string
}

func (r *GridPage) GetPrimitive() tview.Primitive {
	return r.grid
}

func (r *GridPage) SetPrimitive(grid *tview.Grid) {
	r.grid = grid
}

func (r *GridPage) GetName() string {
	return r.name
}

func (gp *GridPage) SetName(name string) {
	gp.name = name
}

func Init() {
	log.Print(PageFactories)
	for _, pageFactory := range PageFactories {
		page := pageFactory.Create()

		ui.InPlainSight.Pages.AddPage(
			page.GetName(),
			page.GetPrimitive(),
			true,
			false,
		)
	}

	err := Navigate("login")
	if err != nil {
		log.Println(err)
	}
}

func Navigate(path string) error {
	if PageFactories[path] != nil {
		ui.InPlainSight.Pages.SwitchToPage(path)
		return nil
	}

	return errors.New(fmt.Sprintf("there's no such page called `%s`", path))
}