package views

import (
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/views"
)

type ContentView struct {
	views.ViewCore
	*app.Context

	nav         views.View
	entriesView views.ResettableView
}

func NewContentView(context *app.Context) *ContentView {
	v := &ContentView{
		Context: context,
	}

	v.ViewCore.El = dom.JQ(templates.Get("content"))

	v.nav = NewNavView(context)
	v.entriesView = NewEntryTableView(context)

	v.ViewCore.El.Append(v.nav.Element())
	v.ViewCore.El.Append(v.entriesView.Element())

	return v
}

func (v *ContentView) Reset() {
	v.entriesView.Reset()
}
