// Copyright 2022 The Ebitengine Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package textinput

import (
	"unicode/utf16"
)

// TODO: Check thread safety
var activeCount int

type TextInput struct {
	text           string
	selectionStart int
	selectionEnd   int
	active         bool
}

func (t *TextInput) IsActive() bool {
	return t.active
}

func (t *TextInput) Update(active bool) {
	if active {
		t.activate()
		return
	}
	t.deactivate()
}

func (t *TextInput) activate() {
	if t.active {
		t.text, t.selectionStart, t.selectionEnd = textInputStatus()
		return
	}

	t.active = true
	activeCount++
	startTextInput()
	setTextInputStatus(t.text, t.selectionStart, t.selectionEnd)
}

func (t *TextInput) deactivate() {
	if !t.active {
		return
	}

	activeCount--
	t.active = false
	if activeCount == 0 {
		endTextInput()
	}
}

func (t *TextInput) SetPosition(x, y int) {
	setTextInputPosition(x, y)
}

func (t *TextInput) Text() string {
	return t.text
}

func (t *TextInput) SetText(text string, selectionStart, selectionEnd int) {
	t.text = text
	t.selectionStart = selectionStart
	t.selectionEnd = selectionEnd
	setTextInputStatus(t.text, t.selectionStart, t.selectionEnd)
}

func (t *TextInput) SelectionStart() int {
	return t.selectionStart
}

func (t *TextInput) SelectionEnd() int {
	return t.selectionEnd
}

func convertUTF16CountToByteCount(text string, c int) int {
	return len(string(utf16.Decode(utf16.Encode([]rune(text))[:c])))
}

func convertByteCountToUTF16Count(text string, c int) int {
	return len(utf16.Encode([]rune(text[:c])))
}
