// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embed"
	"io/fs"

	"goki.dev/gi/v2/gi"
	"goki.dev/gi/v2/gimain"
	"goki.dev/grr"
	"goki.dev/webki"
)

//go:embed content
var content embed.FS

func main() { gimain.Run(app) }

func app() {
	b := gi.NewBody("webki-basic")
	pg := webki.NewPage(b).SetSource(grr.Log(fs.Sub(content, "content")))
	b.AddTopAppBar(pg.TopAppBar)
	pg.OpenURL("", true)
	b.NewWindow().Run().Wait()
}
