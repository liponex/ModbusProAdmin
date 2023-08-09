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

type responsibleRowDistributed struct {
	maxColumns int
	actualSize fyne.Size
}

type layoutRow struct {
	width     float32
	maxHeight float32
	maxWidth  float32
	items     []fyne.CanvasObject
	alignment string
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
			alignment: "right",
		},
	)
}

func splitRows(objects []fyne.CanvasObject, containerWidth float32, maxColumns int) *layoutRowSlice {
	rows := new(layoutRowSlice)
	rows.addRow()

	rowsIter := 0

	for _, o := range objects {
		if o == nil || !o.Visible() {
			continue
		}

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
	minSize := fyne.NewSize(0, 0)
	for _, child := range objects {
		if !child.Visible() {
			continue
		}
		minSize.Width = fyne.Max(minSize.Width, child.MinSize().Width)
	}
	var rows *layoutRowSlice
	if resRD.actualSize.Width < minSize.Width {
		rows = splitRows(objects, minSize.Width, resRD.maxColumns)
	} else {
		rows = splitRows(objects, resRD.actualSize.Width, resRD.maxColumns)
	}
	for _, row := range *rows {
		minSize.Height += row.maxHeight
	}
	minSize.Height += theme.InnerPadding()
	minSize.Width += theme.InnerPadding()
	return minSize
}

func (resRD *responsibleRowDistributed) Layout(objects []fyne.CanvasObject, containerSize fyne.Size) {
	if len(objects) == 0 || objects[0] == nil {
		return
	}

	canvasWidth := containerSize.Width - theme.InnerPadding()

	rows := splitRows(objects, canvasWidth, resRD.maxColumns)

	prevPos := fyne.NewPos(0, 0)
	for _, row := range *rows {
		switch row.alignment {
		case "center":
			row.rowItemsCenter(prevPos, canvasWidth)
			break
		case "left":
			row.rowItemsLeft(prevPos, canvasWidth)
			break
		case "right":
			row.rowItemsRight(prevPos, canvasWidth)
			break
		case "distributed":
			row.rowItemsDistributed(prevPos, canvasWidth)
			break
		default:
			row.rowItemsDistributed(prevPos, canvasWidth)
			break
		}
		prevPos.Y += row.maxHeight + theme.Padding()
	}

	resRD.actualSize = fyne.NewSize(canvasWidth, prevPos.Y)
}

func NewResponsibleRowDistributedLayout(maxColumns int, objects ...fyne.CanvasObject) *fyne.Container {
	r := &responsibleRowDistributed{
		maxColumns: maxColumns,
		actualSize: fyne.NewSize(0, 0),
	}
	return container.New(r, objects...)
}

/*
rowItemsDistributed distributes items in a row with equal spacing between them
and the first and last items are aligned with the left and right edges of the row
*/
func (row *layoutRow) rowItemsDistributed(initialPos fyne.Position, canvasWidth float32) {
	if len(row.items) == 0 {
		return
	}

	row.items[0].Move(initialPos)
	if len(row.items) == 1 {
		return
	}
	prevItemXBound := row.items[0].Size().Width + theme.Padding()

	splitSize := canvasWidth / float32(len(row.items)-1)
	for j, item := range row.items[1 : len(row.items)-1] {
		if prevItemXBound+theme.Padding() > float32(j+1)*splitSize-(item.Size().Width/2) {
			item.Move(initialPos.AddXY(
				prevItemXBound+theme.Padding(),
				0,
			))
			prevItemXBound += item.Size().Width + theme.Padding()
			continue
		}
		item.Move(initialPos.AddXY(
			float32(j+1)*splitSize-(item.Size().Width/2)+theme.Padding(),
			0,
		))
	}

	row.items[len(row.items)-1].Move(initialPos.AddXY(
		canvasWidth-row.items[len(row.items)-1].Size().Width,
		0,
	))
}

/*
rowItemsLeft aligns items in a row with the left edge of the row
*/
func (row *layoutRow) rowItemsLeft(initialPos fyne.Position, canvasWidth float32) {
	if len(row.items) == 0 {
		return
	}

	row.items[0].Move(initialPos)
	if len(row.items) == 1 {
		return
	}
	prevItemXBound := row.items[0].Size().Width + theme.Padding()

	for _, item := range row.items[1:] {
		item.Move(fyne.NewPos(
			prevItemXBound,
			initialPos.Y,
		))
		prevItemXBound += item.Size().Width + theme.Padding()
	}

}

/*
rowItemsRight aligns items in a row with the right edge of the row
*/
func (row *layoutRow) rowItemsRight(initialPos fyne.Position, canvasWidth float32) {
	if len(row.items) == 0 {
		return
	}

	canvasWidth += theme.InnerPadding()
	prevItemXBound := canvasWidth - row.items[len(row.items)-1].Size().Width - theme.Padding()
	row.items[len(row.items)-1].Move(fyne.NewPos(
		prevItemXBound,
		initialPos.Y,
	))
	if len(row.items) == 1 {
		return
	}

	var flippedItems = row.items
	for i, j := 0, len(flippedItems)-1; i < j; i, j = i+1, j-1 {
		flippedItems[i], flippedItems[j] = flippedItems[j], flippedItems[i]
	}
	for _, item := range flippedItems[1:] {
		prevItemXBound -= item.Size().Width + theme.Padding()
		item.Move(fyne.NewPos(
			prevItemXBound,
			initialPos.Y,
		))
	}
}

/*
rowItemsCenter aligns items in a row with the center of the row
*/
func (row *layoutRow) rowItemsCenter(initialPos fyne.Position, canvasWidth float32) {
	if len(row.items) == 0 {
		return
	}

	prevItemXBound := (canvasWidth / 2) - ((row.width + theme.Padding()*float32(len(row.items))) / 2)
	row.items[0].Move(fyne.NewPos(
		prevItemXBound,
		initialPos.Y,
	))
	if len(row.items) == 1 {
		return
	}
	prevItemXBound += row.items[0].Size().Width + theme.Padding()

	for _, item := range row.items[1:] {
		item.Move(fyne.NewPos(
			prevItemXBound,
			initialPos.Y,
		))
		prevItemXBound += item.Size().Width + theme.Padding()
	}
}
