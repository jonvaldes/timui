package timui

import (
	"github.com/nsf/termbox-go"
)

// Coords defines screen or UI coordinates
type Coords struct {
	X, Y int
}

// State carries all global information for the UI
type State struct {
	Cursor      Coords
	MouseCursor Coords
	MouseClick  bool
	Colors      struct {
		Default  termbox.Attribute
		Selected termbox.Attribute
		Cursor   termbox.Attribute
	}
	KeyState          map[termbox.Key]bool
	InputChar         rune
	NeedsRedraw       bool
	boxesCnt          int
	CurrentElemCursor Coords
}

// NewState returns a new State
func NewState() State {
	result := State{}
	result.KeyState = make(map[termbox.Key]bool)
	return result
}

// HandleEvent registers an event in the UI system
func (state *State) HandleEvent(ev termbox.Event) {
	switch ev.Type {
	case termbox.EventKey:
		state.KeyState[ev.Key] = true
		state.InputChar = ev.Ch
	case termbox.EventMouse:
		if ev.Key == termbox.MouseLeft {
			state.MouseCursor.X = ev.MouseX
			state.MouseCursor.Y = ev.MouseY
			state.MouseClick = true
		}
	}
}

// Flush resets the UI for the next frame
func (state *State) Flush() {
	if state.Cursor.X >= state.boxesCnt {
		state.Cursor.X = state.boxesCnt - 1
		state.NeedsRedraw = true
	}

	state.boxesCnt = 0
	// Reset input
	state.KeyState = make(map[termbox.Key]bool)
	state.MouseClick = false
}

// WriteText writes a line of text at the given position, with the specified attribs
func WriteText(x, y int, text string, fg, bg termbox.Attribute) {
	dx := 0
	for _, r := range text {
		termbox.SetCell(x+dx, y, r, fg, bg)
		dx++
	}
}

// Elem is a UI element
type Elem interface {
	// MaxWidth returns the maximum with the element expects to need to draw correctly
	MaxWidth() int

	// CanBeSelected returns whether the element can be manually selected by the user
	CanBeSelected() bool

	// Draw executes all drawing commands and input handling
	Draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool)
}

// TextEdit is an element that allows the user to edit a one-line string
type TextEdit struct {
	Text *string
}

func (t *TextEdit) MaxWidth() int {
	return len(*t.Text) + 2
}

func (t *TextEdit) CanBeSelected() bool {
	return true
}

func (t *TextEdit) Draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {
	fgColor := state.Colors.Default

	if isElemSelected {
		fgColor = state.Colors.Cursor
	}

	if isElemSelected && state.InputChar != 0 {
		*t.Text += string(state.InputChar)
		state.InputChar = 0
		state.NeedsRedraw = true
	}

	if isElemSelected && (state.KeyState[termbox.KeySpace]){
		*t.Text += " "
		state.NeedsRedraw = true
	}

	if isElemSelected &&
		len(*t.Text) > 0 &&
		(state.KeyState[termbox.KeyBackspace] ||
			state.KeyState[termbox.KeyBackspace2]) {
		*t.Text = (*t.Text)[:len(*t.Text)-1]
		state.NeedsRedraw = true
	}

	termbox.SetCell(pos.X, pos.Y, '[', fgColor, state.Colors.Default)
	termbox.SetCell(pos.X+maxWidth-1, pos.Y, ']', fgColor, state.Colors.Default)

	WriteText(pos.X+1, pos.Y, *t.Text, fgColor, state.Colors.Default)
	for x := 1 + len(*t.Text); x < maxWidth-1; x++ {
		// ․‥…▁
		termbox.SetCell(pos.X+x, pos.Y, '․', state.Colors.Default, state.Colors.Default)
	}
}

// RadioBox is an element that allows the user to select one of several alternatives
type RadioBox struct {
	ID    int
	Value *int
	Text  string
}

func (r *RadioBox) MaxWidth() int {
	return len(r.Text) + 4
}

func (r *RadioBox) CanBeSelected() bool {
	return true
}

func (r *RadioBox) Draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {
	if isElemSelected && (state.KeyState[termbox.KeySpace] || state.MouseClick){
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

	WriteText(pos.X, pos.Y, brackets+r.Text, fgColor, bgColor)

	// checkMark := '◉'
	checkMark := '●'
	// ○
	if *r.Value == r.ID {
		termbox.SetCell(pos.X+1, pos.Y, checkMark, fgColor, bgColor)
	}
}

// Separator is a visual element with no interactive functionality.
type Separator struct {
	Text string
}

func (s *Separator) MaxWidth() int {
	return len(s.Text)
}

func (s *Separator) CanBeSelected() bool {
	return false
}

func (s *Separator) Draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {
	fgColor := state.Colors.Default
	if isBoxSelected {
		fgColor = state.Colors.Selected
	}
	for dx := -1; dx < maxWidth+1; dx++ {
		termbox.SetCell(pos.X+dx, pos.Y, '─', fgColor, state.Colors.Default)
	}
	termbox.SetCell(pos.X-2, pos.Y, '⎬', fgColor, state.Colors.Default)
	termbox.SetCell(pos.X+maxWidth+1, pos.Y, '⎨', fgColor, state.Colors.Default)

	textPos := (maxWidth - len(s.Text)) / 2
	WriteText(pos.X+textPos, pos.Y, s.Text, fgColor, state.Colors.Default)
}

// Button is an element that allows for running arbitrary code when pressed
type Button struct {
	Text     string
	Callback func()
}

func (b *Button) MaxWidth() int {
	return len(b.Text)
}

func (b *Button) CanBeSelected() bool {
	return true
}

func (b *Button) Draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {

	fgColor := state.Colors.Default
	bgColor := state.Colors.Default

	if isElemSelected {
		if state.KeyState[termbox.KeySpace] || state.KeyState[termbox.KeyEnter] || state.MouseClick{
			b.Callback()
		}
		bgColor = state.Colors.Selected
	}
	textPos := (maxWidth - len(b.Text)) / 2
	WriteText(pos.X+textPos, pos.Y, b.Text, fgColor, bgColor)
}

// CheckBox is an element that allows the user to toggle a boolean value
type CheckBox struct {
	Value *bool
	Text  string
}

func (c *CheckBox) MaxWidth() int {
	return len(c.Text) + 4 // "[X] <text>"
}

func (c *CheckBox) CanBeSelected() bool {
	return true
}

func (c *CheckBox) Draw(state *State, pos Coords, maxWidth int, isBoxSelected, isElemSelected bool) {
	if isElemSelected && (state.KeyState[termbox.KeySpace] || state.MouseClick){
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
	WriteText(pos.X, pos.Y, brackets+c.Text, fgColor, bgColor)
}

// DrawBox draws a rectangular box with the given size, title and attributes
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

func checkMouseClick(state *State, pos Coords, width int) bool {
	return state.MouseCursor.Y == pos.Y &&
		state.MouseCursor.X >= pos.X &&
		state.MouseCursor.X <= pos.X+width
}

// Box groups elements visually and functionally.
func Box(state *State, x, y int, title string, elems ...Elem) {

	// Make the mouse click visible only to clicked elements
	pendingMouseClick := state.MouseClick
	state.MouseClick = false


	bgColor := state.Colors.Default
	fgColor := state.Colors.Default

	// Calculate box size
	maxWidth := len(title)
	selectableElems := 0
	for _, e := range elems {
		width := e.MaxWidth()
		if width > maxWidth {
			maxWidth = width
		}
		if e.CanBeSelected() {
			selectableElems++
		}
	}

	if state.Cursor.X == state.boxesCnt {
		if state.KeyState[termbox.KeyArrowUp] && state.Cursor.Y > 0 {
			state.Cursor.Y--
		}

		if state.KeyState[termbox.KeyArrowDown] {
			state.Cursor.Y++
		}

		if state.Cursor.Y >= selectableElems {
			state.Cursor.Y = selectableElems - 1
		}

		if state.KeyState[termbox.KeyArrowLeft] && state.Cursor.X > 0 {
			state.Cursor.X--
			state.KeyState[termbox.KeyArrowLeft] = false
			state.NeedsRedraw = true
		}

		if state.KeyState[termbox.KeyArrowRight] {
			state.Cursor.X++
			state.KeyState[termbox.KeyArrowRight] = false
			state.NeedsRedraw = true
		}
	}

	selected := state.Cursor.X == state.boxesCnt
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
		if e.CanBeSelected() {
			elem++
		}
		elemSelected := (selected && state.Cursor.Y == elem)
		elemScreenPos := Coords{x + 2, y + 1 + i}
		if pendingMouseClick {
			if checkMouseClick(state, elemScreenPos, maxWidth) {
				elemSelected = true
				state.Cursor = Coords{state.boxesCnt, elem}
				state.NeedsRedraw = true
				state.MouseClick = true
				pendingMouseClick = false
			}
		}
		e.Draw(state, elemScreenPos, maxWidth, selected, elemSelected)
		state.MouseClick = false
	}

	state.boxesCnt++
	state.MouseClick = pendingMouseClick
}
