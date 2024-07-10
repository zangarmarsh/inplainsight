package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui"
)

func ModalSuccess(text string) tview.Primitive {
	pageName := "modal-success"
	modalSuccess := tview.NewModal()

	modalSuccess.SetText(text).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.InPlainSight.Pages.RemovePage(pageName)
		})

	ui.InPlainSight.Pages.AddAndSwitchToPage(pageName, modalSuccess, true)

	return modalSuccess
}

func ModalError(text string) tview.Primitive {
	pageName := "modal-error"

	modalError := tview.NewModal()
	modalError.AddButtons([]string{"OK"})

	modalError.
		SetText(text).
		SetFocus(0).
		SetBackgroundColor(tcell.ColorRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.InPlainSight.Pages.RemovePage(pageName)
		})

	ui.InPlainSight.Pages.AddAndSwitchToPage(pageName, modalError, true)

	return modalError
}
