package duoui

import (
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/cmd/gui/helpers"
	"github.com/p9c/pod/cmd/gui/models"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gui/layout"
)

func DuoUIoverview(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.DuoUIcontext, func() {
		viewport := layout.Flex{Axis: layout.Horizontal}
		if duo.DuoUIcontext.Constraints.Width.Max < 1024 {
			viewport = layout.Flex{Axis: layout.Vertical}
		}
		page := duo.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 500, 500, 10, 10)
		page.Layout(duo.DuoUIcontext, func() {
			viewport.Layout(duo.DuoUIcontext,
				layout.Flexed(0.5, func() {
					cs := duo.DuoUIcontext.Constraints
					helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ffcfcfcf", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					components.DuoUIbalanceWidget(duo, rc)
				}),
				layout.Flexed(0.5, func() {
					cs := duo.DuoUIcontext.Constraints
					helpers.DuoUIdrawRectangle(duo.DuoUIcontext, cs.Width.Max, cs.Height.Max, "ff424242", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
					components.DuoUIlatestTxsWidget(duo, cx, rc)
				}),
			)
		})
	}
}
func DuoUIsend(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.DuoUIcontext, func() {
		page := duo.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 500, 500, 10, 10)
		page.Layout(duo.DuoUIcontext, func() {
			components.DuoUIsend(duo, cx, rc)
		})
	}
}
func DuoUIreceive(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.DuoUIcontext, func() {
		page := duo.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 500, 500, 10, 10)
		page.Layout(duo.DuoUIcontext, func() {
			duo.DuoUItheme.H5("receive :").Layout(duo.DuoUIcontext)
		})
	}
}

func DuoUIaddressbook(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.DuoUIcontext, func() {
		page := duo.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 500, 500, 10, 10)
		page.Layout(duo.DuoUIcontext, func() {
			duo.DuoUItheme.H5("addressbook :").Layout(duo.DuoUIcontext)
		})
	}
}
func DuoUIsettings(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.DuoUIcontext, func() {
		page := duo.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 500, 500, 0, 0)
		page.Layout(duo.DuoUIcontext, func() {
			components.DuoUIsettingsWidget(duo, cx, rc)
		})
	}
}

func DuoUInetwork(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.DuoUIcontext, func() {
		page := duo.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 500, 500, 10, 10)
		page.Layout(duo.DuoUIcontext, func() {
			duo.DuoUItheme.H5("network :").Layout(duo.DuoUIcontext)
		})
	}
}

func DuoUIhistory(duo *models.DuoUI, cx *conte.Xt, rc *rcd.RcVar) (*layout.Context, func()) {
	return duo.DuoUIcontext, func() {
		page := duo.DuoUItheme.DuoUIpage("ffcf30cf", "ffcf3030", 500, 500, 10, 10)
		page.Layout(duo.DuoUIcontext, func() {
			components.DuoUItransactionsWidget(duo, cx, rc)
		})
	}
}