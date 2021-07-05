package main

import (
	"context"
    "crypto/rand"
    "fmt"
    "io"
    mathrand "math/rand"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/libp2p/go-libp2p"
    "github.com/libp2p/go-libp2p-connmgr"
    "github.com/libp2p/go-libp2p-core/crypto"
    "github.com/libp2p/go-libp2p-discovery"
    "github.com/libp2p/go-libp2p-core/host"
    "github.com/libp2p/go-libp2p-core/network"
    "github.com/libp2p/go-libp2p-core/peer"
    "github.com/libp2p/go-libp2p-core/routing"
    dht "github.com/libp2p/go-libp2p-kad-dht"
    
    "github.com/multiformats/go-multiaddr"
    log "github.com/sirupsen/logrus"
)

func initAgent(config *Config) {
    log.Debug("initializing agent")

    ctx, cancel := context.WithCancel(context.Background())

    h, err := createHost(ctx, config.Seed, config.Port)
	if err != nil {
		log.Fatal(err)
	}

    log.WithFields(log.Fields {
        "id": h.ID().Pretty,
	}).Debug("Agent initialized")
    
    
	for _, addr := range h.Addrs() {
        log.WithFields(log.Fields {
        "address": fmt.Sprintf("%s/p2p/%s", addr, h.ID().Pretty()),
        }).Info("Listening on")
	}
	
	var options []dht.Option
	if len(config.Bootstrap) == 0 {
        options = append(options, dht.Mode(dht.ModeServer))
    }
    
    idht, err := dht.New(ctx, h, options...)
    if err != nil {
		log.Fatal(err)
	}
	
    err = idht.Bootstrap(ctx)
    if err != nil {
		log.Fatal(err)
	}

	if len(config.Bootstrap) > 0 {
        bootstrapAddr, err := multiaddr.NewMultiaddr(config.Bootstrap)
        if err != nil {
            log.Fatal(err)
        }
        
        peer, err := peer.AddrInfoFromP2pAddr(bootstrapAddr)
        if err != nil {
            log.Fatal(err)
        }
        
        err = h.Connect(ctx, *peer)
        if err != nil {
            log.Fatal(err)
        } else {
            log.WithFields(log.Fields {
                "node": peer,
            }).Info("Connected to bootstrap node")
        }
    }
    
    go discover(ctx, h, idht, "cypner-network")
    
    interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	<-interrupt

	log.Info("Terminating agent")

	cancel()

	err = h.Close()
    if err != nil {
		log.Fatal(err)
	}
	
	os.Exit(0)
}

func createHost(ctx context.Context, seed int64, port int) (host.Host, error) {
	var r io.Reader
	if seed == 0 {
		r = rand.Reader
	} else {
		r = mathrand.New(mathrand.NewSource(seed))
	}

	priv, _, err := crypto.GenerateEd25519Key(r)
	if err != nil {
		return nil, err
	}

	ip4addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/0.0.0.0/tcp/%d", port))
    ip6addr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip6/::/tcp/%d", port))

	return libp2p.New(ctx,
		libp2p.ListenAddrs(ip4addr, ip6addr),
		libp2p.Identity(priv),
        libp2p.ConnectionManager(connmgr.NewConnManager(
            10,
            50,
            time.Minute,
        )),
        libp2p.Routing(func(h host.Host) (routing.PeerRouting, error) {
            idht, err := dht.New(ctx, h)
            return idht, err
        }),
	)
}

func discover(ctx context.Context, h host.Host, dht *dht.IpfsDHT, networkId string) {
	var routingDiscovery = discovery.NewRoutingDiscovery(dht)
   
	discovery.Advertise(ctx, routingDiscovery, networkId, discovery.TTL(time.Minute * 1))

	for {
        log.Debug("Discovering peers")
        peers, err := discovery.FindPeers(ctx, routingDiscovery, networkId, discovery.Limit(1000))
        if err != nil {
            log.Fatal(err)
        }
    
        for _, peer := range peers {
            if peer.ID == h.ID() {
                continue
            }
            
            log.WithFields(log.Fields {
                "id": peer.ID.Pretty(),
            }).Debug("Found peer")
            
            if h.Network().Connectedness(peer.ID) != network.Connected {
                _, err = h.Network().DialPeer(ctx, peer.ID)
                if err == nil {
                    log.Info("Connected to peer", peer.ID.Pretty())
                }
            }
        }
		
		time.Sleep(time.Second * 10)
	}
}
