package input

import (
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ardnew/bases/num"
)

type View struct {
	*tview.Flex
	edit  *tview.InputField
	base  *tview.DropDown
	debug *os.File
}

func New() *View {
	// Allocate all primitives before configuring
	return (&View{
		Flex: tview.NewFlex(),
		edit: tview.NewInputField(),
		base: tview.NewDropDown(),
	}).init()
}

const (
	editIndex = iota
	baseIndex
)

func (v *View) init() *View {
	var err error
	v.debug, err = os.OpenFile("debug", os.O_WRONLY, 0)
	if err != nil {
		panic(err)
	}
	v.edit.SetLabel("> ").
		SetChangedFunc(v.change)
	v.base.SetLabel(" ").
		SetOptions(num.BaseID[:], v.setBase).
		SetCurrentOption(int(num.Auto))
	v.AddItem(v.edit, 0, 1, true)
	v.AddItem(v.base, 5, 0, false)
	return v
}

func (v *View) Focus(delegate func(p tview.Primitive)) {
	if v.edit != nil {
		delegate(v.edit)
	} else {
		v.Box.Focus(delegate)
	}
}

func (v *View) HasFocus() bool {
	return (v.edit != nil && v.edit.HasFocus()) ||
		(v.base != nil && v.base.HasFocus()) ||
		v.Box.HasFocus()
}

func (v *View) InputHandler() func(*tcell.EventKey, func(p tview.Primitive)) {
	return v.WrapInputHandler(func(
		event *tcell.EventKey, setFocus func(p tview.Primitive),
	) {
		if v.base != nil && v.base.HasFocus() {
			if h := v.base.InputHandler(); h != nil {
				h(event, setFocus)
				// Return to InputField edit once a base is selected via Enter key.
				if event.Key() == tcell.KeyEnter {
					h(tcell.NewEventKey(tcell.KeyTab, '\t', tcell.ModShift), setFocus)
				}
			}
		} else {
			mod := event.Modifiers()
			switch event.Key() {
			case tcell.KeyBacktab, tcell.KeyTab:
				// Open the base DropDown via Tab key.
				if v.base != nil {
					v.base.SetDoneFunc(func(tcell.Key) { setFocus(v.edit) })
					if h := v.base.InputHandler(); h != nil {
						h(tcell.NewEventKey(tcell.KeyEnter, '\r', tcell.ModNone), setFocus)
					}
				}
			case tcell.KeyRune:
				if r := event.Rune(); mod == tcell.ModCtrl || (mod == tcell.ModNone &&
					('0' <= r && r <= '9' || 'a' <= r && r <= 'f' || 'A' <= r && r <= 'F')) {
					// Pass event on to InputField edit.
					if v.edit != nil && v.edit.HasFocus() {
						if h := v.edit.InputHandler(); h != nil {
							h(event, setFocus)
						}
					}
				}
			case tcell.KeyEnter:
			default:
			}
		}
	})
}

func (v *View) MouseHandler() func(
	tview.MouseAction, *tcell.EventMouse, func(p tview.Primitive),
) (bool, tview.Primitive) {
	return v.WrapMouseHandler(func(
		action tview.MouseAction, event *tcell.EventMouse, setFocus func(p tview.Primitive),
	) (consumed bool, capture tview.Primitive) {
		if !v.InRect(event.Position()) {
			return false, nil
		}
		// Pass mouse events down.
		for _, v := range []tview.Primitive{v.edit, v.base} {
			if v != nil {
				consumed, capture = v.MouseHandler()(action, event, setFocus)
				if consumed {
					return
				}
			}
		}
		return true, nil
	})
}

func (v *View) change(text string) {
}

func (v *View) setBase(option string, index int) {
}
