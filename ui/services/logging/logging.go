package logging

import (
	"github.com/rivo/tview"
	"log"
)

type LogLevel uint8

const (
	Info     LogLevel = iota
	Warning
	Critical
	Debug
)

type Log struct {
	text  string
	level LogLevel
}

type LogsBox struct {
	box *tview.TextView
	logs []Log
}

func NewLogsBox(box *tview.TextView) *LogsBox {
	return &LogsBox{box, []Log{}}
}

func (l *LogsBox) AddLine(text string, level LogLevel) {
	var joined string
	l.logs = append(l.logs, Log{
		text, level,
	})

	for _, line := range l.logs {
		var colorPrefix string

		switch line.level {
			case Warning:
				colorPrefix = "[orange:i]"
			case Critical:
				colorPrefix = "[red:i]"
			case Debug:
				colorPrefix = "[white:i]"
			case Info:
			default:
				colorPrefix = "[orange:i]"
		}

		log.Println(line)

		joined = joined + colorPrefix + line.text + "[-]\n"
	}

	l.box.SetText(joined)
	l.box.SetScrollable(true).ScrollTo(len(l.logs) + 1, 0)
}

func (l *LogsBox) AddSeparator() {
	l.AddLine("***", l.GetLatestLogLevel())
}

func (l *LogsBox) GetLatestLogLevel() LogLevel {
	ll := Info

	if len(l.logs) > 0 {
		ll = l.logs[len(l.logs) - 1].level
	}

	return ll
}

func (l *LogsBox) GetPrimitive() tview.Primitive {
	return l.box
}