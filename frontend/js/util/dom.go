package util

import "github.com/gopherjs/gopherjs/js"

func Prompt(label, defaultValue string) string {
	return js.Global.Get("window").Call("prompt", label, defaultValue).String()
}
