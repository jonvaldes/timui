# timui
Immediate Mode UI lib for terminals

Still rough and unfinished, but allows for some quick UI with nice output

Example:

    ╭───────────────────────────────────────────Commands───────────────────────────────────────────╮
    │ ╭───Commands───╮ ╭─Platforms─╮ ╭────────Levels─────────╮                                     │
    │ │ [ ]generate  │ │ [✖]win64  │ │ ❪●❫ Autotest_Post     │                                     │
    │ │ [ ]open      │ │ [ ]ps4    │ │ ❪ ❫ Autotest_Lighting │       ╭──────╮                      │
    │ │ [✖]build     │ │ [✖]xb1    │ │ ❪ ❫ Autotest_Shadows  │       │ Run! │                      │
    │ ⎬──────────────⎨ ╰───────────╯ │ ❪ ❫ Bilbao            │       ╰──────╯                      │
    │ │ [ ]cook      │               │ ❪ ❫ Östermalm         │                                     │
    │ │ [ ]run       │               │ ❪ ❫ Other:            │                                     │
    │ │ [ ]test      │               │ [....................]│                                     │
    │ ╰──────────────╯               ╰───────────────────────╯                                     │
    ╰──────────────────────────────────────────────────────────────────────────────────────────────╯
    ╭────────────────────────────────────────────Output────────────────────────────────────────────╮
    │                                                                                              │
    │ Building...                                                                                  │
    │   ...                                                                                        │
    │   ...                                                                                        │
    │                                                                                              │
    │                                                                                              │
    │                                                                                              │
    │                                                                                              │
    │                                                                                              │
    │                                                                                              │
    ╰──────────────────────────────────────────────────────────────────────────────────────────────╯

Simple example with code:

    ╭─Commands─╮  ╭──────Dirs───────╮   ╭──────╮
    │ [ ]tree  │  │ ❪●❫ /           │   │ Run! │
    │ [ ]ls    │  │ ❪ ❫ ~           │   ╰──────╯
    ╰──────────╯  │ ❪ ❫ ~/Downloads │
                  │ ❪ ❫ Other:      │
                  │ [․․․․․․․․․․․․․․]│
                  ╰─────────────────╯

Drawing code (for the full example see examples/example.go):

```go
func redraw(state *timui.State, data *Data) {
	termbox.Clear(state.Colors.Default, state.Colors.Default)

	timui.Box(state, 2, 1, "Commands",
		&timui.CheckBox{&data.tree, "tree"},
		&timui.CheckBox{&data.ls, "ls"},
	)

	timui.Box(state, 16, 1, "Dirs",
		&timui.RadioBox{0, &data.selectedDir, "/"},
		&timui.RadioBox{1, &data.selectedDir, "~"},
		&timui.RadioBox{2, &data.selectedDir, "~/Downloads"},
		&timui.RadioBox{3, &data.selectedDir, "Other:"},
		&timui.TextEdit{&data.otherDir},
	)

	timui.Box(state, 38, 1, "",
		&timui.Button{"Run!", func() {
			termbox.Close()
			os.Exit(0)
		}},
	)

	state.Flush()
	termbox.Flush()
}
```

