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

type NavigationHistory struct {
	stack []string
}

func (n *NavigationHistory) Push(fragment string) {
	n.stack = append(n.stack, fragment)
}

func (n *NavigationHistory) Pop() (fragment string) {
	lastItemIndex := len(n.stack) - 1
	fragment = n.stack[lastItemIndex]
	n.stack = n.stack[:lastItemIndex]

	return
}

var History = NavigationHistory{}

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

func Navigate(in any) error {
	var page PageInterface
	var pageName string

	switch sType := in.(type) {
	case string:
		if !inplainsight.InPlainSight.Pages.HasPage(in.(string)) {
			if PageFactories[in.(string)] != nil {
				page = PageFactories[in.(string)].Create()
			}
		}

		pageName = in.(string)

	case *GridPage:
		page = (in).(*GridPage)
		pageName = page.GetName()
	default:
		// Todo clean'em up
		log.Fatalln("Unsupported type of page", sType)
		return errors.New("Unsupported type of page")
	}

	if page != nil {
		inplainsight.InPlainSight.Pages.AddPage(
			page.GetName(),
			page.GetPrimitive(),
			true,
			false,
		)
	}

	if inplainsight.InPlainSight.Pages.HasPage(pageName) {
		History.Push(pageName)

		inplainsight.InPlainSight.Trigger(
			events.Event{
				CreatedAt: time.Now(),
				EventType: events.Navigation,
				Data: map[string]interface{}{
					"slug": pageName,
				},
			})

		inplainsight.InPlainSight.Pages.SwitchToPage(pageName)

		if page != nil {
			inplainsight.InPlainSight.App.SetFocus(page.GetPrimitive())
		}

		return nil
	} else {
		return errors.New(fmt.Sprintf("there's no such page called `%v`", in))
	}
}

func GoBack() error {
	page := History.stack[len(History.stack)-2]

	log.Printf("navigating back to %s [%v]", page, History)
	err := Navigate(page)

	return err
}
