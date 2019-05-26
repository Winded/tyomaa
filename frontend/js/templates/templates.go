package templates

import (
	"github.com/gopherjs/gopherjs/js"
)

func Get(name string) string {
	return js.Global.Get("TEMPLATES").Get(name).String()
}
