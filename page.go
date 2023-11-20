// Copyright (c) 2023, The GoKi Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package webki provides a framework designed for easily building content-focused sites
package webki

import (
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"strings"

	"goki.dev/gi/v2/gi"
	"goki.dev/girl/styles"
	"goki.dev/glide/gidom"
	"goki.dev/glop/dirs"
	"goki.dev/goosi"
	"goki.dev/goosi/events"
	"goki.dev/grr"
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
		s.Direction = styles.Column
	})
}

// OpenURL sets the content of the page from the given url. If the given URL
// has no scheme (eg: "/about"), then it sets the content of the page to the
// file specified by the URL. This is either the "_index.md" file in the
// corresponding directory (eg: "/about/_index.md") or the corresponding
// md file (eg: "/about.md"). If it has a scheme, (eg: "https://example.com"),
// then it opens it in the user's default browser.
func (pg *Page) OpenURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}
	if u.Scheme != "" {
		goosi.TheApp.OpenURL(u.String())
		return nil
	}

	pg.PgURL = rawURL
	pg.History = append(pg.History, rawURL)

	if pg.Source == nil {
		return fmt.Errorf("page source must not be nil")
	}

	// the paths in the fs are never rooted, so we trim a rooted one
	rawURL = strings.TrimPrefix(rawURL, "/")

	fsPath := path.Join(rawURL, "_index.md")
	exists, err := dirs.FileExistsFS(pg.Source, fsPath)
	if err != nil {
		return err
	}
	if !exists {
		fsPath = path.Clean(rawURL) + ".md"
	}

	b, err := fs.ReadFile(pg.Source, fsPath)
	if err != nil {
		return err
	}

	updt := pg.UpdateStart()
	pg.DeleteChildren(true)
	err = gidom.ReadMD(pg, pg, b)
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

// TopAppBar is the default [gi.TopAppBar] for a [Page]
func (pg *Page) TopAppBar(tb *gi.TopAppBar) {
	gi.DefaultTopAppBarStd(tb)

	back := tb.ChildByName("back").(*gi.Button)
	back.OnClick(func(e events.Event) {
		if len(pg.History) > 1 {
			pg.OpenURL(pg.History[len(pg.History)-2])
		}
	})

	ch := tb.ChildByName("nav-bar").(*gi.Chooser)
	ch.AllowNew = true
	ch.ItemsFunc = func() {
		ch.Items = make([]any, len(pg.History))
		for i, u := range pg.History {
			// we reverse the order
			ch.Items[len(pg.History)-i-1] = u
		}
	}
	ch.OnChange(func(e events.Event) {
		grr.Log0(pg.OpenURL(ch.CurLabel))
		e.SetHandled()
	})
}
