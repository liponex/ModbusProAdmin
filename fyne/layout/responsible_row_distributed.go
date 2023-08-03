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
)

type responsibleRowDistributed struct {
	maxColumns int
}

type layoutRow struct {
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
			width:     0,
			maxHeight: 0,
			maxWidth:  0,
			items:     []fyne.CanvasObject{},
		},
	)
}

func splitRows(objects []fyne.CanvasObject, containerWidth float32, maxColumns int) *layoutRowSlice {
	rows := new(layoutRowSlice)
	rows.addRow()

	rowsIter := 0

	for _, o := range objects {
		if len((*rows)[rowsIter].items) == maxColumns {
			rowsIter++
			rows.addRow()
		}

		size := o.MinSize()
		o.Resize(size)

		if (*rows)[rowsIter].width+size.Width > containerWidth {
			rowsIter++
			rows.addRow()
		}

		(*rows)[rowsIter].width += size.Width

		if (*rows)[rowsIter].maxWidth < size.Width {
			(*rows)[rowsIter].maxWidth = size.Width
		}

		if (*rows)[rowsIter].maxHeight < size.Height {
			(*rows)[rowsIter].maxHeight = size.Height
		}

		(*rows)[rowsIter].items = append((*rows)[rowsIter].items, o)
	}

	return rows
}

func (resRD *responsibleRowDistributed) MinSize(objects []fyne.CanvasObject) fyne.Size {
	w, h := float32(0), float32(0)
	for _, o := range objects {
		childSize := o.MinSize()

		if w < childSize.Width {
			w = childSize.Width
		}
	}
	rows := splitRows(objects, w, resRD.maxColumns)
	for _, row := range *rows {
		h += row.maxHeight
	}
	return fyne.NewSize(w, h)
}

func (resRD *responsibleRowDistributed) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 || objects[0] == nil {
		return
	}

	canvasWidth := containerSize.Width

	rows := splitRows(objects, canvasWidth, resRD.maxColumns)

	prevPos := fyne.NewPos(0, 0)
	for _, row := range *rows {
		if len(row.items) == 0 {
			continue
		}

		row.items[0].Move(prevPos)
		if len(row.items) == 1 {
			prevPos = prevPos.AddXY(0, row.maxHeight)
			continue
		}
		prevItemsBounds := prevPos.AddXY(row.items[0].Size().Width, 0)

		splitSize := containerSize.Width / float32(len(row.items)-1)
		for j, item := range row.items[1 : len(row.items)-1] {
			if prevItemsBounds.X > float32(j+1)*splitSize-(item.Size().Width/2) {
				item.Move(prevPos.AddXY(
					prevItemsBounds.X,
					0,
				))
				prevItemsBounds.X += item.Size().Width
				continue
			}
			item.Move(prevPos.AddXY(
				float32(j+1)*splitSize-(item.Size().Width/2),
				0,
			))
		}

		row.items[len(row.items)-1].Move(prevPos.AddXY(
			canvasWidth-row.items[len(row.items)-1].Size().Width,
			0,
		))

		prevPos.Y += row.maxHeight
	}
}

func NewResponsibleHDistributedLayout(maxColumns int, objects ...fyne.CanvasObject) *fyne.Container {
	r := &responsibleRowDistributed{
		maxColumns: maxColumns,
	}
	return container.New(r, objects...)
}
