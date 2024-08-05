package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
)

func ModalSuccess(text string) tview.Primitive {
	pageName := "modal-success"
	modalSuccess := tview.NewModal()

	modalSuccess.SetText(text).
		AddButtons([]string{"OK"}).
		SetFocus(0).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			inplainsight.InPlainSight.Pages.RemovePage(pageName)
		})

	inplainsight.InPlainSight.Pages.AddAndSwitchToPage(pageName, modalSuccess, true)

	return modalSuccess
}

func ModalAlert(text string, callback func()) tview.Primitive {
	pageName := "modal-error"
	okButtonTxt := "OK"

	modalError := tview.NewModal()
	modalError.AddButtons([]string{"Cancel", okButtonTxt})

	modalError.
		SetText(text).
		SetFocus(0).
		SetBackgroundColor(tcell.ColorRed).
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if buttonLabel == okButtonTxt && callback != nil {
				callback()
			}
			inplainsight.InPlainSight.Pages.RemovePage(pageName)
		})

	inplainsight.InPlainSight.Pages.AddAndSwitchToPage(pageName, modalError, true)

	return modalError
}
