package pages

import (
	"errors"
	"fmt"
	"log"

	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui"
)

type PageFactoryDictionary map[string]PageFactoryInterface

var PageFactories = make(PageFactoryDictionary)
var history []string

type PageFactoryInterface interface {
	Create() PageInterface
}

type PageInterface interface {
	GetName() string
	GetPrimitive() tview.Primitive
}

type GridPage struct {
	grid tview.Primitive
	name string
}

func (r *GridPage) GetPrimitive() tview.Primitive {
	return r.grid
}

func (r *GridPage) SetPrimitive(grid tview.Primitive) {
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

	err := Navigate("register")
	if err != nil {
		log.Println(err)
	}
}

func Navigate(path string) error {
	history = append(history, path)

	if PageFactories[path] != nil {
		ui.InPlainSight.Pages.SwitchToPage(path)
		return nil
	}

	return errors.New(fmt.Sprintf("there's no such page called `%s`", path))
}

func GoBack() error {
	page := history[len(history)-1]

	err := Navigate(page)

	return err
}
