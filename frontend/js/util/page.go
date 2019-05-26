package util

import (
	"github.com/gopherjs/gopherjs/js"
)

func SetPageTitle(title string) {
	js.Global.Get("document").Set("title", title)
}
