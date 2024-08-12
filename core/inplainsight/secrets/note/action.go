package note

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"strings"
)

func (s *Note) DoAction() {
	var txtView *tview.TextView
	var commandList *tview.List

	page := &pages.GridPage{}
	page.SetName("show secret")

	grid := tview.NewGrid().
		SetRows(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0).
		SetColumns(0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)

	flex := tview.NewFlex()
	flex.
		SetBorder(true)

	flex.SetDirection(tview.FlexColumn)
	flex.SetBorderPadding(0, 0, 1, 1)

	grid.
		AddItem(flex, 2, 2, 8, 8, 30, 200, false).
		AddItem(flex, 1, 1, 8, 10, 0, 0, false)

	{
		commandList = tview.NewList()
		commandList.
			SetBorder(true).
			SetBorderPadding(2, 2, 2, 2)

		commandList.
			AddItem("Show note", "", 's', func() {
				txtView.SetText(fmt.Sprintf("%s\n\n%s", s.title, s.note))
			}).
			AddItem("Hide note", "", 'h', func() {
				txtView.SetText(fmt.Sprintf("%s\n\n%s", s.title, strings.Repeat("*", len(s.note)%120)))
			}).
			AddItem("Go back to secrets list", "", 'q', func() {
				pages.GoBack()
			})
	}

	{
		text := s.note
		if s.isHidden {
			text = strings.Repeat("*", len(s.note)%120)
		}
		txtView = tview.NewTextView().
			SetText(fmt.Sprintf("%s\n\n%s", s.title, text))

		txtView.
			SetBorder(true).
			SetBorderPadding(2, 2, 4, 2)
	}

	flex.AddItem(commandList, 0, 1, true)
	flex.AddItem(txtView, 0, 1, false)

	page.SetPrimitive(grid)
	pages.Navigate(page)

	inplainsight.InPlainSight.App.SetFocus(commandList)
}
