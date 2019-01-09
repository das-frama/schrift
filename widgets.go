package schrift

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"
)

type PromptWidget struct {
	name string
	x, y float32
	w, h int
	body string
}

func NewPromptWidget(name string, x, y float32, body string) *PromptWidget {
	lines := strings.Split(body, "\n")
	w := 0
	for _, l := range lines {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(lines) + 1
	w = w + 1

	return &PromptWidget{name: name, x: x, y: y, w: w, h: h, body: body}
}

func (wt *PromptWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	wt.x = float32(maxX) * wt.x
	wt.y = float32(maxY) * wt.y
	v, err := g.SetView(wt.name, int(wt.x), int(wt.y), int(wt.x)+wt.w, int(wt.y)+wt.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprint(v, wt.body)
	}
	return nil
}

type SideWidget struct {
	name  string
	x, y  float32
	w, h  float32
	lines []string
}

func NewSideWidget(name string, x, y, w, h float32, lines []string) *SideWidget {
	return &SideWidget{name: name, x: x, y: y, w: w, h: h, lines: lines}
}

func (wt *SideWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	x := int(wt.x * float32(maxX))
	y := int(wt.y * float32(maxY))
	w := int(wt.w * float32(maxX))
	h := int(wt.h * float32(maxY))

	v, err := g.SetView(wt.name, x, y, w-1, h-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, line := range wt.lines {
			fmt.Fprintln(v, line)
		}
	}

	return nil
}

type CmdWidget struct {
	name string
	x, y float32
	w    float32
}

func NewCmdWidget(name string, x, y, w float32) *CmdWidget {
	return &CmdWidget{name: name, x: x, y: y, w: w}
}

func (wt *CmdWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	x := int(wt.x * float32(maxX))
	y := int(wt.y * float32(maxY))
	w := int(wt.w * float32(maxX))
	h := y + 2

	v, err := g.SetView(wt.name, x, y, w-1, h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editable = true
		v.Wrap = false
	}

	return nil
}

type MainWidget struct {
	name string
	x, y float32
	w, h float32
}

func NewMainWidget(name string, x, y, w, h float32) *MainWidget {
	return &MainWidget{name: name, x: x, y: y, w: w, h: h}
}

func (wt *MainWidget) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	x := int(wt.x * float32(maxX))
	y := int(wt.y * float32(maxY))
	w := int(wt.w * float32(maxX))
	h := int(wt.h * float32(maxY))

	v, err := g.SetView(wt.name, x, y, w-1, h-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
	}

	return nil
}
