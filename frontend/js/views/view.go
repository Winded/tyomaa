package views

import (
	"github.com/gopherjs/jquery"
)

// Views are a programmatic representation of a view linked to a DOM element.
type View interface {
	// Element returns the underlying JQuery DOM element reference
	Element() jquery.JQuery
}

type ViewCore struct {
	El jquery.JQuery
}

func (v *ViewCore) Element() jquery.JQuery {
	return v.El
}
