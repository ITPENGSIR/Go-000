/**
* @Author: 彭光豪
* @Date: 12/5/20 8:42 PM
 */
//基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够 一个退出，全部注销退出
package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())
	g.Go(func() error {
		if err := startServer(ctx, ":8080"); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		if err := startServer(ctx, ":8888"); err != nil {
			return err
		}
		return nil
	})
	g.Go(func() error {
		ch := make(chan os.Signal)
		signal.Notify(ch)
		select {
		case <-ch:
			return errors.New("signal error")
		case <-ctx.Done():
			return ctx.Err()
		}
		return nil
	})
	if err := g.Wait(); err != nil {
		fmt.Println("error group return err:", err.Error())
	}
	fmt.Println("server is exit")
}

func startServer(ctx context.Context, address string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte("hello world "))
	})
	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}
	go func() {
		<-ctx.Done()
		fmt.Println("server received close signal")
		timeout, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		server.Shutdown(timeout)
		fmt.Println("server shutdown")
		server.Shutdown(ctx)
	}()
	fmt.Println("server: listening on port :"+address)
	return server.ListenAndServe()
}



