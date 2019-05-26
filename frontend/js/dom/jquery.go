package dom

import (
	"github.com/gopherjs/jquery"
)

func JQ(args ...interface{}) jquery.JQuery {
	return jquery.NewJQuery(args...)
}
