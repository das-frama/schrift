package main

import (
	"fmt"
	"log"
	"time"

	"github.com/marcusolsson/tui-go"
)

const (
	caret  = "$ "
	myName = "frama"
)

var ui tui.UI

type uiChanel struct {
	label    *tui.Label
	messages []Message
	history  *tui.Box
	cursor   int
}

type historyWidget struct {
	*tui.Box
	historyBox *tui.Box

	chanels []*uiChanel
	active  int
}

func (h *historyWidget) OnKeyEvent(ev tui.KeyEvent) {
	switch ev.Key {
	case tui.KeyTab:
		h.Next()
	case tui.KeyBacktab:
		h.Prev()
	}

	h.Box.OnKeyEvent(ev)
}

func (h *historyWidget) SetChanel(chanel *uiChanel) {
	h.historyBox.Remove(0)
	historyScroll := tui.NewScrollArea(chanel.history)
	historyScroll.SetAutoscrollToBottom(true)
	h.historyBox.Append(historyScroll)
	go h.PrintMessage(chanel.messages[0].text)
}

func (h *historyWidget) Next() {
	h.active = clamp(h.active+1, 0, len(h.chanels)-1)
	h.style()
	h.SetChanel(h.chanels[h.active])
}

func (h *historyWidget) Prev() {
	h.active = clamp(h.active-1, 0, len(h.chanels)-1)
	h.style()
	h.SetChanel(h.chanels[h.active])
}

func clamp(n, min, max int) int {
	if n < min {
		return max
	}
	if n > max {
		return min
	}
	return n
}

func (h *historyWidget) style() {
	for i := 0; i < len(h.chanels); i++ {
		if i == h.active {
			h.chanels[i].label.SetStyleName("chanel-selected")
			continue
		}
		h.chanels[i].label.SetStyleName("chanel")
	}
}

func (h *historyWidget) PrintMessage(message string) {
	textLabel := tui.NewLabel("")
	// textLabel.SetWordWrap(true)
	h.chanels[h.active].history.Append(tui.NewHBox(
		tui.NewLabel(fmt.Sprintf("[%s]: ", myName)),
		textLabel,
		tui.NewSpacer(),
	))
	for _, s := range message {
		textLabel.SetText(textLabel.Text() + string(s))
		ui.Update(func() {})
		time.Sleep(50 * time.Millisecond)
	}
}

func newHistoryWidget(chanels ...*uiChanel) *historyWidget {
	view := &historyWidget{chanels: chanels}

	sidebar := tui.NewVBox()
	sidebar.SetBorder(true)
	sidebar.SetSizePolicy(tui.Maximum, tui.Minimum)
	for i := 0; i < len(chanels); i++ {
		sidebar.Append(chanels[i].label)
	}
	sidebar.Append(tui.NewSpacer())

	historyScroll := tui.NewScrollArea(chanels[0].history)
	historyScroll.SetAutoscrollToBottom(true)
	historyBox := tui.NewVBox(historyScroll)
	historyBox.SetBorder(true)
	view.historyBox = historyBox

	input := tui.NewEntry()
	input.SetFocused(true)
	input.SetSizePolicy(tui.Expanding, tui.Maximum)
	inputBox := tui.NewHBox(
		tui.NewLabel(caret),
		input,
	)
	inputBox.SetBorder(true)
	inputBox.SetSizePolicy(tui.Expanding, tui.Maximum)

	chat := tui.NewVBox(historyBox, inputBox)
	chat.SetSizePolicy(tui.Expanding, tui.Expanding)

	input.OnSubmit(func(e *tui.Entry) {
		view.chanels[view.active].history.Append(tui.NewHBox(
			tui.NewLabel("> "),
			tui.NewLabel(e.Text()),
			tui.NewSpacer(),
		))
		input.SetText("")
	})
	view.Box = tui.NewHBox(sidebar, chat)
	view.style()

	return view
}

func main() {
	historyLayout := newHistoryWidget(
		&uiChanel{label: tui.NewLabel("#0. Пролог"), messages: Intro, history: tui.NewVBox()},
		&uiChanel{label: tui.NewLabel("#1. В доме приглушён огонь"), messages: Test, history: tui.NewVBox()},
		// &uiChanel{label: tui.NewLabel("#2. Пробежавший лось"), messases: &Test, history: tui.NewVBox()},
	)
	var err error
	ui, err = tui.New(historyLayout)
	if err != nil {
		log.Fatal(err)
	}
	ui.SetKeybinding("Esc", func() { ui.Quit() })

	theme := tui.NewTheme()
	theme.SetStyle("label.chanel", tui.Style{Reverse: tui.DecorationOff})
	theme.SetStyle("label.chanel-selected", tui.Style{Reverse: tui.DecorationOn, Fg: tui.ColorRed, Bg: tui.ColorWhite})
	ui.SetTheme(theme)

	go historyLayout.PrintMessage("Запуск!")

	if err := ui.Run(); err != nil {
		log.Fatal(err)
	}
}
