package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"fmt"
	// "cherry-blossom-hunters-app/logger"
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
	
	err := h.memberService.CheckMemberComplianceWithContext(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "メンバーチェックに失敗しました",
			"details": err.Error(),
		})
		return
	}
}