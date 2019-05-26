package widgets

import (
	"strconv"
	"time"

	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/util"
	"github.com/winded/tyomaa/frontend/js/views"
)

const (
	TAB_KEY   = 9
	SHIFT_KEY = 16
)

type DateTimeWidget struct {
	views.ViewCore

	day, month, year     jquery.JQuery
	hour, minute, second jquery.JQuery

	dateSeparators jquery.JQuery
	timeSeparators jquery.JQuery
}

func NewDateTimeWidget() *DateTimeWidget {
	v := &DateTimeWidget{}

	v.El = dom.JQ(templates.Get("datetime_widget"))

	v.day = dom.JQ("input[name=day]", v.El)
	v.month = dom.JQ("input[name=month]", v.El)
	v.year = dom.JQ("input[name=year]", v.El)

	v.hour = dom.JQ("input[name=hour]", v.El)
	v.minute = dom.JQ("input[name=minute]", v.El)
	v.second = dom.JQ("input[name=second]", v.El)

	v.dateSeparators = dom.JQ("#date-separator", v.El)
	v.timeSeparators = dom.JQ("#time-separator", v.El)

	return v
}

func (v *DateTimeWidget) Reset() {
	v.day.SetVal("")
	v.month.SetVal("")
	v.year.SetVal("")
	v.hour.SetVal("")
	v.minute.SetVal("")
	v.second.SetVal("")
}

func (v *DateTimeWidget) SetSeparators(dateSeparator, timeSeparator string) {
	v.dateSeparators.SetText(dateSeparator)
	v.timeSeparators.SetText(timeSeparator)
}

func (v *DateTimeWidget) Time() (time.Time, error) {
	day, err := strconv.Atoi(v.day.Val())
	if err != nil {
		return time.Time{}, err
	}
	month, err := strconv.Atoi(v.month.Val())
	if err != nil {
		return time.Time{}, err
	}
	year, err := strconv.Atoi(v.year.Val())
	if err != nil {
		return time.Time{}, err
	}
	hour, err := strconv.Atoi(v.hour.Val())
	if err != nil {
		return time.Time{}, err
	}
	minute, err := strconv.Atoi(v.minute.Val())
	if err != nil {
		return time.Time{}, err
	}
	second, err := strconv.Atoi(v.second.Val())
	if err != nil {
		return time.Time{}, err
	}

	t := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
	return t.UTC(), nil
}

func (v *DateTimeWidget) SetTime(value time.Time) {
	value = value.Local()

	setVal := func(element jquery.JQuery, value int) {
		element.SetVal(util.LeftPad(strconv.Itoa(value), "0", 2))
	}

	setVal(v.day, value.Day())
	setVal(v.month, int(value.Month()))
	setVal(v.year, value.Year())

	setVal(v.hour, value.Hour())
	setVal(v.minute, value.Minute())
	setVal(v.second, value.Second())
}
