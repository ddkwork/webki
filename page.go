// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package webki provides a framework designed for easily building content-focused sites
package glide

import (
	"fmt"
	"io/fs"
	"path"

	"goki.dev/gi/v2/gi"
	"goki.dev/girl/styles"
	"goki.dev/glide/gidom"
	"goki.dev/glop/dirs"
	"goki.dev/ki/v2"
)

// Page represents one site page
type Page struct {
	gi.Frame
	gidom.ContextBase

	// Source is the filesystem in which the content is located.
	Source fs.FS

	// The history of URLs that have been visited. The oldest page is first.
	History []string

	// PgURL is the current page URL
	PgURL string
}

var _ ki.Ki = (*Page)(nil)

func (pg *Page) OnInit() {
	pg.Frame.OnInit()
	pg.Style(func(s *styles.Style) {
		s.Direction = styles.Col
	})
}

// OpenURL sets the content of the page from the given url.
func (pg *Page) OpenURL(url string) error {
	pg.PgURL = url
	pg.History = append(pg.History, url)

	if pg.Source == nil {
		return fmt.Errorf("page source must not be nil")
	}

	fsPath := path.Join(url, "_index.md")
	exists, err := dirs.FileExistsFS(pg.Source, fsPath)
	if err != nil {
		return err
	}
	if !exists {
		fsPath = path.Clean(url) + ".md"
	}

	file, err := pg.Source.Open(fsPath)
	if err != nil {
		return err
	}
	defer file.Close()

	updt := pg.UpdateStart()
	pg.DeleteChildren(true)
	err = gidom.ReadHTML(pg, pg, file)
	if err != nil {
		return err
	}
	pg.Update()
	pg.UpdateEndLayout(updt)
	return nil
}

// PageURL returns the current page URL
func (pg *Page) PageURL() string {
	return pg.PgURL
}
