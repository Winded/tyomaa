package views

import (
	"time"

	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/models"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/util"
	"github.com/winded/tyomaa/frontend/js/views"
	"github.com/winded/tyomaa/shared/api"
)

type clockView struct {
	views.ViewCore

	context *app.Context

	activeProject string
	activeEntry   api.TimeEntry

	buttonEl   jquery.JQuery
	projectEl  jquery.JQuery
	durationEl jquery.JQuery

	timerChannel chan time.Time
}

func NewClockView(context *app.Context) views.View {
	v := &clockView{
		context: context,
	}

	v.El = dom.JQ(templates.Get("clock"))

	v.buttonEl = dom.JQ("#start-stop-btn", v.El)
	v.projectEl = dom.JQ("#project-name", v.El)
	v.durationEl = dom.JQ("#duration", v.El)

	v.projectEl.On(jquery.CLICK, v.selectProject)

	v.buttonEl.On(jquery.CLICK, v.startStop)

	v.context.Dispatcher.AddListener(models.EVENT_CHANGE_TOKEN, v.tokenChanged)
	v.context.Dispatcher.AddListener(models.EVENT_CLOCK_START, v.clockStarted)
	v.context.Dispatcher.AddListener(models.EVENT_CLOCK_STOP, v.clockStopped)

	v.timerChannel = make(chan time.Time, 5)
	go timerLoop(v.durationEl, v.timerChannel)

	return v
}

func (v *clockView) startStop(event jquery.Event) {
	if v.activeProject == "" {
		v.selectProject(event)
		return
	}

	event.PreventDefault()

	go func() {
		_, err := v.context.Clock.ActiveEntry()

		if err == api.ClockActiveEntryNotFoundErr {
			v.context.Clock.Start(v.activeProject)
		} else if err == nil {
			_, err = v.context.Clock.Stop()
			if err != nil {
				// TODO error handle
			}
		} else {
			// TODO error handle
		}
	}()
}

func timerLoop(durationEl jquery.JQuery, ch chan time.Time) {
	var (
		running     bool
		currentTime time.Time
	)

	timerChan := make(chan bool, 5)
	go func() {
		for {
			time.Sleep(time.Second)
			timerChan <- true
		}
	}()

	for {
		select {
		case <-timerChan:
		case t := <-ch:
			if t != (time.Time{}) {
				currentTime = t
				running = true
			} else {
				running = false
			}
		}

		if running {
			d := time.Since(currentTime)
			durationEl.SetText(util.FormatDuration(d))
		} else {
			durationEl.SetText(util.FormatDuration(0))
		}
	}
}

func (v *clockView) updateTimer() {
	if v.activeProject != "" {
		v.projectEl.SetText(v.activeProject)
	} else {
		v.projectEl.SetHtml("<i>No active project</i>")
	}

	if v.activeEntry.ID != 0 {
		v.timerChannel <- v.activeEntry.Start
		v.setButtonState(true)
	} else {
		v.timerChannel <- time.Time{}
		v.setButtonState(false)
	}
}

func (v *clockView) tokenChanged(data interface{}) {
	token := data.(string)
	if token == "" {
		return
	}

	go func() {
		entry, err := v.context.Clock.ActiveEntry()
		if err == api.ClockActiveEntryNotFoundErr {
			v.activeProject = v.context.Clock.ActiveProject()
			v.updateTimer()
		} else if err == nil {
			v.activeEntry = entry
			v.activeProject = entry.Project
			v.updateTimer()
		} else {
			// TODO error handle
		}
	}()
}

func (v *clockView) clockStarted(data interface{}) {
	entry := data.(api.TimeEntry)

	v.activeProject = entry.Project
	v.activeEntry = entry

	v.updateTimer()
}

func (v *clockView) clockStopped(data interface{}) {
	v.activeEntry = api.TimeEntry{}

	v.updateTimer()
}

func (v *clockView) setButtonState(running bool) {
	if running {
		v.buttonEl.RemoveClass("btn-green")
		v.buttonEl.AddClass("btn-red")
		v.buttonEl.SetHtml(`<i class="fa fa-stop"></i> Stop`)
	} else {
		v.buttonEl.AddClass("btn-green")
		v.buttonEl.RemoveClass("btn-red")
		v.buttonEl.SetHtml(`<i class="fa fa-play"></i> Start`)
	}
}

func (v *clockView) selectProject(event jquery.Event) {
	event.PreventDefault()

	if v.activeEntry.ID != 0 {
		return
	}

	v.context.Dispatcher.Dispatch(models.EVENT_MODAL_OPEN, "project_name")
}
