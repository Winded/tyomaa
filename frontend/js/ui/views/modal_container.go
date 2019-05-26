package views

import (
	"github.com/gopherjs/jquery"
	"github.com/winded/tyomaa/frontend/js/app"
	"github.com/winded/tyomaa/frontend/js/dom"
	"github.com/winded/tyomaa/frontend/js/models"
	"github.com/winded/tyomaa/frontend/js/templates"
	"github.com/winded/tyomaa/frontend/js/views"
)

type ModalContainer struct {
	views.ViewCore
	context *app.Context

	containerEl jquery.JQuery

	modalViews   map[string]views.ResettableView
	currentModal views.ResettableView
}

func NewModalContainer(context *app.Context, modalViews map[string]views.ResettableView) views.View {
	v := &ModalContainer{
		context:    context,
		modalViews: modalViews,
	}

	v.El = dom.JQ(templates.Get("modal_container"))
	v.containerEl = dom.JQ(".modal-container", v.El)

	v.context.Dispatcher.AddListener(models.EVENT_MODAL_OPEN, v.openModal)
	v.context.Dispatcher.AddListener(models.EVENT_MODAL_CLOSE, v.closeModal)

	v.Element().Hide()
	return v
}

func (v *ModalContainer) openModal(data interface{}) {
	modalName := data.(string)

	view := v.modalViews[modalName]
	if view == nil {
		// TODO error handling
		return
	}

	if v.currentModal != nil {
		v.closeModal(nil)
	}

	v.containerEl.Append(view.Element())
	view.Reset()
	v.Element().Show()
	v.currentModal = view
}

func (v *ModalContainer) closeModal(data interface{}) {
	if v.currentModal == nil {
		return
	}

	v.currentModal.Element().Detach()
	v.Element().Hide()
	v.currentModal = nil
}
