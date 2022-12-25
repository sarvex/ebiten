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

//go:build !ios

package textinput

// #cgo CFLAGS: -x objective-c
// #cgo LDFLAGS: -framework Cocoa
//
// #import <Cocoa/Cocoa.h>
//
// static NSTextView* textView() {
//   static NSTextView* textView;
//   if (!textView) {
//     textView = [[NSTextView alloc] initWithFrame:NSMakeRect(0, 0, 100, 10)];
//     [textView setFont:[NSFont userFontOfSize:0.01]];
//     [textView setRichText:NO];
//     [textView setUsesFontPanel:NO];
//     [textView setUsesFindPanel:NO];
//     [textView setDisplaysLinkToolTips:NO];
//     [textView setHorizontallyResizable:NO];
//     [textView setVerticallyResizable:NO];
//     //[textView setBackgroundColor:[NSColor clearColor]];
//   }
//   return textView;
// }
//
// static void startTextInput() {
//   NSTextView* tv = textView();
//   NSWindow* window = [[NSApplication sharedApplication] mainWindow];
//   [[window contentView] addSubview: tv];
//   [window makeFirstResponder: tv];
// }
//
// static void endTextInput() {
//   NSTextView* tv = textView();
//   [tv removeFromSuperview];
// }
//
// static void setTextInputPosition(int x, int y) {
//   NSWindow* window = [[NSApplication sharedApplication] mainWindow];
//   y = [[window contentView] frame].size.height - y - 10 - 4;
//   NSTextView* tv = textView();
//   [tv setFrame:NSMakeRect(x, y, 100, 10)];
// }
//
// static void textInputStatus(const char** text, int* start, int* end) {
//   NSTextView* tv = textView();
//   *text = [[tv string] UTF8String];
//   NSRange range = [tv selectedRange];
//   *start = range.location;
//   *end = range.location + range.length;
// }
//
// static void setTextInputStatus(const char* text, int start, int end) {
//   NSTextView* tv = textView();
//   [tv setString: [NSString stringWithUTF8String: text]];
//   [tv setSelectedRange: NSMakeRange(start, end-start)];
// }
import "C"

import (
	"unsafe"

	"github.com/hajimehoshi/ebiten/v2/internal/ui"
)

func startTextInput() {
	ui.RunOnMainThread(func() {
		C.startTextInput()
	})
}

func endTextInput() {
	ui.RunOnMainThread(func() {
		C.endTextInput()
	})
}

func setTextInputPosition(x, y int) {
	xf, yf := ui.LogicalPositionToClientPosition(float64(x), float64(y))
	ui.RunOnMainThread(func() {
		C.setTextInputPosition(C.int(xf), C.int(yf))
	})
}

func textInputStatus() (string, int, int) {
	var text string
	var start, end int
	ui.RunOnMainThread(func() {
		var ctext *C.char
		var cstart, cend C.int
		C.textInputStatus(&ctext, &cstart, &cend)
		text = C.GoString(ctext)
		start = int(cstart)
		end = int(cend)
	})
	return text, convertUTF16CountToByteCount(text, start), convertUTF16CountToByteCount(text, end)
}

func setTextInputStatus(text string, start, end int) {
	cstart := C.int(convertByteCountToUTF16Count(text, start))
	cend := C.int(convertByteCountToUTF16Count(text, end))
	ctext := C.CString(text)
	defer C.free(unsafe.Pointer(ctext))
	ui.RunOnMainThread(func() {
		C.setTextInputStatus(ctext, cstart, cend)
	})
}
