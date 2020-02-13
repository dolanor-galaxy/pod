package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"image/color"
)

var (
	buttonDialogCancel = new(controller.Button)
	buttonDialogOK     = new(controller.Button)
	buttonDialogClose  = new(controller.Button)

	list               = &layout.List{
		Axis: layout.Vertical,
	}
)

// Main wallet screen
func (ui *DuoUI) DuoUIdialog() {
	cs := ui.ly.Context.Constraints
	theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, "ee000000", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
	layout.Align(layout.Center).Layout(ui.ly.Context, func() {
		//cs := ui.ly.Context.Constraints
		theme.DuoUIdrawRectangle(ui.ly.Context, 408, 150, ui.ly.Theme.Color.Primary, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})

		layout.Flex{
			Axis:      layout.Vertical,
			Alignment: layout.Middle,
		}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(ui.ly.Context,
					layout.Rigid(func() {
						layout.Align(layout.Center).Layout(ui.ly.Context, func() {
							layout.Inset{Top: unit.Dp(24), Bottom: unit.Dp(8), Left: unit.Dp(0), Right: unit.Dp(4)}.Layout(ui.ly.Context, func() {
								cur := ui.ly.Theme.H4(ui.rc.Dialog.Text)
								cur.Font.Typeface = "bariol"
								cur.Color = color.RGBA{A: 0xff, R: 0xcf, G: 0xcf, B: 0xcf}
								cur.Alignment = text.Start
								cur.Layout(ui.ly.Context)
							})
						})
					}),
				)
			}),
			layout.Rigid(func() {
				layout.Flex{
					Axis:      layout.Horizontal,
					Alignment: layout.Middle,
				}.Layout(ui.ly.Context,
					layout.Rigid(ui.dialogButon(func(){ui.rc.Dialog.Cancel()},"CANCEL", "ffcfcfcf", "ffcf3030", "ffcfcfcf", buttonDialogCancel, ui.ly.Theme.Icons["iconCancel"])),
					layout.Rigid(ui.dialogButon(func(){ui.rc.Dialog.Ok()},"OK", "ffcfcfcf", "ff308030", "ffcfcfcf", buttonDialogOK, ui.ly.Theme.Icons["iconOK"])),
					layout.Rigid(ui.dialogButon(func(){ui.rc.Dialog.Show = false},"CLOSE", "ffcfcfcf", "ffcf8030", "ffcfcfcf", buttonDialogClose, ui.ly.Theme.Icons["iconClose"])),
				)
			}),
		)
	})
}

func (ui *DuoUI)dialogButon(f func(),t, txtColor, bgColor, iconColor string, button *controller.Button, icon *theme.DuoUIicon) func() {
	var b theme.DuoUIbutton
	return func() {
		layout.Inset{Top: unit.Dp(8), Bottom: unit.Dp(8), Left: unit.Dp(8), Right: unit.Dp(8)}.Layout(ui.ly.Context, func() {
			b = ui.ly.Theme.DuoUIbutton(t, txtColor, bgColor, iconColor, 24, 120, 60, 0, 0, icon)
			for button.Clicked(ui.ly.Context) {
				f()
			}
			b.Layout(ui.ly.Context, button)
		})
	}
}
