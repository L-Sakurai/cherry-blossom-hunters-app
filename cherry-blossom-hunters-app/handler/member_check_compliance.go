package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	"cherry-blossom-hunters-app/logger"
	"cherry-blossom-hunters-app/service"
	"cherry-blossom-hunters-app/appConfig"

)

type MemberSeviceHandler struct {
	memberService *service.MemberService
}

func NewMemberServiceHandler(appConfig *appConfig.Config) *MemberSeviceHandler {
	return &MemberSeviceHandler{
		memberService: service.NewMemberService(appConfig),
	}
}

func (h *MemberSeviceHandler) CheckCompliance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CheckCompliance called")
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

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Logging("Failed to encode JSON: "+err.Error(), logger.Error)
		http.Error(w, `{"error": "Encoding Failed"}`, http.StatusInternalServerError)
	}
}