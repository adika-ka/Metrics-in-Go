package handlers

import (
	"encoding/json"
	"net/http"
	"task4.2.3/internal/models"
)

// @Summary Регистрация пользователя
// @Description Регистрирует нового пользователя
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Данные пользователя"
// @Success 200 {object} string "User registered successfully"
// @Failure 400 {string} string "Bad request"
// @Router /api/register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request, store *UserStore) {
	var registerReq models.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&registerReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if store.UserExists(registerReq.Login) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "User already exists"})
		return
	}

	store.AddUser(registerReq.Login, registerReq.Password)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}
