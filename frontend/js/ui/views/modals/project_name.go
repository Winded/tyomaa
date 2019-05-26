package modals

import (
	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/ui/views/widgets"
	"github.com/winded/tyomaa/frontend/js/views"
)

type projectNameModal struct {
	views.ViewCore
	context *app.Context

	errorView  *widgets.ErrorMessageWidget
	actionView *widgets.ActionCancelWidget

	formEl   jquery.JQuery
	inputEl  jquery.JQuery
	cancelEl jquery.JQuery
}

func NewProjectNameModal(context *app.Context) views.ResettableView {
	v := &projectNameModal{
		context: context,
	}

	v.El = dom.JQ(templates.Get("project_name_modal"))
	v.formEl = dom.JQ("form", v.El)
	v.inputEl = dom.JQ("input[name=project-name]", v.formEl)
	v.cancelEl = dom.JQ("#cancel-btn", v.formEl)

	v.errorView = widgets.NewErrorMessageWidget()
	v.formEl.Prepend(v.errorView.Element())

	v.actionView = widgets.NewActionCancelWidget(context, v.submit, "Start", false)
	v.formEl.Append(v.actionView.Element())
	v.formEl.Submit(v.actionView.Action)

	return v
}

func (v *projectNameModal) Reset() {
	v.inputEl.SetVal("")
	v.errorView.Hide()
}

func (v *projectNameModal) submit() bool {
	projectName := v.inputEl.Val()
	if projectName == "" {
		return false
	}

	_, err := v.context.Clock.Start(projectName)
	if err != nil {
		v.errorView.Show(err)
		return false
	}

	return true
}
