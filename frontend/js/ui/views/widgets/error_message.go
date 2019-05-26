package widgets

import (
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/views"
	"github.com/winded/tyomaa/shared/api"
)

type ErrorMessageWidget struct {
	views.ViewCore
}

func NewErrorMessageWidget() *ErrorMessageWidget {
	w := &ErrorMessageWidget{}

	w.El = dom.JQ(templates.Get("error_message_widget"))
	w.El.Hide()

	return w
}

func (w *ErrorMessageWidget) Show(err error) {
	if apiErr, ok := err.(api.ApiError); ok {
		w.El.SetText(apiErr.Message)
	} else {
		w.El.SetText(err.Error())
	}
	w.El.Show()
}

func (w *ErrorMessageWidget) Hide() {
	w.El.Hide()
}
