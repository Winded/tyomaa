package main

import (
	"time"

	"github.com/winded/tyomaa/frontend/js/dispatcher"
	"github.com/winded/tyomaa/frontend/js/dom"
	uiviews "github.com/winded/tyomaa/frontend/js/ui/views"
	"github.com/winded/tyomaa/frontend/js/ui/views/modals"
	"github.com/winded/tyomaa/frontend/js/views"

	"github.com/gopherjs/gopherjs/js"
	"github.com/winded/tyomaa/frontend/js/app"
	_ "github.com/winded/tyomaa/frontend/js/inc"
	"github.com/winded/tyomaa/frontend/js/util/xhr_http_client"
	"github.com/winded/tyomaa/shared/api/client"
)

func main() {
	apiEndpoint := js.Global.Get("window").Get("API_ENDPOINT").String()

	apiClient := client.NewApiClient(xhr_http_client.New(time.Second*10), client.Settings{
		Host: apiEndpoint,
	})

	context := app.NewContext(dispatcher.New(), apiClient)

	mainView := uiviews.NewMainView(context)
	modalContainer := uiviews.NewModalContainer(context, map[string]views.ResettableView{
		"project_name": modals.NewProjectNameModal(context),
		"new_entry":    modals.NewNewEntryModal(context),
		"edit_entry":   modals.NewEditEntryModal(context),
		"delete_entry": modals.NewDeleteEntryModal(context),
	})

	appEl := dom.JQ("#app")
	appEl.Append(mainView.Element())
	appEl.Append(modalContainer.Element())

	context.Auth.LoadToken()

	dom.JQ("#loading-indicator").Hide()
}
