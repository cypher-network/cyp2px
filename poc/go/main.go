package main

import (
    "flag"
    
    log "github.com/sirupsen/logrus"
)

func main() {
    config := Config{}
    
	flag.String(config.Bootstrap, "bootstrap", "Bootstrap node")
    flag.StringVar(&config.LogLevel, "loglevel", "info", "Logging level (debug, info, warn, error)")
    flag.StringVar(&config.Mode, "mode", "cli" ,"Run cyp2px as agent, cli client or text UI client ('agent', 'cli', 'tui')")
	flag.IntVar(&config.Port, "port", 0, "Port to listen on (0 for auto-assigned port)")
    flag.IntVar(&config.RpcPort, "rpc-port", 9123, "RPC port")
	flag.Int64Var(&config.Seed, "seed", 0, "PeerID seed (0 for random seed)")
	flag.Parse()
    
   	switch config.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.Fatal("Invalid loglevel: ", config.LogLevel)
	}

	log.WithFields(log.Fields {
        "bootstrap": config.Bootstrap,
        "loglevel": config.LogLevel,
        "mode": config.Mode,
		"port": config.Port,
        "rpc-port": config.RpcPort,
		"seed": config.Seed,
	}).Debug("Command-line parameters")

    switch config.Mode {
    case "agent":
        initAgent(&config)
    default:
        log.Fatal("Invalid mode: ", config.Mode)
    }
}
