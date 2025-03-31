package handlers

import (
	"encoding/json"
	"net/http"

	"task4.2.3/internal/middleware"
	"task4.2.3/internal/models"

	"golang.org/x/crypto/bcrypt"
)

// @Summary Авторизация пользователя
// @Description Проверяет логин и пароль, возвращает JWT-токен
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Данные пользователя"
// @Success 200 {object} map[string]string "Токен"
// @Failure 400 {string} string "Bad request"
// @Failure 401 {string} string "Unauthorized"
// @Router /api/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request, store *UserStore) {
	var loginReq models.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	hash, exist := store.users[loginReq.Login]
	if !exist {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"error": "User does not exist"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(loginReq.Password))
	if err != nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"error": "Password is incorrect"})
		return
	}

	_, tokenString, _ := middleware.TokenAuth.Encode(map[string]interface{}{
		"login": loginReq.Login,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
