package modals

import (
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/models"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/ui/views/widgets"
	"github.com/winded/tyomaa/frontend/js/views"
	"github.com/winded/tyomaa/shared/api"
)

type deleteEntryModal struct {
	views.ViewCore
	context *app.Context

	errorView  *widgets.ErrorMessageWidget
	actionView *widgets.ActionCancelWidget

	entry api.TimeEntry
}

func NewDeleteEntryModal(context *app.Context) views.ResettableView {
	v := &deleteEntryModal{
		context: context,
	}

	v.El = dom.JQ(templates.Get("delete_entry_modal"))

	formEl := dom.JQ("form", v.El)

	v.errorView = widgets.NewErrorMessageWidget()
	formEl.Prepend(v.errorView.Element())

	v.actionView = widgets.NewActionCancelWidget(context, v.submit, "Delete", true)
	formEl.Append(v.actionView.Element())
	formEl.Submit(v.actionView.Action)

	v.context.Dispatcher.AddListener(models.EVENT_ENTRY_SELECT, v.entrySelect)

	return v
}

func (v *deleteEntryModal) Reset() {
	v.errorView.Hide()
}

func (v *deleteEntryModal) entrySelect(data interface{}) {
	v.entry = data.(api.TimeEntry)
}

func (v *deleteEntryModal) submit() bool {
	if v.entry.ID == 0 {
		return false
	}

	_, err := v.context.Entries.Delete(v.entry.ID)
	if err != nil {
		v.errorView.Show(err)
		return false
	}

	return true
}
