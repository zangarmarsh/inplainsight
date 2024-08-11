package pages

import (
	"errors"
	"fmt"
	"github.com/zangarmarsh/inplainsight/core/events"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"log"
	"time"

	"github.com/rivo/tview"
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
	err := Navigate("register")
	if err != nil {
		log.Println(err)
	}
}

func Navigate(path string) error {
	history = append(history, path)

	if PageFactories[path] != nil {
		if !inplainsight.InPlainSight.Pages.HasPage(path) {
			page := PageFactories[path].Create()

			inplainsight.InPlainSight.Pages.AddPage(
				page.GetName(),
				page.GetPrimitive(),
				true,
				false,
			)

			inplainsight.InPlainSight.Trigger(
				events.Event{
					CreatedAt: time.Now(),
					EventType: events.Navigation,
					Data: map[string]interface{}{
						"slug": path,
					},
				})
		}

		inplainsight.InPlainSight.Pages.SwitchToPage(path)
		return nil
	}

	return errors.New(fmt.Sprintf("there's no such page called `%s`", path))
}

func GoBack() error {
	page := history[len(history)-1]

	err := Navigate(page)

	return err
}
