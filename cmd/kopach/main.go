package kopach

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"time"

	"github.com/urfave/cli"
	"go.uber.org/atomic"

	"github.com/p9c/pod/cmd/kopach/client"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/kopachctrl"
	"github.com/p9c/pod/pkg/kopachctrl/job"
	"github.com/p9c/pod/pkg/kopachctrl/pause"
	"github.com/p9c/pod/pkg/kopachctrl/sol"
	"github.com/p9c/pod/pkg/stdconn/worker"
	"github.com/p9c/pod/pkg/transport"
	"github.com/p9c/pod/pkg/util/interrupt"
)

type HashCount struct {
	uint64
	Time time.Time
}

type Worker struct {
	active        atomic.Bool
	conn          *transport.Channel
	ctx           context.Context
	quit          chan struct{}
	cx            *conte.Xt
	sendAddresses []*net.UDPAddr
	workers       []*client.Client
	FirstSender   atomic.String
	lastSent      atomic.Int64
	Status        atomic.String
	HashTick      chan HashCount
	LastHash      *chainhash.Hash
}

func KopachHandle(cx *conte.Xt) func(c *cli.Context) error {
	return func(c *cli.Context) (err error) {
		L.Debug("miner controller starting")
		ctx, cancel := context.WithCancel(context.Background())
		w := &Worker{
			ctx:           ctx,
			cx:            cx,
			quit:          cx.KillAll,
			sendAddresses: []*net.UDPAddr{},
		}
		w.lastSent.Store(time.Now().UnixNano())
		w.active.Store(false)
		L.Debug("opening broadcast channel listener")
		w.conn, err = transport.
			NewBroadcastChannel("kopachmain", w, *cx.Config.MinerPass,
				transport.DefaultPort, kopachctrl.MaxDatagramSize, handlers,
				cx.KillAll)
		if err != nil {
			L.Error(err)
			cancel()
			return
		}
		var wks []*worker.Worker
		// start up the workers
		L.Debug("starting up kopach workers")
		for i := 0; i < *cx.Config.GenThreads; i++ {
			L.Debug("starting worker", i)
			cmd := worker.Spawn(os.Args[0], "worker",
				cx.ActiveNet.Name, *cx.Config.LogLevel)
			wks = append(wks, cmd)
			w.workers = append(w.workers, client.New(cmd.StdConn))
		}
		interrupt.AddHandler(func() {
			w.active.Store(false)
			L.Debug("KopachHandle interrupt")
			for i := range w.workers {
				if err = wks[i].Process.Signal(syscall.SIGKILL); !L.Check(err) {
				}
				//if err := wks[i].Stop(); L.Check(err) {
				//}
				//if err := wks[i].Kill(); L.Check(err) {
				//}
				L.Debug("stopped worker", i)
			}
		})
		for i := range w.workers {
			L.Debug("sending pass to worker", i)
			err := w.workers[i].SendPass(*cx.Config.MinerPass)
			if err != nil {
				L.Error(err)
			}
		}
		w.active.Store(true)
		// controller watcher thread
		go func() {
			L.Debug("starting controller watcher")
			ticker := time.NewTicker(time.Second)
		out:
			for {
				select {
				case <-ticker.C:
					// if the last message sent was 3 seconds ago the server is
					// almost certainly disconnected or crashed so clear
					// FirstSender
					since := time.Now().Sub(time.Unix(0, w.lastSent.Load()))
					wasSending := since > time.Second*3 && w.FirstSender.Load() != ""
					if wasSending {
						L.Debug("previous current controller has stopped"+
							" broadcasting", since, w.FirstSender.Load())
						// when this string is clear other broadcasts will be listened to
						w.FirstSender.Store("")
						// pause the workers
						for i := range w.workers {
							L.Trace("sending pause to worker", i)
							err := w.workers[i].Pause()
							if err != nil {
								L.Error(err)
							}
						}
					}
				case <-cx.KillAll:
					break out
				}
			}
		}()
		L.Debug("listening on", kopachctrl.UDP4MulticastAddress)
		<-cx.KillAll
		L.Info("kopach shutting down")
		return
	}
}

// these are the handlers for specific message types.
var handlers = transport.Handlers{
	string(job.Magic): func(ctx interface{}, src net.Addr, dst string,
		b []byte) (err error) {
		w := ctx.(*Worker)
		if !w.active.Load() {
			L.Debug("not active")
			return
		}
		j := job.LoadContainer(b)
		ips := j.GetIPs()
		cP := j.GetControllerListenerPort()
		addr := net.JoinHostPort(ips[0].String(), fmt.Sprint(cP))
		firstSender := w.FirstSender.Load()
		otherSent := firstSender != addr && firstSender != ""
		if otherSent {
			// L.Debug("ignoring other controller job")
			// ignore other controllers while one is active and received first
			return
		}
		if firstSender == "" {
			L.Warn("new sender", addr)
		}
		w.FirstSender.Store(addr)
		w.lastSent.Store(time.Now().UnixNano())
		for i := range w.workers {
			err := w.workers[i].NewJob(&j)
			if err != nil {
				L.Error(err)
			}
		}
		return
	},
	string(pause.PauseMagic): func(ctx interface{}, src net.Addr, dst string, b []byte) (err error) {
		w := ctx.(*Worker)
		p := pause.LoadPauseContainer(b)
		fs := w.FirstSender.Load()
		ni := p.GetIPs()[0].String()
		np := p.GetControllerListenerPort()
		ns := net.JoinHostPort(ni, fmt.Sprint(np))
		if fs == ns {
			for i := range w.workers {
				L.Debug("sending pause to worker", i, fs, ns)
				err := w.workers[i].Pause()
				if err != nil {
					L.Error(err)
				}
			}
		}
		return
	},
	string(sol.SolutionMagic): func(ctx interface{}, src net.Addr, dst string,
		b []byte) (err error) {
		w := ctx.(*Worker)
		// port := strings.Split(w.FirstSender.Load(), ":")[1]
		// j := sol.LoadSolContainer(b)
		// senderPort := j.GetSenderPort()
		// if fmt.Sprint(senderPort) == port {
		// 	L.Warn("we found a solution")
		// }
		w.FirstSender.Store("")
		return
	},
}
