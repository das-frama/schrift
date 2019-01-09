package schrift

import (
	"fmt"
	"time"

	"github.com/jroimartin/gocui"
)

type Command int8
type State int8

// Commands.
const (
	IdleCommand Command = iota
	ContinueCommand
	QuitCommand
)

// States.
const (
	IdleState State = iota
	RunningState
	UserState
)

// Message represents a single message that appears in terminal.
// It could have speed and delay.
type Message struct {
	Text     string
	Speed    time.Duration
	Delay    time.Duration
	State    State
	NewLines int

	Questions []*Question
}

type Question struct {
	UserAnswers []string
	Answer      string
	Command     Command
}

func (m Message) Send(g *gocui.Gui, v *gocui.View) {
	// Delay before sending the message.
	time.Sleep(m.Delay)

	if m.Speed == 0 {
		g.Update(func(g *gocui.Gui) error {
			fmt.Fprintf(v, "%s", m.Text)
			return nil
		})
		return
	}

	chnl := make(chan rune, len(m.Text))
	go func() {
		for {
			c, more := <-chnl
			if more {
				g.Update(func(g *gocui.Gui) error {
					fmt.Fprintf(v, "%c", c)
					return nil
				})
			} else {
				fmt.Fprintf(v, " ")
				for i := 0; i < m.NewLines; i++ {
					fmt.Fprintln(v)
				}
				return
			}
		}
	}()

	for _, r := range m.Text {
		chnl <- r
		time.Sleep(m.Speed)
	}

	close(chnl)
}
