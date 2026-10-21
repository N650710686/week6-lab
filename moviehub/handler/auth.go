package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"unicode/utf8"

	"moviehub/auth"
	"moviehub/mw"
	"moviehub/store"
)

type registerReq struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"displayName"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "อ่านกล่อง JSON ไม่ออก")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.DisplayName = strings.TrimSpace(req.DisplayName)

	if req.Email == "" || !strings.Contains(req.Email, "@") {
		writeError(w, http.StatusBadRequest, "BAD_EMAIL", "กรุณากรอกอีเมลให้ถูกต้อง")
		return
	}
	if utf8.RuneCountInString(req.Password) < 8 {
		writeError(w, http.StatusBadRequest, "SHORT_PASSWORD", "รหัสผ่านต้องยาวอย่างน้อย 8 ตัวอักษร")
		return
	}
	if req.DisplayName == "" || utf8.RuneCountInString(req.DisplayName) > 50 {
		writeError(w, http.StatusBadRequest, "BAD_NAME", "กรุณากรอกชื่อที่แสดง (ไม่เกิน 50 ตัวอักษร)")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		internal(w, err)
		return
	}

	m, err := h.Store.CreateMember(r.Context(), req.Email, hash, req.DisplayName)
	if errors.Is(err, store.ErrEmailTaken) {
		writeError(w, http.StatusConflict, "EMAIL_TAKEN", "อีเมลนี้ถูกใช้งานแล้ว")
		return
	}
	if err != nil {
		internal(w, err)
		return
	}

	tok, err := h.Tickets.Issue(m.ID, m.Role)
	if err != nil {
		internal(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, map[string]any{
		"token":  tok,
		"member": m,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "BAD_JSON", "อ่านกล่อง JSON ไม่ออก")
		return
	}
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))

	m, err := h.Store.MemberByEmail(r.Context(), req.Email)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusUnauthorized, "BAD_CREDENTIALS", "อีเมลหรือรหัสผ่านไม่ถูกต้อง")
		return
	}
	if err != nil {
		internal(w, err)
		return
	}

	if !auth.CheckPassword(req.Password, m.PasswordHash) {
		writeError(w, http.StatusUnauthorized, "BAD_CREDENTIALS", "อีเมลหรือรหัสผ่านไม่ถูกต้อง")
		return
	}

	tok, err := h.Tickets.Issue(m.ID, m.Role)
	if err != nil {
		internal(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"token":  tok,
		"member": m,
	})
}

// Me (ด่านที่ 4) ตอบข้อมูลของคนที่ถือตั๋ว เลขสมาชิกมาจากตั๋วที่ RequireLogin ตรวจแล้วเท่านั้น
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	id, ok := mw.MemberID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "NO_TOKEN", "กรุณาเข้าสู่ระบบก่อน")
		return
	}
	m, err := h.Store.MemberByID(r.Context(), id)
	if errors.Is(err, store.ErrNotFound) {
		writeError(w, http.StatusUnauthorized, "MEMBER_GONE", "ไม่พบบัญชีนี้แล้ว กรุณาเข้าสู่ระบบใหม่")
		return
	}
	if err != nil {
		internal(w, err)
		return
	}
	writeJSON(w, http.StatusOK, m)
}
