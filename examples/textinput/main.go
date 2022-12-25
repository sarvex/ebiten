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

package main

import (
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/hajimehoshi/bitmapfont/v2"
	"golang.org/x/image/font"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/exp/textinput"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type TextField struct {
	bounds     image.Rectangle
	multilines bool
	textInput  textinput.TextInput
}

func NewTextField(bounds image.Rectangle, multilines bool) *TextField {
	return &TextField{
		bounds:     bounds,
		multilines: multilines,
	}
}

func (t *TextField) Contains(x, y int) bool {
	r := t.bounds
	return r.Min.X <= x && x < r.Max.X && r.Min.Y <= y && y < r.Max.Y
}

func (t *TextField) Update(isEditing bool) {
	t.textInput.Update(isEditing)

	if !isEditing {
		return
	}

	if !t.multilines {
		orig := t.textInput.Text()
		new := strings.ReplaceAll(orig, "\n", "")
		if new != orig {
			start := t.textInput.SelectionStart()
			end := t.textInput.SelectionEnd()
			start -= strings.Count(orig[:start], "\n")
			end -= strings.Count(orig[:end], "\n")
			t.textInput.SetText(new, start, end)
		}
	}

	x, y := t.bounds.Min.X, t.bounds.Min.Y
	cx, cy := t.cursorPos()
	px, py := textFieldPadding()
	ox, oy := textFieldOffsetToDot()
	x += cx + px + ox
	y += cy + py + oy
	t.textInput.SetPosition(x, y)
}

func (t *TextField) cursorPos() (int, int) {
	var nlCount int
	lastNLPos := -1
	for i, r := range t.textInput.Text()[:t.textInput.SelectionStart()] {
		if r == '\n' {
			nlCount++
			lastNLPos = i
		}
	}
	x := font.MeasureString(bitmapfont.Face, t.textInput.Text()[lastNLPos+1:t.textInput.SelectionStart()]).Ceil()
	y := nlCount * bitmapfont.Face.Metrics().Height.Ceil()
	return x, y
}

func (t *TextField) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, float32(t.bounds.Min.X), float32(t.bounds.Min.Y), float32(t.bounds.Dx()), float32(t.bounds.Dy()), color.White)
	var clr color.Color = color.Black
	if t.textInput.IsActive() {
		clr = color.RGBA{0, 0, 0xff, 0xff}
	}
	vector.StrokeRect(screen, float32(t.bounds.Min.X), float32(t.bounds.Min.Y), float32(t.bounds.Dx()), float32(t.bounds.Dy()), 0.5, clr)

	px, py := textFieldPadding()
	if t.textInput.IsActive() && t.textInput.SelectionStart() >= 0 {
		x, y := t.bounds.Min.X, t.bounds.Min.Y
		cx, cy := t.cursorPos()
		x += px + cx
		y += py + cy
		h := bitmapfont.Face.Metrics().Height.Ceil()
		vector.StrokeLine(screen, float32(x), float32(y), float32(x), float32(y+h), 0.5, color.Black)
	}

	ox, oy := textFieldOffsetToDot()
	tx := t.bounds.Min.X + px + ox
	ty := t.bounds.Min.Y + py + oy
	text.Draw(screen, t.textInput.Text(), bitmapfont.Face, tx, ty, color.Black)
}

const textFieldHeight = 24

func textFieldPadding() (int, int) {
	m := bitmapfont.Face.Metrics()
	return 4, (textFieldHeight - m.Height.Ceil()) / 2
}

func textFieldOffsetToDot() (int, int) {
	m := bitmapfont.Face.Metrics()
	// TODO: Avoid a magic number
	return 4, m.Ascent.Ceil()
}

type Game struct {
	textFields      []*TextField
	activeTextField *TextField
}

func (g *Game) Update() error {
	if g.textFields == nil {
		g.textFields = append(g.textFields, NewTextField(image.Rect(16, 16, screenWidth-16, 16+textFieldHeight), false))
		g.textFields = append(g.textFields, NewTextField(image.Rect(16, 48, screenWidth-16, 48+textFieldHeight), false))
		g.textFields = append(g.textFields, NewTextField(image.Rect(16, 80, screenWidth-16, screenHeight-16), true))
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.activeTextField = nil

		x, y := ebiten.CursorPosition()
		for _, tf := range g.textFields {
			if tf.Contains(x, y) {
				g.activeTextField = tf
				break
			}
		}
	}

	for _, tf := range g.textFields {
		tf.Update(g.activeTextField == tf)
	}

	ebiten.SetCursorShape(ebiten.CursorShapeDefault)
	x, y := ebiten.CursorPosition()
	for _, tf := range g.textFields {
		if tf.Contains(x, y) {
			ebiten.SetCursorShape(ebiten.CursorShapeText)
			break
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0xcc, 0xcc, 0xcc, 0xff})
	for _, tf := range g.textFields {
		tf.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Text Input (Ebitengine Demo)")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
