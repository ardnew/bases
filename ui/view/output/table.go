package output

import (
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"

	"github.com/ardnew/bases/num"
)

type column int

const (
	off column = iota
	hex
	oct
	bin
	dec
	colCount
)

var (
	columnID = [colCount]string{
		"",
		num.BaseID[num.Hex],
		num.BaseID[num.Oct],
		num.BaseID[num.Bin],
		num.BaseID[num.Dec],
	}
	columnBase = [colCount]num.Base{
		num.Auto, num.Hex, num.Oct, num.Bin, num.Dec,
	}
)

// RowBits defines the supported bit widths for rows of a table.
type RowBits int

const Byte, Word, Long RowBits = 8, 16, 32

func (b RowBits) maxOffset() int { return num.MaxBits - int(b) }
func (b RowBits) maxRow() int    { return num.MaxBits / int(b) }

type table struct {
	*tview.Table
	bits RowBits
}

func newTable(bits RowBits) *table {
	return (&table{
		Table: tview.NewTable(),
		bits:  bits,
	}).init()
}

func (t *table) init() *table {
	for i, s := range columnID {
		t.SetCell(0, i, tview.NewTableCell(s).
			SetTextColor(tcell.ColorGreen).
			SetAlign(tview.AlignCenter))
	}
	for o := 0; o < t.bits.maxRow(); o++ {
		t.SetCell(o+1, 0, tview.NewTableCell("+"+strconv.Itoa(o*int(t.bits))).
			SetTextColor(tcell.ColorYellow).
			SetAlign(tview.AlignRight))
		for b := off + 1; b < colCount; b++ {
			t.SetCell(o+1, int(b), tview.NewTableCell("-").
				SetTextColor(tcell.ColorWhite).
				SetAlign(tview.AlignCenter))
		}
	}
	t.Select(0, 0).
		SetFixed(1, 1).
		SetSelectable(false, false)
	return t
}
