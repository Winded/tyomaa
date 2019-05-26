package widgets

import (
	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/models"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/views"
)

type ActionFunc func() bool

type ActionCancelWidget struct {
	views.ViewCore
	context *app.Context

	loadingEl jquery.JQuery

	actionFunc       ActionFunc
	performingAction bool
}

func NewActionCancelWidget(context *app.Context, actionFunc ActionFunc, actionText string, dangerous bool) *ActionCancelWidget {
	w := &ActionCancelWidget{
		context:    context,
		actionFunc: actionFunc,
	}

	w.El = dom.JQ(templates.Get("action_cancel_widget"))
	w.loadingEl = dom.JQ("#loading", w.El)
	w.loadingEl.Hide()

	actionBtn := dom.JQ("#action-btn", w.El)
	actionBtn.On(jquery.CLICK, w.Action)
	dom.JQ("#action-text", actionBtn).SetText(actionText)
	if dangerous {
		actionBtn.AddClass("btn-red")
		w.loadingEl.AddClass("btn-red")
	} else {
		actionBtn.AddClass("btn-green")
		w.loadingEl.AddClass("btn-green")
	}

	cancelBtn := dom.JQ("#cancel-btn", w.El)
	cancelBtn.On(jquery.CLICK, w.Cancel)

	return w
}

func NewActionOnlyWidget(context *app.Context, actionFunc ActionFunc, actionText string, dangerous bool) *ActionCancelWidget {
	w := NewActionCancelWidget(context, actionFunc, actionText, dangerous)
	dom.JQ("#cancel-btn", w.El).Remove()
	return w
}

func (w *ActionCancelWidget) Action(event jquery.Event) {
	event.PreventDefault()

	if w.performingAction {
		return
	}

	go func() {
		w.performingAction = true

		w.loadingEl.Show()
		closeModal := w.actionFunc()
		w.loadingEl.Hide()

		if closeModal {
			w.context.Dispatcher.Dispatch(models.EVENT_MODAL_CLOSE, nil)
		}

		w.performingAction = false
	}()
}

func (w *ActionCancelWidget) Cancel(event jquery.Event) {
	event.PreventDefault()
	w.context.Dispatcher.Dispatch(models.EVENT_MODAL_CLOSE, nil)
}
