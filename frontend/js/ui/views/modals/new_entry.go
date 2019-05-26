package modals

import (
	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/ui/views/widgets"
	"github.com/winded/tyomaa/frontend/js/views"
	"github.com/winded/tyomaa/shared/api"
)

type newEntryModal struct {
	views.ViewCore
	context *app.Context

	formEl    jquery.JQuery
	projectEl jquery.JQuery
	startView *widgets.DateTimeWidget
	endView   *widgets.DateTimeWidget

	errorView  *widgets.ErrorMessageWidget
	actionView *widgets.ActionCancelWidget
}

func NewNewEntryModal(context *app.Context) views.ResettableView {
	v := &newEntryModal{
		context: context,
	}

	v.El = dom.JQ(templates.Get("new_entry_modal"))

	v.formEl = dom.JQ("form", v.El)

	v.projectEl = dom.JQ("input[name=project]", v.formEl)

	v.startView = widgets.NewDateTimeWidget()
	dom.JQ("#start-container", v.formEl).Append(v.startView.Element())
	v.endView = widgets.NewDateTimeWidget()
	dom.JQ("#end-container", v.formEl).Append(v.endView.Element())

	v.errorView = widgets.NewErrorMessageWidget()
	v.formEl.Prepend(v.errorView.Element())

	v.actionView = widgets.NewActionCancelWidget(context, v.submit, "Create", false)
	v.formEl.Append(v.actionView.Element())
	v.formEl.Submit(v.actionView.Action)

	return v
}

func (v *newEntryModal) Reset() {
	v.projectEl.SetVal("")
	v.startView.Reset()
	v.endView.Reset()
	v.errorView.Hide()
}

func (v *newEntryModal) submit() bool {
	project := v.projectEl.Val()
	startTime, err := v.startView.Time()
	if err != nil {
		v.errorView.Show(err)
		return false
	}
	endTime, err := v.endView.Time()
	if err != nil {
		v.errorView.Show(err)
		return false
	}

	_, err = v.context.Entries.Create(api.TimeEntry{
		Project: project,
		Start:   startTime,
		End:     &endTime,
	})
	if err != nil {
		v.errorView.Show(err)
		return false
	}

	return true
}
