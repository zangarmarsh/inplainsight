package website

import (
	"fmt"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"github.com/zangarmarsh/inplainsight/ui/widgets"
	"golang.design/x/clipboard"
	"os/exec"
)

func (s *WebsiteCredential) DoAction() {
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
		AddItem(flex, 2, 2, 4, 8, 30, 200, false).
		AddItem(flex, 1, 1, 4, 10, 0, 0, false)

	{
		commandList = tview.NewList()
		commandList.
			SetBorder(true).
			SetBorderPadding(2, 2, 2, 2)

		commandList.
			AddItem("Go to website", "", 'g', func() {
				cmd := exec.Command("xdg-open", s.website)
				err := cmd.Run()

				if err != nil {
					widgets.NewModal(widgets.ModalAlert, err.Error(), "", nil)
				}
			}).
			AddItem("Copy username", "", 'u', func() {
				clipboard.Write(clipboard.FmtText, []byte(s.account))
			}).
			AddItem("Copy password", "", 'p', func() {
				clipboard.Write(clipboard.FmtText, []byte(s.password))
			}).
			AddItem("Go back to secrets list", "", 'q', func() {
				pages.GoBack()
			})
	}

	{
		txtView = tview.NewTextView().
			SetText(fmt.Sprintf("%s\nu: %s\np: %s\n", s.website, s.account, s.password))

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
