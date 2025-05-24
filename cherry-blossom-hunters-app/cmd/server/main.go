package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cherry-blossom-hunters-app/routes"
	"cherry-blossom-hunters-app/logger"
	"cherry-blossom-hunters-app/notify"
)

func main() {
	logger.SetUp()

	// グレイスフルシャットダウン用のチャネル
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// HTTPシャットダウン用のチャネル
	httpShutdown := make(chan bool, 1)

	// HTTPサーバーの設定
	server := &http.Server{
		Addr:    ":8080",
		Handler: routes.SetupRoutes(httpShutdown),
	}

	// サーバー起動用のゴルーチン
	go func() {
		logger.Logging("Server listening on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logging("Server error: "+err.Error(), logger.Error)
			os.Exit(1)
		}
	}()

	// シャットダウン待機
	select {
	case sig := <-shutdown:
		logger.Logging("Received signal: " + sig.String() + ". Initiating graceful shutdown...")
	case <-httpShutdown:
		logger.Logging("HTTP shutdown request received. Initiating graceful shutdown...")
	}

	// グレイスフルシャットダウンの実行
	gracefulShutdown(server)

}


func gracefulShutdown(server *http.Server) {
	// シャットダウンのタイムアウトを設定（30秒）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var wg sync.WaitGroup

	// 外部通知を並行して実行
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Logging("Sending shutdown notifications...")
		notify.NotifyUserShutdown()
	}()

	// HTTPサーバーのシャットダウン
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Logging("Shutting down HTTP server...")
		if err := server.Shutdown(ctx); err != nil {
			logger.Logging("Server shutdown error: "+err.Error(), logger.Error)
		} else {
			logger.Logging("HTTP server shutdown complete")
		}
	}()

	// すべてのシャットダウン処理の完了を待機
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Logging("Graceful shutdown completed successfully")
	case <-ctx.Done():
		logger.Logging("Shutdown timeout exceeded, forcing exit", logger.Warn)
	}

	// 最終クリーンアップ
	time.Sleep(100 * time.Millisecond)
	logger.Logging("Application exiting...")
}