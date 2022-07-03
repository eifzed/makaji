package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi"
)

func ListenAndServe(port string, router *chi.Mux) error {
	srv := http.Server{
		ReadTimeout:  10 * time.Second, // TODO: read it from config
		WriteTimeout: 3 * time.Second,  // TODO: read it from config
		Handler:      router,
	}
	listener, err := GetListener(port)
	if err != nil {
		return err
	}
	go func() {
		fmt.Println("listening with pid " + fmt.Sprint(os.Getpid()))
		srv.Serve(listener)
	}()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	fmt.Println("got signal", <-sigs)
	fmt.Println("shutting down..")
	srv.Shutdown(context.Background())
	return nil
}

func GetListener(port string) (net.Listener, error) {
	var listener net.Listener

	EINHORN_FDS := os.Getenv("EINHORN_FDS")
	if EINHORN_FDS != "" {
		fds, err := strconv.Atoi(EINHORN_FDS)
		if err != nil {
			return nil, err
		}
		log.Println("using socket master, listening on", EINHORN_FDS)
		f := os.NewFile(uintptr(fds), "listener")
		listener, err = net.FileListener(f)
		if err != nil {
			log.Fatalln("error create listener", err)
			return nil, err
		}
		return listener, nil
	}
	return net.Listen("tcp4", port)
}
