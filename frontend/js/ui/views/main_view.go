package views

import (
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/models"
	"github.com/winded/tyomaa/frontend/js/views"
)

type MainView struct {
	views.ViewCore

	contentView views.ResettableView
	loginView   views.ResettableView
}

func NewMainView(context *app.Context) views.View {
	v := &MainView{}

	v.ViewCore.El = dom.JQ("<div></div>")

	v.contentView = NewContentView(context)
	v.loginView = NewLoginView(context)

	v.ViewCore.El.Append(v.loginView.Element())
	v.ViewCore.El.Append(v.contentView.Element())

	context.Dispatcher.AddListener(models.EVENT_CHANGE_TOKEN, v.tokenChanged)

	return v
}

func (v *MainView) tokenChanged(data interface{}) {
	token := data.(string)

	if token != "" {
		v.loginView.Element().Hide()
		v.contentView.Element().Show()
		v.contentView.Reset()
	} else {
		v.contentView.Element().Hide()
		v.loginView.Element().Show()
		v.loginView.Reset()
	}
}
