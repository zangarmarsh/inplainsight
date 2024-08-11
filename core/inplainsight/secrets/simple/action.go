package simple

import (
	"github.com/zangarmarsh/inplainsight/core/events"
	"github.com/zangarmarsh/inplainsight/core/inplainsight"
	"golang.design/x/clipboard"
	"log"
	"time"
)

func (s SimpleSecret) DoAction() {
	log.Println("Copying into clipboard")
	clipboard.Write(clipboard.FmtText, []byte(s.secret))

	inplainsight.InPlainSight.Trigger(events.Event{
		CreatedAt: time.Time{},
		EventType: events.SecretCopiedIntoClipboard,
		Data:      nil,
	})
}
