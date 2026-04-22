//go:build !js

package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/NimbleMarkets/go-booba-example/internal/model"
	boobaServe "github.com/NimbleMarkets/go-booba/serve"
)

func main() {
	listen := flag.String("listen", "127.0.0.1:8080", "listen address")
	flag.Parse()

	host, portStr, err := net.SplitHostPort(*listen)
	if err != nil {
		log.Fatalf("invalid listen address %q: %v", *listen, err)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("invalid port %q: %v", portStr, err)
	}

	cfg := boobaServe.DefaultConfig()
	cfg.Host = host
	cfg.Port = port

	srv := boobaServe.NewServer(cfg)

	handler := func(sess boobaServe.Session) (tea.Model, []tea.ProgramOption) {
		return model.InitialModel(), boobaServe.MakeOptions(sess)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := srv.Serve(ctx, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
