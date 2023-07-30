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

package widget

import (
	"fyne.io/fyne/v2/widget"
)

type SelectWithPrevLink[T interface{}] struct {
	widget.Select
	curSelected T
	hasPrev     bool
	//OnChanged   func(*T, string) `json:"-"`
}

func NewSelectWithPrevLink[T interface{}](options []string, changed func(*bool, *T, string)) *SelectWithPrevLink[T] {
	sWithPrevLink := &SelectWithPrevLink[T]{}
	sWithPrevLink.ExtendBaseWidget(sWithPrevLink)
	sWithPrevLink.Options = options
	sWithPrevLink.hasPrev = false
	sWithPrevLink.curSelected = *new(T)
	sWithPrevLink.Select.OnChanged = func(selected string) {
		changed(&sWithPrevLink.hasPrev, &sWithPrevLink.curSelected, selected)
	}
	return sWithPrevLink
}

func (s *SelectWithPrevLink[T]) GetSelected() *T {
	return &s.curSelected
}
