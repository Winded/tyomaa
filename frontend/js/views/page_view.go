package views

import (
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/models"
)

// ResettableView represents a View that can be reset to it's original state by calling Reset()
type ResettableView interface {
	View
	Reset()
}

type PageMap map[string]ResettableView

type PageView struct {
	ViewCore

	routes      PageMap
	currentView ResettableView
}

func NewPageView(routes PageMap, context *app.Context) *PageView {
	if routes == nil {
		routes = make(PageMap)
	}

	v := &PageView{}

	v.ViewCore.El = dom.JQ("<div></div>")
	v.routes = routes

	context.Dispatcher.AddListener(models.EVENT_CHANGE_PAGE, v.onPageChange)

	return v
}

func (this *PageView) onPageChange(data interface{}) {
	newPage := data.(string)
	newView := this.routes[newPage]

	if newView == nil {
		newView = this.routes["*"]
		if newView == nil {
			if this.currentView != nil {
				this.currentView.Element().Detach()
				this.currentView = nil
			}
			return
		}
	}

	if this.currentView != nil && this.currentView == newView {
		return
	} else if this.currentView != nil {
		this.currentView.Element().Detach()
		this.currentView = nil
	}

	this.Element().Append(newView.Element())
	this.currentView = newView

	// Run reset in non-event goroutine
	go func() {
		newView.Reset()
	}()
}
