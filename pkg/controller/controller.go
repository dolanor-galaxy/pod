package controller

import (
	"context"
	"math/rand"
	"sync/atomic"
	"time"

	chain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

type Blocks []*mining.BlockTemplate

func Run(cx *conte.Xt) (cancel context.CancelFunc) {
	log.WARN("starting controller")
	var ctx context.Context
	var lastBlock atomic.Value // Blocks
	ctx, cancel = context.WithCancel(context.Background())
	blockChan := make(chan Blocks)
	go func() {
		for {
			// work loop
			select {
			case lb := <-blockChan:
				lastBlock.Store(lb)
				// send out block broadcast
				log.DEBUG("sending out block broadcast")
				// log.DEBUG(spew.Sdump(lb))

			case <-ctx.Done():
				// cancel has been called
				return
			// default:
			}
		}
	}()
	connCount := cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
	current := cx.RPCServer.Cfg.SyncMgr.IsCurrent()
	// if out of sync or disconnected,
	// once a second send out empty blocks
	for (connCount < 1 && !*cx.Config.Solo) || !current {
		time.Sleep(time.Second)
		connCount = cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
		current = cx.RPCServer.Cfg.SyncMgr.IsCurrent()
		log.WARN("waiting for sync/peers", connCount, current)
		select {
		case <-ctx.Done():
			log.WARN("cancelled before initial connection/sync")
			return
		default:
		}
	}
	// load initial blocks and send them out
	lastBlock.Store(Blocks{})
	var blocks Blocks
	// generate Blocks
	for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.
		BestSnapshot().Height+1)].Algos {
		// Choose a payment address at random.
		rand.Seed(time.Now().UnixNano())
		payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.
			StateCfg.ActiveMiningAddrs))]
		template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0,
			payToAddr, algo)
		if err != nil {
			log.ERROR("failed to create new block template:", err)
			continue
		}
		blocks = append(blocks, template)
	}
	lastBlock.Store(blocks)
	log.WARN("sending initial block broadcast")
	blockChan <- lastBlock.Load().(Blocks)
	// create subscriber for new block event
	cx.RPCServer.Cfg.Chain.Subscribe(func(n *chain.
	Notification) {
		switch n.Type {
		case chain.NTBlockConnected:
			log.WARN("new block found")
			lastBlock.Store(Blocks{})
			var blocks Blocks
			// generate Blocks
			for algo := range fork.List[fork.GetCurrent(cx.RPCServer.Cfg.Chain.
				BestSnapshot().Height+1)].Algos {
				// Choose a payment address at random.
				rand.Seed(time.Now().UnixNano())
				payToAddr := cx.StateCfg.ActiveMiningAddrs[rand.Intn(len(cx.
					StateCfg.ActiveMiningAddrs))]
				template, err := cx.RPCServer.Cfg.Generator.NewBlockTemplate(0,
					payToAddr, algo)
				if err != nil {
					log.ERROR("failed to create new block template:", err)
					continue
				}
				blocks = append(blocks, template)
			}
			lastBlock.Store(blocks)
			blockChan <- lastBlock.Load().(Blocks)
		}
	})
	// goroutine loop checking for connection and sync status
	go func() {
		for {
			time.Sleep(time.Second)
			connCount := cx.RPCServer.Cfg.ConnMgr.ConnectedCount()
			current := cx.RPCServer.Cfg.SyncMgr.IsCurrent()
			// if out of sync or disconnected,
			// once a second send out empty blocks
			if connCount < 1 || !current {
				lastBlock.Store(Blocks{})
				blockChan <- lastBlock.Load().(Blocks)
			}
			select {
			case <-ctx.Done():
				break
			}
		}
	}()
	return
}
