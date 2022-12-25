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
	"fmt"
	"syscall/js"

	"github.com/hajimehoshi/ebiten/v2/internal/ui"
)

var (
	document        = js.Global().Get("document")
	textareaElement js.Value
)

func init() {
	if !document.Truthy() {
		return
	}

	textareaElement = document.Call("createElement", "textarea")
	textareaElement.Set("autocapitalize", "off")
	textareaElement.Set("spellcheck", false)
	textareaElement.Set("translate", "no")
	textareaElement.Set("wrap", "off")

	style := textareaElement.Get("style")
	style.Set("position", "absolute")
	style.Set("left", "0")
	style.Set("top", "0")
	style.Set("opacity", "0")
	style.Set("resize", "none")
	style.Set("cursor", "normal")
	style.Set("pointerEvents", "none")
	style.Set("overflow", "hidden")
	style.Set("tabindex", "-1")
	style.Set("width", "1px")
	style.Set("height", "1px")

	textareaElement.Call("addEventListener", "focusout", js.FuncOf(func(this js.Value, args []js.Value) any {
		e := args[0]
		t := e.Get("target")
		// A mouse click might unfocus the text area unexpectdly.
		// This hack keeps the textarea focused if necessary.
		if document.Get("body").Call("contains", t).Bool() {
			t.Call("focus")
		}
		return nil
	}))
	textareaElement.Call("addEventListener", "keydown", js.FuncOf(func(this js.Value, args []js.Value) any {
		e := args[0]
		if e.Get("code").String() == "Tab" {
			e.Call("preventDefault")
			// TODO: Insert a tab char?
		}
		return nil
	}))
}

func startTextInput() {
	if !textareaElement.Truthy() {
		return
	}

	body := document.Get("body")
	if body.Call("contains", textareaElement).Bool() {
		return
	}

	body.Call("appendChild", textareaElement)
	textareaElement.Call("focus")
}

func endTextInput() {
	if !textareaElement.Truthy() {
		return
	}

	body := document.Get("body")
	if !body.Call("contains", textareaElement).Bool() {
		return
	}

	textareaElement.Set("value", "")
	body.Call("removeChild", textareaElement)
}

func setTextInputPosition(x, y int) {
	if !textareaElement.Truthy() {
		return
	}

	body := document.Get("body")
	if !body.Call("contains", textareaElement).Bool() {
		return
	}

	xf, yf := ui.LogicalPositionToClientPosition(float64(x), float64(y))

	style := textareaElement.Get("style")
	style.Set("left", fmt.Sprintf("%0.2fpx", xf))
	style.Set("top", fmt.Sprintf("%0.2fpx", yf))
}

func textInputStatus() (string, int, int) {
	if !textareaElement.Truthy() {
		return "", 0, 0
	}

	body := document.Get("body")
	if !body.Call("contains", textareaElement).Bool() {
		return "", 0, 0
	}

	text := textareaElement.Get("value").String()
	start := textareaElement.Get("selectionStart").Int()
	end := textareaElement.Get("selectionEnd").Int()
	return text, convertUTF16CountToByteCount(text, start), convertUTF16CountToByteCount(text, end)
}

func setTextInputStatus(text string, start, end int) {
	if !textareaElement.Truthy() {
		return
	}

	body := document.Get("body")
	if !body.Call("contains", textareaElement).Bool() {
		return
	}

	textareaElement.Set("value", text)
	textareaElement.Set("selectionStart", convertByteCountToUTF16Count(text, start))
	textareaElement.Set("selectionEnd", convertByteCountToUTF16Count(text, end))
}
