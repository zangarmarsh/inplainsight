package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/ui"
)

func ModalSuccess(text string) tview.Primitive {
	modal := tview.NewModal()

	modal.SetText(text).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {

		})

	return modal
}

func ModalError(text string) tview.Primitive {
	modal := tview.NewModal()

	modal.SetText( "Wrong password" ).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetBackgroundColor(tcell.ColorRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			ui.InPlainSight.Pages.HidePage("modal-error")
		})

	return modal
}

