package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MizushimaToshihiko/quaketest/xmls"
)

var osExit = os.Exit

func main() {
	rand.Seed(time.Now().UnixNano())

	ch := runServer(":8080", xmls.WeatherHandler)
	// Accept and print the error from the handler.
	log.Println(<-ch)
}

func runServer(addr string, handler func(w http.ResponseWriter, r *http.Request)) chan error {
	// Listen to port 8080, and set handler to 'OmikujiHandler'.
	listener, ch := server(addr, handler)
	fmt.Println("Quake info test Server started at", listener.Addr())

	// 'ctrl+c' signal interrupt
	ctx := context.Background()
	_, cancel := interruptListen(ctx, listener)
	defer cancel()

	return ch
}

// server function creates a net.Listener that listens from 'addr',
// and invoke 'router' function by go routine,
// and reteuns the net.Listener that server created, and error channel from the server handler.
func server(addr string,
	handler func(w http.ResponseWriter, r *http.Request)) (net.Listener, chan error) {
	ch := make(chan error)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		mux := router(handler)
		ch <- http.Serve(listener, mux)
	}()

	return listener, ch
}

// router function returns the pointer of http.ServerMux that has 'handler'.
func router(handler func(w http.ResponseWriter, r *http.Request)) *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	return mux
}

// listen accepts 'ctrl+c' signal, and stop the 'Omikuji' server,
// and returns context.Context and function for clean.
func interruptListen(ctx context.Context, listener net.Listener) (context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig
		if _, err := fmt.Println("\n^Csignal : interrupt."); err != nil {
			cancel()
			log.Fatalf("listen: fmt.Println error: %s", err)
		}
		if err := listener.Close(); err != nil {
			cancel()
			log.Fatalf("listen: listener.Close error: %s\n", err)
		}
		cancel()
		osExit(0)
	}()

	return ctx, cancel
}
