package views

import (
	"log"

	"github.com/winded/tyomaa/frontend/js/models"
	"github.com/winded/tyomaa/frontend/js/views"

	"github.com/winded/tyomaa/frontend/js/util"

	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/shared/api"
)

type EntryRowView struct {
	views.ViewCore
	context *app.Context

	projectEl, startEl, endEl jquery.JQuery

	editEl, deleteEl jquery.JQuery

	data *api.TimeEntry
}

func NewEntryRowView(context *app.Context) *EntryRowView {
	v := &EntryRowView{
		context: context,
	}

	v.El = dom.JQ(`
		<tr>
			<td></td>
			<td></td>
			<td></td>
			<td class="icon-column"><a id="edit-btn" title="Edit" class="fa fa-edit" href="#"></a><a id="delete-btn" title="Delete" class="fa fa-trash" href="#"></a></td>
		</tr>
	`)

	cols := dom.JQ("td", v.El)
	v.projectEl = cols.Eq(0)
	v.startEl = cols.Eq(1)
	v.endEl = cols.Eq(2)
	links := cols.Eq(3)

	v.editEl = dom.JQ("#edit-btn", links)
	v.editEl.On(jquery.CLICK, v.edit)
	v.deleteEl = dom.JQ("#delete-btn", links)
	v.deleteEl.On(jquery.CLICK, v.del)

	return v
}

func (v *EntryRowView) SetData(data *api.TimeEntry) {
	v.data = data

	if v.data != nil {
		v.projectEl.SetText(v.data.Project)
		v.startEl.SetText(util.FormatDate(v.data.Start))
		if v.data.End != nil {
			v.endEl.SetText(util.FormatDate(*v.data.End))
		} else {
			v.endEl.SetText("")
		}
	} else {
		v.projectEl.SetText("")
		v.startEl.SetText("")
		v.endEl.SetText("")
	}
}

func (v *EntryRowView) Remove() {
	v.El.Remove()
}

func (v *EntryRowView) edit(event jquery.Event) {
	event.PreventDefault()

	if v.data == nil {
		return
	}

	v.context.Dispatcher.Dispatch(models.EVENT_ENTRY_SELECT, *v.data)
	v.context.Dispatcher.Dispatch(models.EVENT_MODAL_OPEN, "edit_entry")
}

func (v *EntryRowView) del(event jquery.Event) {
	event.PreventDefault()

	if v.data == nil {
		return
	}

	v.context.Dispatcher.Dispatch(models.EVENT_ENTRY_SELECT, *v.data)
	v.context.Dispatcher.Dispatch(models.EVENT_MODAL_OPEN, "delete_entry")
}

type EntryTableView struct {
	views.ViewCore

	context *app.Context

	tableEl jquery.JQuery
	tbodyEl jquery.JQuery

	data []api.TimeEntry
	rows []*EntryRowView
}

func NewEntryTableView(context *app.Context) *EntryTableView {
	v := &EntryTableView{
		context: context,
	}

	v.context = context

	v.El = dom.JQ(templates.Get("entry_table"))

	v.tableEl = dom.JQ("#entry-table", v.El)
	v.tbodyEl = dom.JQ("tbody", v.tableEl)

	dom.JQ("#add-entry-btn", v.El).On(jquery.CLICK, v.newEntry)

	v.context.Dispatcher.AddListener(models.EVENT_CLOCK_START, v.doRefresh)
	v.context.Dispatcher.AddListener(models.EVENT_CLOCK_STOP, v.doRefresh)
	v.context.Dispatcher.AddListener(models.EVENT_ENTRY_CREATED, v.doRefresh)
	v.context.Dispatcher.AddListener(models.EVENT_ENTRY_UPDATED, v.doRefresh)
	v.context.Dispatcher.AddListener(models.EVENT_ENTRY_DELETED, v.doRefresh)

	return v
}

func (v *EntryTableView) Reset() {
	v.refresh()
}

func (v *EntryTableView) setData(data []api.TimeEntry) {
	v.data = data

	for _, entry := range v.data {
		row := NewEntryRowView(v.context)
		row.SetData(&entry)
		v.rows = append(v.rows, row)

		v.tbodyEl.Append(row.Element())
	}
}

func (v *EntryTableView) newEntry(event jquery.Event) {
	event.PreventDefault()
	v.context.Dispatcher.Dispatch(models.EVENT_MODAL_OPEN, "new_entry")
}

func (v *EntryTableView) doRefresh(data interface{}) {
	go v.refresh()
}

func (v *EntryTableView) refresh() {
	for _, row := range v.rows {
		row.Remove()
	}
	v.rows = []*EntryRowView{}

	entries, err := v.context.Entries.List()
	if err != nil {
		log.Println(err)
		return
	}

	v.setData(entries)
}
