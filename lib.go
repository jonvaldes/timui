package timgui

import (
	"github.com/nsf/termbox-go"
)

type Coords struct {
	x, y int
}

type State struct {
	Cursor Coords
	Colors struct {
		Default  termbox.Attribute
		Selected termbox.Attribute
		Cursor   termbox.Attribute
	}
	KeyState    map[termbox.Key]bool
	InputChar   rune
	NeedsRedraw bool
	boxesCnt    int
}

func NewState() State {
	result := State{}
	result.KeyState = make(map[termbox.Key]bool)
	return result
}

func (state *State) HandleEvent(ev termbox.Event) {
	if ev.Type == termbox.EventKey {
		state.KeyState[ev.Key] = true
		state.InputChar = ev.Ch
	}
}

func (state *State) Flush() {
	if state.Cursor.x >= state.boxesCnt {
		state.Cursor.x = state.boxesCnt - 1
		state.NeedsRedraw = true
	}

	state.boxesCnt = 0
	state.KeyState = make(map[termbox.Key]bool)
}

func WriteText(x, y int, text string, fg, bg termbox.Attribute) {
	dx := 0
	for _, r := range text {
		termbox.SetCell(x+dx, y, r, fg, bg)
		dx++
	}
}

type Elem interface {
	maxWidth() int
	canBeSelected() bool
	draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool)
}

type TextEdit struct {
	Text *string
}

func (t *TextEdit) maxWidth() int {
	return len(*t.Text) + 2
}

func (t *TextEdit) canBeSelected() bool {
	return true
}

func (t *TextEdit) draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {
	fgColor := state.Colors.Default
	if isElemSelected {
		fgColor = state.Colors.Cursor
	}

	if isElemSelected && state.InputChar != 0 {
		*t.Text += string(state.InputChar)
		state.InputChar = 0
	}

	if isElemSelected && state.KeyState[termbox.KeySpace] {
		*t.Text += " "
	}

	if isElemSelected &&
		len(*t.Text) > 0 &&
		(state.KeyState[termbox.KeyBackspace] ||
		state.KeyState[termbox.KeyBackspace2] ) {
		*t.Text = (*t.Text)[:len(*t.Text) - 1]
	}

	termbox.SetCell(pos.x, pos.y, '[', fgColor, state.Colors.Default)
	termbox.SetCell(pos.x+maxWidth, pos.y, ']', fgColor, state.Colors.Default)

	WriteText(pos.x+1, pos.y, *t.Text, fgColor, state.Colors.Default)
	for x := 1 + len(*t.Text); x < maxWidth; x++ {
		// ․‥…▁
		termbox.SetCell(pos.x+x, pos.y, '․', state.Colors.Default, state.Colors.Default)
	}
}

type RadioBox struct {
	ID    int
	Value *int
	Text  string
}

func (r *RadioBox) maxWidth() int {
	return len(r.Text) + 4
}

func (r *RadioBox) canBeSelected() bool {
	return true
}

func (r *RadioBox) draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {
	if isElemSelected && state.KeyState[termbox.KeySpace] {
		*r.Value = r.ID
		state.NeedsRedraw = true
	}

	fgColor := state.Colors.Default
	if isElemSelected {
		fgColor = state.Colors.Cursor
	}
	bgColor := state.Colors.Default

	//brackets := "( )"
	brackets := "❪ ❫ "

	WriteText(pos.x, pos.y, brackets+r.Text, fgColor, bgColor)

	// checkMark := '◉'
	checkMark := '●'
	// ○
	if *r.Value == r.ID {
		termbox.SetCell(pos.x+1, pos.y, checkMark, fgColor, bgColor)
	}
}

type Separator struct {
	Text string
}

func (s *Separator) maxWidth() int {
	return len(s.Text)
}

func (s *Separator) canBeSelected() bool {
	return false
}

func (s *Separator) draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {
	fgColor := state.Colors.Default
	if isBoxSelected {
		fgColor = state.Colors.Selected
	}
	for dx := -1; dx < maxWidth+1; dx++ {
		termbox.SetCell(pos.x+dx, pos.y, '─', fgColor, state.Colors.Default)
	}
	termbox.SetCell(pos.x-2, pos.y, '⎬', fgColor, state.Colors.Default)
	termbox.SetCell(pos.x+maxWidth+1, pos.y, '⎨', fgColor, state.Colors.Default)

	textPos := (maxWidth - len(s.Text)) / 2
	WriteText(pos.x+textPos, pos.y, s.Text, fgColor, state.Colors.Default)
}

type Button struct {
	Text     string
	Callback func()
}

func (b *Button) maxWidth() int {
	return len(b.Text)
}

func (b *Button) canBeSelected() bool {
	return true
}

func (b *Button) draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {

	fgColor := state.Colors.Default
	bgColor := state.Colors.Default

	if isElemSelected {
		if state.KeyState[termbox.KeySpace] || state.KeyState[termbox.KeyEnter] {
			b.Callback()
		}
		bgColor = state.Colors.Selected
	}
	textPos := (maxWidth - len(b.Text)) / 2
	WriteText(pos.x+textPos, pos.y, b.Text, fgColor, bgColor)
}

type CheckBox struct {
	Value *bool
	Text  string
}

func (c *CheckBox) maxWidth() int {
	return len(c.Text) + 4 // "[X] <text>"
}

func (c *CheckBox) canBeSelected() bool {
	return true
}

func (c *CheckBox) draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {
	if isElemSelected && state.KeyState[termbox.KeySpace] {
		*c.Value = !*c.Value
	}

	fgColor := state.Colors.Default
	if isElemSelected {
		fgColor = state.Colors.Cursor
	}
	bgColor := state.Colors.Default

	brackets := "[ ]"

	//checkMark := '✓'
	// checkMark := '✔'
	// checkMark := '✖'
	// checkMark := '✘'
	if *c.Value {
		brackets = "[✖]"
	}
	WriteText(pos.x, pos.y, brackets+c.Text, fgColor, bgColor)
}

func Box(state *State, x, y int, title string, elems ...Elem) {

	bgColor := state.Colors.Default
	fgColor := state.Colors.Default

	// Calculate box size
	maxWidth := len(title)
	selectableElems := 0
	for _, e := range elems {
		width := e.maxWidth()
		if width > maxWidth {
			maxWidth = width
		}
		if e.canBeSelected() {
			selectableElems++
		}
	}

	if state.Cursor.x == state.boxesCnt {
		if state.KeyState[termbox.KeyArrowUp] && state.Cursor.y > 0 {
			state.Cursor.y--
		}

		if state.KeyState[termbox.KeyArrowDown] {
			state.Cursor.y++
		}

		if state.Cursor.y >= selectableElems {
			state.Cursor.y = selectableElems - 1
		}

		if state.KeyState[termbox.KeyArrowLeft] && state.Cursor.x > 0 {
			state.Cursor.x--
			state.KeyState[termbox.KeyArrowLeft] = false
			state.NeedsRedraw = true
		}

		if state.KeyState[termbox.KeyArrowRight] {
			state.Cursor.x++
			state.KeyState[termbox.KeyArrowRight] = false
			state.NeedsRedraw = true
		}
	}

	selected := state.Cursor.x == state.boxesCnt
	if selected {
		fgColor = state.Colors.Selected
	}

	// Draw box
	boxWidth := maxWidth + 4
	boxHeight := len(elems) + 2

	DrawBox(x, y, boxWidth, boxHeight, title, fgColor, bgColor)

	// Draw elements
	elem := -1
	for i, e := range elems {
		if e.canBeSelected() {
			elem++
		}
		elemSelected := (selected && state.Cursor.y == elem)
		e.draw(state, Coords{x: x + 2, y: y + 1 + i}, maxWidth, selected, elemSelected)
	}

	state.boxesCnt++
}

func DrawBox(x, y, w, h int, title string, fgColor, bgColor termbox.Attribute) {
	termbox.SetCell(x, y, '╭', fgColor, bgColor)
	termbox.SetCell(x+w-1, y, '╮', fgColor, bgColor)
	termbox.SetCell(x, y+h-1, '╰', fgColor, bgColor)
	termbox.SetCell(x+w-1, y+h-1, '╯', fgColor, bgColor)
	for dx := 1; dx < w-1; dx++ {
		termbox.SetCell(x+dx, y, '─', fgColor, bgColor)
		termbox.SetCell(x+dx, y+h-1, '─', fgColor, bgColor)
	}

	for dy := 0; dy < h-2; dy++ {
		termbox.SetCell(x, y+dy+1, '│', fgColor, bgColor)
		termbox.SetCell(x+w-1, y+dy+1, '│', fgColor, bgColor)
	}

	titlePos := (w - len(title)) / 2
	WriteText(x+titlePos, y, title, fgColor, bgColor)
}
