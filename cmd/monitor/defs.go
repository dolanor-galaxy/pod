package monitor

import (
	"encoding/json"
	"gioui.org/app"
	"gioui.org/layout"
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/gel"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/ring"
	"github.com/p9c/pod/pkg/stdconn/worker"
	"io/ioutil"
	"path/filepath"
)

const ConfigFileName = "monitor.json"

type State struct {
	Ctx                       *conte.Xt
	Gtx                       *layout.Context
	W                         *app.Window
	Worker                    *worker.Worker
	Rc                        *rcd.RcVar
	Theme                     *gelook.DuoUItheme
	Config                    *Config
	Buttons                   map[string]*gel.Button
	FilterLevelsButtons       []gel.Button
	FilterButtons             []gel.Button
	Lists                     map[string]*layout.List
	ModesButtons              map[string]*gel.Button
	CommandEditor             gel.Editor
	WindowWidth, WindowHeight int
	Loggers                   *Node
	RunningInRepo             bool
	HasGo                     bool
	HasOtherGo                bool
	CannotRun                 bool
	RunCommandChan            chan string
	EntryBuf                  *ring.Entry
	FilterRoot                *Node
}

func NewMonitor(cx *conte.Xt, gtx *layout.Context, rc *rcd.RcVar) (s *State) {
	s = &State{
		Ctx:                 cx,
		Gtx:                 gtx,
		Rc:                  rc,
		Theme:               gelook.NewDuoUItheme(),
		ModesButtons:        make(map[string]*gel.Button),
		Config:              &Config{FilterNodes: make(map[string]*Node)},
		WindowWidth:         800,
		WindowHeight:        600,
		RunCommandChan:      make(chan string),
		EntryBuf:            ring.NewEntry(65536),
		FilterLevelsButtons: make([]gel.Button, 7),
		Buttons:             make(map[string]*gel.Button),
		Lists:               make(map[string]*layout.List),
	}
	modes := []string{
		"node", "wallet", "shell", "gui", "mon",
	}
	for i := range modes {
		s.ModesButtons[modes[i]] = new(gel.Button)
	}
	buttons := []string{
		"Close", "Restart", "Logo", "RunMenu", "StopMenu", "PauseMenu",
		"RestartMenu", "KillMenu", "RunModeFold", "SettingsFold",
		"SettingsClose", "SettingsZoom", "BuildFold",
		"BuildClose", "BuildZoom", "BuildTitleClose", "Filter",
		"FilterHeader", "FilterAll", "FilterHide", "FilterShow",
		"FilterNone", "FilterClear", "FilterSend", "RunningInRepo",
		"RunFromProfile", "UseBuiltinGo", "InstallNewGo",}
	for i := range buttons {
		s.Buttons[buttons[i]] = new(gel.Button)
	}
	lists := []string{
		"Modes", "FilterLevel", "Groups", "Filter", "Log",
		"SettingsFields",
	}
	for i := range lists {
		s.Lists[lists[i]] = new(layout.List)
	}
	s.Lists = map[string]*layout.List{
		"Modes": {
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		},
		"Groups": {
			Axis:      layout.Horizontal,
			Alignment: layout.Start,
		},
		"SettingsFields": {
			Axis: layout.Vertical,
		},
		"Filter":      {},
		"FilterLevel": {},
		"Log":         {},
	}
	s.Config.RunMode = "node"
	s.Config.DarkTheme = true
	return
}

type TreeNode struct {
	Closed, Hidden bool
}

type Config struct {
	Width          int
	Height         int
	RunMode        string
	RunModeOpen    bool
	RunModeZoomed  bool
	SettingsOpen   bool
	SettingsZoomed bool
	SettingsTab    string
	BuildOpen      bool
	BuildZoomed    bool
	DarkTheme      bool
	RunInRepo      bool
	UseBuiltinGo   bool
	Running        bool
	Pausing        bool
	FilterOpen     bool
	FilterNodes    map[string]*Node
	FilterLevel    int
	ClickCommand   string
}

func (s *State) LoadConfig() (isNew bool) {
	Debug("loading config")
	var err error
	cnf := &Config{}
	filename := filepath.Join(*s.Ctx.Config.DataDir, ConfigFileName)
	if apputil.FileExists(filename) {
		var b []byte
		if b, err = ioutil.ReadFile(filename); !Check(err) {
			if err = json.Unmarshal(b, cnf); Check(err) {
				s.SaveConfig()
			}
			if s.Config.FilterNodes == nil {
				s.Config.FilterNodes = make(map[string]*Node)
			}
			for i := range cnf.FilterNodes {
				if s.Config.FilterNodes[i] == nil {
					s.Config.FilterNodes[i] = &Node{}
				}
				s.Config.FilterNodes[i].Hidden = cnf.FilterNodes[i].Hidden
				s.Config.FilterNodes[i].Closed = cnf.FilterNodes[i].Closed
			}
			s.Config.Width = cnf.Width
			s.Config.Height = cnf.Height
			s.Config.RunMode = cnf.RunMode
			s.Config.RunModeOpen = cnf.RunModeOpen
			s.Config.RunModeZoomed = cnf.RunModeZoomed
			s.Config.SettingsOpen = cnf.SettingsOpen
			s.Config.SettingsZoomed = cnf.SettingsZoomed
			s.Config.SettingsTab = cnf.SettingsTab
			s.Config.BuildOpen = cnf.BuildOpen
			s.Config.BuildZoomed = cnf.BuildZoomed
			s.Config.DarkTheme = cnf.DarkTheme
			s.Config.RunInRepo = cnf.RunInRepo
			s.Config.UseBuiltinGo = cnf.UseBuiltinGo
			s.Config.Running = cnf.Running
			s.Config.Pausing = cnf.Pausing
			s.Config.FilterOpen = cnf.FilterOpen
			s.Config.FilterLevel = cnf.FilterLevel
			s.Config.ClickCommand = cnf.ClickCommand
			s.CommandEditor.SetText(s.Config.ClickCommand)
		}
	} else {
		Warn("creating new configuration")
		s.Config.UseBuiltinGo = s.HasGo
		s.Config.RunInRepo = s.RunningInRepo
		isNew = true
		s.SaveConfig()
	}
	if s.Config.Width < 1 || s.Config.Height < 1 {
		s.Config.Width = 800
		s.Config.Height = 600
	}
	if s.Config.SettingsTab == "" {
		s.Config.SettingsTab = "config"
	}
	s.Rc.Settings.Tabs.Current = s.Config.SettingsTab
	s.SetTheme(s.Config.DarkTheme)
	return
}

func (s *State) SaveConfig() {
	Debug("saving monitor config")
	s.Config.Width = s.WindowWidth
	s.Config.Height = s.WindowHeight
	filename := filepath.Join(*s.Ctx.Config.DataDir, ConfigFileName)
	if yp, e := json.MarshalIndent(s.Config, "", "  "); !Check(e) {
		apputil.EnsureDir(filename)
		if e := ioutil.WriteFile(filename, yp, 0600); Check(e) {
			panic(e)
		}
	}
}
