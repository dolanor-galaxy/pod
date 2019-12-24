package helpers

import (
	"image"
	"image/color"

	"github.com/p9c/pod/pkg/gio/f32"
	"github.com/p9c/pod/pkg/gio/layout"
	"github.com/p9c/pod/pkg/gio/op/clip"
	"github.com/p9c/pod/pkg/gio/op/paint"
	"github.com/p9c/pod/pkg/gio/unit"
)


func DuoUIdrawRectangle(gtx *layout.Context, w, h int, color color.RGBA, ne, nw, se, sw float32, inset unit.Value) {
	in := layout.UniformInset(inset)
	in.Layout(gtx, func() {
		square := f32.Rectangle{
			Max: f32.Point{
				X: float32(w),
				Y: float32(h),
			},
		}
		paint.ColorOp{Color: color}.Add(gtx.Ops)

		clip.Rect{Rect: square,
			NE: ne, NW: nw, SE: se, SW: sw}.Op(gtx.Ops).Add(gtx.Ops) // HLdraw
		paint.PaintOp{Rect: square}.Add(gtx.Ops)
		gtx.Dimensions = layout.Dimensions{Size: image.Point{X: w, Y: h}}
	})
}