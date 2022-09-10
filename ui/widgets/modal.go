package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui"
)

var modalSuccess *tview.Modal
func ModalSuccess(text string) tview.Primitive {
	var newlyCreated bool
	pageName := "modal-success"

	if modalSuccess == nil {
		modalSuccess = tview.NewModal()
		newlyCreated = true
	}

	modalSuccess.SetText(text).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.InPlainSight.Pages.HidePage(pageName)
		})

	if newlyCreated {
		ui.InPlainSight.Pages.AddAndSwitchToPage( pageName, modalSuccess, true )
	} else {
		ui.InPlainSight.Pages.SwitchToPage(pageName)
	}

	return modalSuccess
}

var modalError *tview.Modal
func ModalError(text string) tview.Primitive {
	var newlyCreated bool
	pageName := "modal-error"

	if modalError == nil {
		modalError = tview.NewModal()
		newlyCreated = true
	}

	modalError.SetText( text ).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetBackgroundColor(tcell.ColorRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.InPlainSight.Pages.HidePage(pageName)
		})

	if newlyCreated {
		ui.InPlainSight.Pages.AddAndSwitchToPage( pageName, modalError, true )
	} else {
		ui.InPlainSight.Pages.SwitchToPage(pageName)
	}

	return modalError
}
