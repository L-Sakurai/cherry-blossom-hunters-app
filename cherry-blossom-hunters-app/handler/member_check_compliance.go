package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"cherry-blossom-hunters-app/service"
	"cherry-blossom-hunters-app/appConfig"
	"cherry-blossom-hunters-app/logger"
	"cherry-blossom-hunters-app/notify"
)

type MemberSeviceHandler struct {
	memberService *service.MemberService
	notifier *notify.Notifier
}

func NewMemberServiceHandler(cfg *appConfig.Config) *MemberSeviceHandler {
    return &MemberSeviceHandler{
        memberService: service.NewMemberService(cfg),
        notifier:      notify.NewDiscordNotifier(cfg.Notify.Webhooks["member"]),
    }
}

func (h *MemberSeviceHandler) CheckCompliance(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	
	res, err := h.memberService.CheckMemberComplianceWithContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "メンバーチェックに失敗しました",
			"details": err.Error(),
		})
		return
	}
	
	response := map[string]interface{}{
		"message": res.Message,
		"status": res.Diff,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	payload := &notify.CustomPayload{
		Title:   "🔔 メンバー差分検知機能",
		Message: res.Message,
		Fields: map[string]interface{}{
			"いいねを確認できていないユーザ詳細": res.Diff,
		},
	}

	// 修正: notifier.Sendメソッドを正しく呼び出し
	if err := h.notifier.Send(payload); err != nil {
		logger.Logging("error", "Webhook送信失敗: %v", err)
	} else {
		logger.Logging("info", "通知に成功しました")
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logging("error", "%v", err)
		http.Error(w, `{"error": "Encoding Failed"}`, http.StatusInternalServerError)
	}
}