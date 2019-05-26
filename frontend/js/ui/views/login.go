package views

import (
	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/ui/views/widgets"
	"github.com/winded/tyomaa/frontend/js/views"
)

type LoginView struct {
	views.ViewCore
	context *app.Context

	formEl     jquery.JQuery
	usernameEl jquery.JQuery
	passwordEl jquery.JQuery

	errorView  *widgets.ErrorMessageWidget
	actionView *widgets.ActionCancelWidget
}

func NewLoginView(context *app.Context) *LoginView {
	v := &LoginView{
		context: context,
	}

	v.ViewCore.El = dom.JQ(templates.Get("login"))

	v.formEl = dom.JQ("#login-form", v.ViewCore.El)
	v.usernameEl = dom.JQ("input[name='username']", v.formEl)
	v.passwordEl = dom.JQ("input[name='password']", v.formEl)

	v.errorView = widgets.NewErrorMessageWidget()
	v.formEl.Prepend(v.errorView.Element())

	v.actionView = widgets.NewActionOnlyWidget(context, v.submit, "Login", false)
	v.formEl.Append(v.actionView.Element())
	v.formEl.Submit(v.actionView.Action)

	return v
}

func (v *LoginView) Reset() {
	v.usernameEl.SetVal("")
	v.passwordEl.SetVal("")
	v.errorView.Hide()
}

func (v *LoginView) submit() bool {
	if err := v.context.Auth.Login(v.usernameEl.Val(), v.passwordEl.Val()); err != nil {
		v.errorView.Show(err)
		return false
	}

	return true
}
