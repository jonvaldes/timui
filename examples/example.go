package main

import (
	"os"
	"github.com/nsf/termbox-go"
	"github.com/jonvaldes/timgui"
)

type Data struct {
	tree bool
	ls bool

	selectedDir int
	otherDir string
}

func redraw(state *timgui.State, data *Data) {
	termbox.Clear(state.Colors.Default, state.Colors.Default)

	timgui.Box(state, 2, 1, "Commands",
		&timgui.CheckBox{&data.tree, "tree"},
		&timgui.CheckBox{&data.ls, "ls"},
	)

	timgui.Box(state, 16, 1, "Dirs",
		&timgui.RadioBox{0, &data.selectedDir, "/"},
		&timgui.RadioBox{1, &data.selectedDir, "~"},
		&timgui.RadioBox{2, &data.selectedDir, "~/Downloads"},
		&timgui.RadioBox{5, &data.selectedDir, "Other:"},
		&timgui.TextEdit{&data.otherDir},
	)

	timgui.Box(state, 38, 1, "",
		&timgui.Button{"Run!", func() {
			termbox.Close()
			os.Exit(0)
		}},
	)

	state.Flush()
	termbox.Flush()
}

func main() {
	termbox.Init()

	data := Data{}

	state := timgui.NewState()
	state.Colors.Selected = termbox.ColorCyan
	state.Colors.Cursor = termbox.ColorCyan | termbox.AttrBold
	redraw(&state, &data)

mainloop:
	for {
		ev := termbox.PollEvent()
		state.HandleEvent(ev)
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			break mainloop
		}
	repeat:
		redraw(&state, &data)
		if state.NeedsRedraw {
			state.NeedsRedraw = false
			goto repeat
		}
	}

	termbox.Close()
}
