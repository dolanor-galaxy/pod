package monitor

import (
	"gioui.org/layout"
)

func (s *State) BottomBar() layout.FlexChild {
	return Rigid(func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
		s.FlexV(
			s.SettingsPage(),
			s.BuildPage(),
			s.StatusBar(),
		)
	})
}

func (s *State) StatusBar() layout.FlexChild {
	return Rigid(func() {
		cs := s.Gtx.Constraints
		s.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg", "ff")
		s.FlexH(
			s.RunControls(),
			s.RunmodeButtons(),
			s.BuildButtons(),
			s.SettingsButtons(),
			s.Spacer(),
			s.Filter(),
		)
	})
}

func (s *State) RunmodeButtons() layout.FlexChild {
	return Rigid(func() {
		s.FlexH(Rigid(func() {
			if !s.Config.RunModeOpen {
				fg, bg := "ButtonText", "ButtonBg"
				if s.Config.Running {
					fg, bg = "ButtonBg", "DocText"
				}
				txt := s.Config.RunMode
				b := s.Buttons["RunModeFold"]
				s.TextButton(txt, "Secondary", 34, fg, bg, b)
				for b.Clicked(s.Gtx) {
					if !s.Config.Running {
						s.Config.RunModeOpen = true
						s.SaveConfig()
					}
				}
			} else {
				modes := []string{
					"node", "wallet", "shell", "gui", "mon",
				}
				s.Lists["Modes"].Layout(s.Gtx, len(modes), func(i int) {
					mm := modes[i]
					fg := "DocBg"
					if modes[i] == s.Config.RunMode {
						fg = "DocText"
					}
					txt := mm
					if s.WindowWidth <= 880 && s.Config.FilterOpen ||
						s.WindowWidth <= 640 && !s.Config.FilterOpen {
						txt = txt[:1]
					}
					cs := s.Gtx.Constraints
					s.Rectangle(cs.Width.Max, cs.Height.Max, "ButtonBg",
						"ff")
					s.TextButton(txt, "Secondary", 34, fg,
						"ButtonBg", s.ModesButtons[modes[i]])
					for s.ModesButtons[modes[i]].Clicked(s.Gtx) {
						Debug(mm, "clicked")
						if s.Config.RunModeOpen {
							s.Config.RunMode = modes[i]
							s.Config.RunModeOpen = false
						}
						s.SaveConfig()
					}
				})
			}
		}),
		)
	})
}

func (s *State) Filter() layout.FlexChild {
	fg, bg := "PanelText", "PanelBg"
	if s.Config.FilterOpen {
		fg, bg = "DocText", "DocBg"
	}
	return Rigid(func() {
		b := s.Buttons["Filter"]
		s.ButtonArea(func() {
			cs := s.Gtx.Constraints
			s.Rectangle(cs.Width.Max, cs.Height.Max, bg, "ff")
			s.Inset(8, func() {
				s.Icon("Filter", fg, bg, 32)
			})
		}, b)
		for b.Clicked(s.Gtx) {
			Debug("clicked filter button")
			if !s.Config.FilterOpen {
				s.Config.SettingsOpen = false
				s.Config.BuildOpen = false
			}
			s.Config.FilterOpen = !s.Config.FilterOpen
			s.SaveConfig()
		}
		//}
	})
}
