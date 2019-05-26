package views

import (
	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/models"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/views"
)

type NavView struct {
	views.ViewCore
	context *app.Context

	clockView views.View

	userEl jquery.JQuery
}

func NewNavView(context *app.Context) *NavView {
	v := &NavView{
		context: context,
	}

	v.El = dom.JQ(templates.Get("nav"))

	v.clockView = NewClockView(context)
	dom.JQ("#clock-container", v.El).Append(v.clockView.Element())

	dom.JQ("#link-logout", v.El).On(jquery.CLICK, v.logout)

	v.userEl = dom.JQ("#user-label", v.El)
	context.Dispatcher.AddListener(models.EVENT_CHANGE_TOKEN, v.tokenChanged)

	return v
}

func (v *NavView) logout(event jquery.Event) {
	event.PreventDefault()
	v.context.Auth.Logout()
}

func (v *NavView) tokenChanged(data interface{}) {
	token := data.(string)
	if token == "" {
		return
	}

	go func() {
		user, err := v.context.Auth.GetUser()
		if err != nil {
			// TODO
			return
		}

		v.userEl.SetText(user.Name)
	}()
}
