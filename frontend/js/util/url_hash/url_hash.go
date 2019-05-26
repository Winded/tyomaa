package url_hash

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/winded/tyomaa/frontend/js/dom"
)

func Get() string {
	hash := js.Global.Get("window").Get("location").Get("hash").String()
	if len(hash) > 0 {
		return hash[1:]
	} else {
		return ""
	}
}

func Set(value string) {
	js.Global.Get("window").Get("location").Set("hash", value)
}

func AddListener(listener func()) {
	window := js.Global.Get("window")
	dom.JQ(window).On("hashchange", listener)
}
