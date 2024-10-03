package tuihelpers

import (
	"github.com/rivo/tview"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"path/filepath"
)

func GenerateContainerSelector(containerPath string) *tview.InputField {
	dedicatedHostInput := tview.NewInputField()
	dedicatedHostInput.
		SetLabel("Container")

	if containerPath != "" {
		dedicatedHostInput.SetText(containerPath)
	}

	dedicatedHostInput.
		SetAutocompleteFunc(func(currentText string) (entries []string) {
			if results := inplainsight.InPlainSight.Hosts.SearchByContainerPath(currentText); results != nil {
				for _, result := range results {
					// Todo check if secret can be contained by chosen container
					entries = append(entries, filepath.Base(result.Host.GetPath()))
				}
			}

			if len(entries) > 1 {
				entries = append([]string{"Random"}, entries...)
			}

			return
		})

	return dedicatedHostInput
}
