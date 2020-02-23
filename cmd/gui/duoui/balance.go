package duoui

import (
	"fmt"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
	"golang.org/x/exp/shiny/materialdesign/icons"
	"image"
)

var (
	itemsList = &layout.List{
		Axis: layout.Vertical,
	}
	singleItem = &layout.Flex{
		Axis:    layout.Horizontal,
		Spacing: layout.SpaceBetween,
	}
	Icon, _ = theme.NewDuoUIicon(icons.EditorMonetizationOn)
)

func (ui *DuoUI) DuoUIbalance() func() {
	return func() {
		cs := ui.ly.Context.Constraints
		theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, cs.Height.Max, ui.ly.Theme.Color.Light, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
		in := layout.UniformInset(unit.Dp(16))
		in.Layout(ui.ly.Context, func() {
			cs := ui.ly.Context.Constraints
			navButtons := []func(){
				listItem(ui.ly.Context, ui.ly.Theme, "BALANCE :", ui.rc.Status.Wallet.Balance+" "+ui.rc.Settings.Abbrevation),
				func() { theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, "ffbdbdbd", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0}) },
				listItem(ui.ly.Context, ui.ly.Theme, "UNCNFIRMED :", ui.rc.Status.Wallet.Unconfirmed+" "+ui.rc.Settings.Abbrevation),
				func() { theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, "ffbdbdbd", [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0}) },
				listItem(ui.ly.Context, ui.ly.Theme, "TRANSACTIONS :", fmt.Sprint(ui.rc.Status.Wallet.TxsNumber)),
			}
			itemsList.Layout(ui.ly.Context, len(navButtons), func(i int) {
				layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, navButtons[i])
			})
		})
	}
}

func listItem(gtx *layout.Context, th *theme.DuoUItheme, name, value string) func() {
	return func() {
		layout.Flex{
			Axis:    layout.Horizontal,
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			layout.Rigid(func() {
				layout.Flex{}.Layout(gtx,
					layout.Rigid(func() {
						layout.Inset{Top: unit.Dp(0), Bottom: unit.Dp(0), Left: unit.Dp(0), Right: unit.Dp(0)}.Layout(gtx, func() {
							if Icon != nil {
								Icon.Color = theme.HexARGB(th.Color.Dark)
								Icon.Layout(gtx, unit.Px(float32(32)))
							}
							gtx.Dimensions = layout.Dimensions{
								Size: image.Point{X: 32, Y: 32},
							}
						})
					}),
					layout.Rigid(func() {
						txt := th.H6(name)
						txt.Font.Typeface = th.Font.Primary
						txt.Color = theme.HexARGB(th.Color.Primary)
						txt.Layout(gtx)
					}),
				)
			}),
			layout.Rigid(func() {
				value := th.H5(value)
				value.Font.Typeface = th.Font.Primary
				value.Color = theme.HexARGB(th.Color.Dark)
				value.Alignment = text.End
				value.Layout(gtx)
			}),
		)
	}
}