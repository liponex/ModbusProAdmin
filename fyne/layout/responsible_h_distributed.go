/*
 * Copyright (C) 2023 liponex
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the terms of the GNU General Public License as published by
 * the  Free Software Foundation, either version 3 of the License, or any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.licenses/>.
 */

package layout

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

type responsibleHDistributed struct {
	maxColumns uint
}

type layoutRow struct {
	height    float32
	width     float32
	maxHeight float32
	maxWidth  float32
	items     []fyne.CanvasObject
}

type layoutRowSlice []layoutRow

func (rows *layoutRowSlice) addRow() {
	*rows = append(
		*rows,
		layoutRow{
			height:    0,
			width:     0,
			maxHeight: 0,
			maxWidth:  0,
			items:     []fyne.CanvasObject{},
		},
	)
}

func (resHD *responsibleHDistributed) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		if w < childSize.Width {
			w = childSize.Width
		}
		h += childSize.Height * 2
	}
	return fyne.NewSize(w, h)
}

func (resHD *responsibleHDistributed) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 || objects[0] == nil {
		return
	}

	canvasWidth := containerSize.Width - theme.Padding()

	rows := layoutRowSlice{}
	rows.addRow()

	rowsIter := 0

	for _, o := range objects {
		if uint(len(rows[rowsIter].items)) == resHD.maxColumns {
			rowsIter++
			rows.addRow()
		}

		size := o.MinSize()
		o.Resize(size)

		if rows[rowsIter].maxWidth < size.Width {
			rows[rowsIter].maxWidth = size.Width
		}

		if rows[rowsIter].maxHeight < size.Height {
			rows[rowsIter].maxHeight = size.Height
		}

		rows[rowsIter].items = append(rows[rowsIter].items, o)

		if rows[rowsIter].maxWidth*float32(len(rows[rowsIter].items)+1) > canvasWidth {
			rowsIter++
			rows.addRow()
		}
	}

	prevPos := fyne.NewPos(0, 0)
	for _, row := range rows {
		if len(row.items) < 1 {
			return
		}
		splitSize := containerSize.Width / float32(len(row.items)-1)
		row.items[0].Move(prevPos)
		for j := 1; j < len(row.items)-1; j++ {
			row.items[j].Move(
				prevPos.AddXY(
					float32(j)*splitSize-(row.items[j].Size().Width/2),
					0,
				),
			)
		}
		if len(row.items) > 1 {
			row.items[len(row.items)-1].Move(
				prevPos.AddXY(
					canvasWidth-row.items[len(row.items)-1].Size().Width,
					0,
				),
			)
		}
		prevPos = prevPos.AddXY(0, prevPos.Y+row.maxHeight)
	}
}

func NewResponsibleHDistributedLayout(maxColumns uint, objects ...fyne.CanvasObject) *fyne.Container {
	r := &responsibleHDistributed{
		maxColumns: maxColumns,
	}
	return container.New(r, objects...)
}
