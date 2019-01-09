package schrift

import (
	"errors"
	"fmt"
	"log"

	"github.com/jroimartin/gocui"
)

// App struct.
type App struct {
	StartMsg string
	YesMsg   string
	NoMsg    string
	Chapters map[string]string
	Delay    int64
	Speed    int64
	Chapter  string

	index    map[string]int
	messages map[string][]*Message
	// userAnswers map[string][]userAnswer
	state   State
	hasInit bool

	gui *gocui.Gui
}

// type userAnswer struct {
// 	Text string
// 	At   time.Time
// }

// NewApp creates application object that contains the whole program.
func NewApp() *App {
	return &App{
		YesMsg: "Yes",
		NoMsg:  "No",
	}
}

func (a *App) Init() error {
	var err error
	// Init gocui.Gui.
	if a.gui, err = gocui.NewGui(gocui.OutputNormal); err != nil {
		return err
	}
	// Setup Gui.
	a.gui.Cursor = true
	a.gui.InputEsc = true

	// Build widgets.
	var managers []gocui.Manager

	var lines []string
	for c := range a.Chapters {
		lines = append(lines, c)
	}
	managers = append(managers, NewSideWidget("side", 0, 0, 0.2, 1.0, lines))
	for chapter := range a.Chapters {
		managers = append(managers, NewMainWidget("main_"+chapter, 0.2, 0, 1.0, 0.92))
		a.gui.SetViewOnBottom("main_" + chapter)
	}
	managers = append(managers, NewCmdWidget("cmd", 0.2, 0.92, 1.0))
	managers = append(managers, NewPromptWidget("prompt", 0.5, 0.5, a.StartMsg))

	a.gui.SetManager(managers...)
	a.gui.SetCurrentView("side")

	if err := a.parseAllChapters(); err != nil {
		return err
	}

	a.hasInit = true

	return nil
}

// Run will run the app untill error or exit command are commit.
func (a *App) Run() error {
	defer a.gui.Close()

	if !a.hasInit {
		return errors.New("you must have call Init() before Run()")
	}

	if err := a.keybindings(); err != nil {
		return err
	}

	if err := a.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func (a *App) parseAllChapters() error {
	var err error
	for chapter, filename := range a.Chapters {
		if a.messages[chapter], err = parseFile(filename); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) switchChapter(chapter string, run bool) {
	_, ok := a.Chapters[chapter]
	if !ok {
		return
	}
	a.Chapter = chapter
	a.gui.SetCurrentView("main_" + chapter)
	a.gui.SetViewOnTop("main_" + chapter)

	if run {
		a.runMessages(a.messages[chapter], 0, a.gui)
	}
}

func (a *App) runMessages(mm []*Message, index int, g *gocui.Gui) {
	g.Cursor = false
	main, err := g.SetCurrentView("main_" + a.Chapter)
	if err != nil {
		log.Fatalln(err)
	}

	a.state = RunningState
	for i, m := range mm {
		m.Send(g, main)
		a.state = m.State
		a.index[a.Chapter] = i
		if m.State == UserState {
			g.Cursor = true
			g.SetCurrentView("cmd")
			break
		}
	}

	a.state = UserState
	g.Cursor = true
	g.SetCurrentView("cmd")
}

func (a *App) keybindings() error {
	if err := a.gui.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := a.gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, nextFocus); err != nil {
		return err
	}
	if err := a.gui.SetKeybinding("cmd", gocui.KeyEnter, gocui.ModNone, a.send); err != nil {
		return err
	}
	if err := a.gui.SetKeybinding("side", gocui.KeyEnter, gocui.ModNone, a.run); err != nil {
		return err
	}
	if err := a.gui.SetKeybinding("side", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}
	if err := a.gui.SetKeybinding("side", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	return nil
}

func (a *App) run(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	l, err := v.Line(cy)
	if err != nil {
		return err
	}

	go a.switchChapter(l, true)

	return nil
}

func (a *App) send(g *gocui.Gui, v *gocui.View) error {
	if a.state != UserState {
		return nil
	}

	main, err := g.View("main_" + a.Chapter)
	if err != nil {
		return err
	}

	input := v.ViewBuffer()
	if input != "" {
		fmt.Fprintf(main, "> %s", input)
		v.Clear()
		v.SetCursor(0, 0)
	}

	return nil
}
