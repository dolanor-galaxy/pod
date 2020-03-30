package monitor

import (
	"gioui.org/layout"
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/logi/consume"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

func (s *State) RunControls() layout.FlexChild {
	return Rigid(func() {
		if s.CannotRun {
			return
		}
		if !s.Config.Running {
			s.IconButton("Run", "PanelBg", "PanelText", &s.RunMenuButton)
			for s.RunMenuButton.Clicked(s.Gtx) {
				Debug("clicked run button")
				if !s.Config.RunModeOpen {
					s.RunCommandChan <- "run"
				}
			}
		} else {
			ic := "Pause"
			fg, bg := "PanelBg", "PanelText"
			if s.Config.Pausing {
				ic = "Run"
				fg, bg = "PanelText", "PanelBg"
			}
			s.FlexH(Rigid(func() {
				s.IconButton("Stop", "PanelBg", "PanelText",
					&s.StopMenuButton)
				for s.StopMenuButton.Clicked(s.Gtx) {
					Debug("clicked stop button")
					s.RunCommandChan <- "stop"
				}
			}), Rigid(func() {
				s.IconButton(ic, fg, bg, &s.PauseMenuButton)
				for s.PauseMenuButton.Clicked(s.Gtx) {
					if s.Config.Pausing {
						Debug("clicked on resume button")
						s.RunCommandChan <- "resume"
					} else {
						Debug("clicked pause button")
						s.RunCommandChan <- "pause"
					}
				}
				//}), Rigid(func() {
				//	s.IconButton("Kill", "PanelBg", "PanelText",
				//		s.KillMenuButton)
				//	for s.KillMenuButton.Clicked(s.Gtx) {
				//		Debug("clicked kill button")
				//		s.RunCommandChan <- "kill"
				//	}
			}), Rigid(func() {
				s.IconButton("Restart", "PanelBg", "PanelText",
					&s.RestartMenuButton)
				for s.RestartMenuButton.Clicked(s.Gtx) {
					Debug("clicked restart button")
					s.RunCommandChan <- "restart"
				}
			}),
			)
		}
	})
}

func (s *State) Build() (exePath string, err error) {
	var c *exec.Cmd
	gt := "goterm"
	if runtime.GOOS == "windows" {
		gt = ""
	}
	exePath = filepath.Join(*s.Ctx.Config.DataDir, "pod_mon")
	c = exec.Command("go", "build", "-v",
		"-tags", gt, "-o", exePath)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err = c.Run(); !Check(err) {
	}
	return
}

func (s *State) Runner() {
	logi.L.SetLevel(*s.Ctx.Config.LogLevel, true, "pod")
	var err error
	var exePath string
	var quit chan struct{}
	for cmd := range s.RunCommandChan {
		switch cmd {
		case "run":
			Debug("run called")
			if s.HasGo && !s.Config.Running {
				if exePath, err = s.Build(); !Check(err) {
					quit = make(chan struct{})
					s.Worker = consume.Log(quit, func(ent *logi.Entry) (
						err error) {
						//Debugf("KOPACH %s %s", ent.Text, ent.Level)
						s.EntryBuf.Add(ent)
						return
					}, exePath, "-D", *s.Ctx.Config.DataDir, s.Config.RunMode)
					consume.Start(s.Worker)
					s.Config.Running = true
					s.Config.Pausing = false
					s.W.Invalidate()
					go func() {
						if err = s.Worker.Wait(); !Check(err) {
							s.Config.Running = false
							s.Config.Pausing = false
							s.W.Invalidate()
						}
					}()
				}
			}
		case "stop":
			Debug("stop called")
			if s.HasGo && s.Worker != nil && s.Config.Running {
				close(quit)
				if err = s.Worker.Interrupt(); !Check(err) {
					s.Config.Running = false
				}
			}
		case "pause":
			Debug("pause called")
			if s.HasGo && s.Worker != nil && s.Config.Running && !s.Config.
				Pausing {
				s.Config.Pausing = !s.Config.Pausing
				consume.Stop(s.Worker)
				if err = s.Worker.Pause(); Check(err) {
				}
			}
		case "resume":
			Debug("resume called")
			if s.HasGo && s.Worker != nil && s.Config.Running && s.Config.
				Pausing {
				s.Config.Pausing = !s.Config.Pausing
				if err = s.Worker.Resume(); Check(err) {
				}
				consume.Start(s.Worker)
			}
		case "kill":
			Debug("kill called")
			if s.HasGo && s.Worker != nil && s.Config.Running {
				close(quit)
				if err = s.Worker.Interrupt(); !Check(err) {
				}
			}
		case "restart":
			Debug("restart called")
			if s.HasGo && s.Worker != nil {
				go func() {
					s.RunCommandChan <- "stop"
					time.Sleep(time.Second)
					s.RunCommandChan <- "run"
				}()
			}
		}
	}
	return
}
