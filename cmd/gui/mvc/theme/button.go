// SPDX-License-Identifier: Unlicense OR MIT

package theme

import (
	"github.com/p9c/pod/pkg/gui/io/pointer"
	"image"
	"image/color"

	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/layout"
	"github.com/p9c/pod/pkg/gui/op"
	"github.com/p9c/pod/pkg/gui/op/clip"
	"github.com/p9c/pod/pkg/gui/op/paint"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/gui/unit"
)

type Button struct {
	Text string
	// Color is the text color.
	Color        color.RGBA
	Font         text.Font
	TextSize     unit.Value
	Background   color.RGBA
	CornerRadius unit.Value
	shaper       text.Shaper
}

type IconButton struct {
	Background color.RGBA
	Color      color.RGBA
	Icon       *DuoUIicon
	Size       unit.Value
	Padding    unit.Value
}

func (t *DuoUItheme) Button(txt string) Button {
	return Button{
		Font: text.Font{
			Typeface: t.Font.Primary,
		},
		Text:       txt,
		Color:      rgb(0xffffff),
		Background: HexARGB(t.Color.Primary),
		TextSize:   t.TextSize.Scale(14.0 / 16.0),
		shaper:     t.Shaper,
	}
}

func (t *DuoUItheme) IconButton(icon *DuoUIicon) IconButton {
	return IconButton{
		Background: HexARGB(t.Color.Primary),
		Color:      HexARGB(t.Color.InvText),
		Icon:       icon,
		Size:       unit.Dp(56),
		Padding:    unit.Dp(16),
	}
}

func (b Button) Layout(gtx *layout.Context, button *controller.Button) {
	col := b.Color
	bgcol := b.Background
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(4)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, bgcol)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin
			layout.Center.Layout(gtx, func() {
				layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func() {
					paint.ColorOp{Color: col}.Add(gtx.Ops)
					controller.Label{}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Text)
				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func (b IconButton) Layout(gtx *layout.Context, button *controller.Button) {
	layout.Stack{}.Layout(gtx,
		layout.Expanded(func() {
			size := float32(gtx.Constraints.Width.Min)
			rr := float32(size) * .5
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{X: size, Y: size}},
				NE:   rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, b.Background)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			layout.UniformInset(b.Padding).Layout(gtx, func() {
				size := gtx.Px(b.Size) - 2*gtx.Px(b.Padding)
				if b.Icon != nil {
					b.Icon.Color = b.Color
					b.Icon.Layout(gtx, unit.Px(float32(size)))
				}
				gtx.Dimensions = layout.Dimensions{
					Size: image.Point{X: size, Y: size},
				}
			})
			pointer.Ellipse(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func toPointF(p image.Point) f32.Point {
	return f32.Point{X: float32(p.X), Y: float32(p.Y)}
}

func toRectF(r image.Rectangle) f32.Rectangle {
	return f32.Rectangle{
		Min: toPointF(r.Min),
		Max: toPointF(r.Max),
	}
}

func drawInk(gtx *layout.Context, c controller.Click) {
	d := gtx.Now().Sub(c.Time)
	t := float32(d.Seconds())
	const duration = 0.5
	if t > duration {
		return
	}
	t = t / duration
	var stack op.StackOp
	stack.Push(gtx.Ops)
	size := float32(gtx.Px(unit.Dp(700))) * t
	rr := size * .5
	col := byte(0xaa * (1 - t*t))
	ink := paint.ColorOp{Color: color.RGBA{A: col, R: col, G: col, B: col}}
	ink.Add(gtx.Ops)
	op.TransformOp{}.Offset(c.Position).Offset(f32.Point{
		X: -rr,
		Y: -rr,
	}).Add(gtx.Ops)
	clip.Rect{
		Rect: f32.Rectangle{Max: f32.Point{
			X: float32(size),
			Y: float32(size),
		}},
		NE: rr, NW: rr, SE: rr, SW: rr,
	}.Op(gtx.Ops).Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rectangle{Max: f32.Point{X: float32(size), Y: float32(size)}}}.Add(gtx.Ops)
	stack.Pop()
	op.InvalidateOp{}.Add(gtx.Ops)
}

var (
	buttonInsideLayoutList = &layout.List{
		Axis: layout.Vertical,
	}
)

type DuoUIbutton struct {
	Text string
	// Color is the text color.
	TxColor           color.RGBA
	Font              text.Font
	TextSize          unit.Value
	Width             int
	Height            int
	BgColor           color.RGBA
	CornerRadius      unit.Value
	Icon              *DuoUIicon
	IconSize          int
	IconColor         color.RGBA
	PaddingVertical   unit.Value
	PaddingHorizontal unit.Value
	shaper            text.Shaper
	hover             bool
}

func (t *DuoUItheme) DuoUIbutton(txtFont text.Typeface, txt, txtColor, bgColor, icon, iconColor string, textSize, iconSize, width, height, paddingVertical, paddingHorizontal int) DuoUIbutton {
	return DuoUIbutton{
		Text: txt,
		Font: text.Font{
			Typeface: txtFont,
		},
		TextSize:          unit.Dp(float32(textSize)),
		Width:             width,
		Height:            height,
		TxColor:           HexARGB(txtColor),
		BgColor:           HexARGB(bgColor),
		Icon:              t.Icons[icon],
		IconSize:          iconSize,
		IconColor:         HexARGB(iconColor),
		PaddingVertical:   unit.Dp(float32(paddingVertical)),
		PaddingHorizontal: unit.Dp(float32(paddingHorizontal)),
		shaper:            t.Shaper,
	}
}

func (b DuoUIbutton) Layout(gtx *layout.Context, button *controller.Button) {
	hmin := gtx.Constraints.Width.Min
	vmin := gtx.Constraints.Height.Min
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(gtx.Constraints.Width.Min),
					Y: float32(gtx.Constraints.Height.Min),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, b.BgColor)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = hmin
			gtx.Constraints.Height.Min = vmin
			layout.Center.Layout(gtx, func() {
				layout.Inset{Top: unit.Dp(10), Bottom: unit.Dp(10), Left: unit.Dp(12), Right: unit.Dp(12)}.Layout(gtx, func() {

					paint.ColorOp{Color: b.TxColor}.Add(gtx.Ops)
					controller.Label{
						Alignment: text.Middle,
					}.Layout(gtx, b.shaper, b.Font, b.TextSize, b.Text)
				})
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func (b DuoUIbutton) IconLayout(gtx *layout.Context, button *controller.Button) {
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(b.Width),
					Y: float32(b.Height),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, b.BgColor)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = b.Width
			gtx.Constraints.Height.Min = b.Height
			layout.Center.Layout(gtx, func() {
				layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
					b.Icon.Color = b.IconColor
					b.Icon.Layout(gtx, unit.Dp(float32(b.IconSize)))
				})
				gtx.Dimensions = layout.Dimensions{
					Size: image.Point{X: b.IconSize, Y: b.IconSize},
				}
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}

func (b DuoUIbutton) MenuLayout(gtx *layout.Context, button *controller.Button) {
	layout.Stack{Alignment: layout.Center}.Layout(gtx,
		layout.Expanded(func() {
			rr := float32(gtx.Px(unit.Dp(0)))
			clip.Rect{
				Rect: f32.Rectangle{Max: f32.Point{
					X: float32(b.Width),
					Y: float32(b.Height),
				}},
				NE: rr, NW: rr, SE: rr, SW: rr,
			}.Op(gtx.Ops).Add(gtx.Ops)
			fill(gtx, b.BgColor)
			for _, c := range button.History() {
				drawInk(gtx, c)
			}
		}),
		layout.Stacked(func() {
			gtx.Constraints.Width.Min = b.Width
			gtx.Constraints.Height.Min = b.Height
			layout.Center.Layout(gtx, func() {
				layout.Flex{Axis: layout.Vertical, Alignment: layout.Middle}.Layout(gtx,
					layout.Rigid(func() {
						layout.Center.Layout(gtx, func() {
							layout.UniformInset(unit.Dp(0)).Layout(gtx, func() {
								b.Icon.Color = b.IconColor
								b.Icon.Layout(gtx, unit.Dp(float32(b.IconSize)))
							})
							gtx.Dimensions = layout.Dimensions{
								Size: image.Point{X: b.IconSize, Y: b.IconSize},
							}
						})
					}),
					layout.Rigid(func() {
						layout.Center.Layout(gtx, func() {
							paint.ColorOp{Color: b.TxColor}.Add(gtx.Ops)
							controller.Label{
								Alignment: text.Middle,
							}.Layout(gtx, b.shaper, b.Font, unit.Dp(12), b.Text)
						})
					}))
			})
			pointer.Rect(image.Rectangle{Max: gtx.Dimensions.Size}).Add(gtx.Ops)
			button.Layout(gtx)
		}),
	)
}
