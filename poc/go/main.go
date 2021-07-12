package main

import (
	"github.com/gohornet/hornet/core/app"
	"github.com/gohornet/hornet/pkg/node"

	"tangrams.io/cyp2px/core/p2p"
	"tangrams.io/cyp2px/plugins/p2pdisc"
)

func main() {
	node.Run(
		node.WithInitPlugin(app.InitPlugin),
		node.WithCorePlugins([]*node.CorePlugin{
			p2p.CorePlugin,
		}...),
		node.WithPlugins([]*node.Plugin{
			p2pdisc.Plugin,
		}...),
	)
}