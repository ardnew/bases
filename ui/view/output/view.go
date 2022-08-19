package output

import (
	"github.com/rivo/tview"
)

type View struct {
	*tview.Grid
	table []*table
}

func New() *View {
	// Allocate all primitives before configuring
	return (&View{
		Grid:  tview.NewGrid(),
		table: []*table{},
	}).init()
}

func (v *View) init() *View {
	v.clear()
	v.table = append(v.table, newTable(32), newTable(16), newTable(8))
	pad := 2 // +1 for header, +1 for margin
	row := []int{5}
	height := 5
	for _, t := range v.table {
		h := t.bits.maxRow() + pad
		row = append(row, h)
		height += h
	}
	v.SetColumns(-1).SetRows(row...).AddItem(
		tview.NewTextView().SetTextAlign(tview.AlignCenter).SetText("summary"),
		0, 0, 1, 1, 1, 1, false)
	for i, t := range v.table {
		v.AddItem(t, i+1, 0, 1, 1, height, 0, false)
		height -= t.bits.maxRow() + pad
		remain := height
		for _, u := range v.table[i+1:] {
			v.AddItem(u, i+1, 0, 1, 1, remain, 0, false)
			remain -= u.bits.maxRow() + pad
		}
	}
	return v
}

func (v *View) clear() *View {
	v.table = []*table{}
	v.Clear()
	return v
}
