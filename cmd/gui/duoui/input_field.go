package duoui

import (
	"github.com/p9c/pod/cmd/gui/mvc/controller"
	"github.com/p9c/pod/pkg/gui/text"
	"github.com/p9c/pod/pkg/log"
	"reflect"
)

func (ui *DuoUI) DuoUIinputField(f *Field) func() {
	return func() {
		e := ui.ly.Theme.DuoUIeditor(f.Field.Name)
		e.Font.Typeface = ui.ly.Theme.Font.Primary
		e.Font.Style = text.Italic
		e.Layout(ui.ly.Context, (ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor))
		(ui.rc.Settings.Daemon.Widgets[f.Field.Name]).(*controller.Editor).SetText(f.Field.Value.(reflect.Value).String())
		log.INFO(f.Field.Value.(reflect.Value).String())
		//for _, e := range lineEditor.Events(ui.ly.Context) {
		//	if _, ok := e.(controller.SubmitEvent); ok {
		//		//topLabel = e.Text
		//		lineEditor.SetText(f.Field.Value.(reflect.Value).String())
		//		log.INFO(f.Field.Value.(reflect.Value).String())
		//	}
		//}
	}
}
