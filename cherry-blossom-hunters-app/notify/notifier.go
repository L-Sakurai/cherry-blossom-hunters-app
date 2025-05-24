package notify

import (
	"context"
	"time"

	"cherry-blossom-hunters-app/logger"
)

// NotifyUserShutdown は既存の関数（変更なし）
func NotifyUserShutdown() {
	// 既存の通知ロジックをここに実装
	logger.Logging("Sending shutdown notifications to users...")
	
	// 例: 実際の通知処理（メール、Slack、Webhookなど）
	// この部分は元のコードに依存します
	time.Sleep(200 * time.Millisecond) // 通知処理のシミュレーション
	
	logger.Logging("User shutdown notifications sent successfully")
}

// NotifyUserShutdownWithContext はコンテキスト対応版
func NotifyUserShutdownWithContext(ctx context.Context) error {
	logger.Logging("Sending shutdown notifications to users...")
	
	// コンテキストのキャンセレーションをチェックしながら通知処理
	select {
	case <-ctx.Done():
		logger.Logging("Notification canceled due to context timeout", logger.Warn)
		return ctx.Err()
	default:
		// 実際の通知処理
		// 例: 外部APIコール、データベース更新など
		time.Sleep(200 * time.Millisecond) // 通知処理のシミュレーション
	}
	
	logger.Logging("User shutdown notifications sent successfully")
	return nil
}

// NotifyUserShutdownWithTimeout はタイムアウト付きの通知
func NotifyUserShutdownWithTimeout(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	return NotifyUserShutdownWithContext(ctx)
}