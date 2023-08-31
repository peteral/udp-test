package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	listenPort int
	sendPort   int
	host       string
	period     time.Duration
)

func sender(ctx context.Context) {
	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", host, sendPort))
	if err != nil {
		slog.Error("Failed resolving send address", "error", err)
		return
	}

	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		slog.Error("Dial", "error", err)
	}
	defer c.Close()

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(period):
			_, err = fmt.Fprint(c, "Hello\000")
			if err != nil {
				slog.Error("failed sending", "error", err)
			}
		}
	}
}

func main() {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt, syscall.SIGTERM)

	flag.IntVar(&listenPort, "listenPort", 2601, "port number")
	flag.IntVar(&sendPort, "sendPort", 2602, "port number")
	flag.StringVar(&host, "host", "", "host or ip address")
	flag.DurationVar(&period, "period", time.Second, "telegram send period")

	flag.Parse()

	s, err := net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", listenPort))
	if err != nil {
		panic(err)
	}

	conn, err := net.ListenUDP("udp4", s)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	slog.Info("Server started", "port", listenPort)

	ctx, cancel := context.WithCancel(context.Background())
	go sender(ctx)

	buffer := make([]byte, 1024)

	for {
		n, addr, _ := conn.ReadFromUDP(buffer)
		slog.Info("Incoming message", "IP", addr.IP, "port", addr.Port, "data", string(buffer[:n-1]))

		select {
		case <-sigChannel:
			slog.Info("Graceful shutdown initiated")
			cancel()

			slog.Info("Exiting process")
			os.Exit(0)
		default:
		}
	}
}
