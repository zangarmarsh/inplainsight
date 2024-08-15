package widgets

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"github.com/zangarmarsh/inplainsight/ui/pages"
	"strconv"
)

type ModalType int

const (
	ModalInfo ModalType = iota
	ModalAlert
	ModalSuccess
)

func NewModal(modalType ModalType, text string, okButtonTxt string, callback func()) *pages.GridPage {
	pageName := "modal-" + strconv.Itoa(int(modalType))

	page := pages.GridPage{}
	page.SetName(pageName)

	modalError := tview.NewModal()

	var buttons []string
	{
		if okButtonTxt != "" {
			buttons = append(buttons, "Cancel")
			buttons = append(buttons, okButtonTxt)
		} else {
			buttons = append(buttons, "Close")
		}
	}
	modalError.AddButtons(buttons)

	var backgroundColor tcell.Color
	switch modalType {
	case ModalAlert:
		backgroundColor = tcell.ColorRed
	case ModalSuccess:
		backgroundColor = tcell.ColorGreen
	case ModalInfo:
		backgroundColor = tcell.ColorBlue
	}

	modalError.
		SetText(text).
		SetBackgroundColor(backgroundColor).
		SetBorder(false).
		SetBorderPadding(0, 0, 0, 0)

	modalError.
		SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			if okButtonTxt != "" && buttonLabel == okButtonTxt && callback != nil {
				callback()
			}
			pages.GoBack()
			inplainsight.InPlainSight.Pages.RemovePage(pageName)
		})

	page.SetPrimitive(modalError)
	pages.Navigate(&page)

	return &page
}
