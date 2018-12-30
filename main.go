package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jroimartin/gocui"
)

var headers = []string{
	"Глава 1",
	"Глава 2",
	"Глава 3",
}

const (
	startMessage = "Sollen wir anfangen (Ja / Nein)?"
)

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyEnter, gocui.ModNone, sendCmd); err != nil {
		return err
	}
	if err := g.SetKeybinding("cmd", gocui.KeyTab, gocui.ModNone, run); err != nil {
		return err
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	// Side.
	if v, err := g.SetView("side", 0, 0, int(0.2*float32(maxX))-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, header := range headers {
			fmt.Fprintln(v, header)
		}
	}
	// Main.
	if v, err := g.SetView("main", int(0.2*float32(maxX)), 0, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "> Приветик")
		v.Wrap = true
		v.Autoscroll = true
	}
	// Cmd.
	if v, err := g.SetView("cmd", int(0.2*float32(maxX)), maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		var lineEditor gocui.Editor = gocui.EditorFunc(lineEditor)
		v.Editable = true
		v.Editor = lineEditor
		v.Wrap = false
		if _, err := g.SetCurrentView("cmd"); err != nil {
			return err
		}
	}

	return nil
}

func lineEditor(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	switch {
	case ch != 0 && mod == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyArrowLeft:
		v.MoveCursor(-1, 0, false)
	case key == gocui.KeyArrowRight:
		v.MoveCursor(1, 0, false)
	}
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func sendCmd(g *gocui.Gui, v *gocui.View) error {
	main, err := g.View("main")
	if err != nil {
		return err
	}
	fmt.Fprintf(main, "%s", v.ViewBuffer())
	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

func run(g *gocui.Gui, v *gocui.View) error {
	go sendMsg("К моей голове подключили радио.\n", g)

	return nil
}

func sendMsg(message string, g *gocui.Gui) {
	ch := make(chan rune, len(message))

	go func() {
		for {
			c, more := <-ch
			if more {
				g.Update(func(g *gocui.Gui) error {
					main, err := g.View("main")
					if err != nil {
						return err
					}
					fmt.Fprintf(main, "%c", c)

					return nil
				})
			} else {
				return
			}
		}
	}()

	for _, r := range message {
		ch <- r
		time.Sleep(40 * time.Millisecond)
	}
	close(ch)
}

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.InputEsc = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
