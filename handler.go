package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func writeMessage(w http.ResponseWriter, statusCode int, message string) error {
	w.Header().Add("content-type", "application/json")
	w.WriteHeader(statusCode)

	var j struct {
		Message string `json:"message"`
	}

	j.Message = message

	err := json.NewEncoder(w).Encode(j)
	if err != nil {
		return err
	}

	return nil
}

func writeError(w http.ResponseWriter, statusCode int, err error) error {
	return writeMessage(w, statusCode, err.Error())
}

func RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	type UserRegister struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var ur UserRegister
	err := json.NewDecoder(r.Body).Decode(&ur)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	defer r.Body.Close()
	if ur.Name == "" || ur.Email == "" || ur.Password == "" {
		writeError(w, http.StatusBadRequest, errors.New("bad request"))
		return
	}

	err = RegisterUserService(r.Context(), ur.Name, ur.Email, ur.Password)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	writeMessage(w, http.StatusCreated, "Registration is success")
}

func LoginUserHandler(w http.ResponseWriter, r *http.Request) {
	type UserLogin struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var ul UserLogin
	json.NewDecoder(r.Body).Decode(&ul)

	chiperPassword, err := LoginUserService(r.Context(), ul.Email)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)

		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(chiperPassword), []byte(ul.Password))
	if err != nil {
		writeError(w, http.StatusUnauthorized, errors.New("password is wrong"))

		return
	}

	tokenJwt, err := GenerateJwtTOken([]byte(os.Getenv("JWT_KEY")))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)

		return
	}

	w.Header().Add("content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenJwt,
	})
}

func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, http.StatusUnauthorized, errors.New("missing token"))

		return
	}

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

	tokenJwt, err := VerifyJwtToken(bearerToken)
	if err != nil {
		writeError(w, http.StatusUnauthorized, errors.New("unauthorized"))

		return
	}

	if !tokenJwt.Valid {
		writeError(w, http.StatusUnauthorized, errors.New("cannot create todo"))

		return
	}

	type Todo struct {
		UserId      uuid.UUID `json:"userId"`
		Title       string    `json:"title"`
		Description string    `json:"description"`
	}
	var t Todo
	json.NewDecoder(r.Body).Decode(&t)

	err = CreateTodoService(r.Context(), t.UserId, t.Title, t.Description)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New(uuid.Nil.String()))

		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"userId":      t.UserId,
		"title":       t.Title,
		"description": t.Description,
	})
}

func UpdateTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, http.StatusUnauthorized, errors.New("missing token"))

		return
	}

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

	tokenJwt, err := VerifyJwtToken(bearerToken)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("failed to verify token"))

		return
	}

	if !tokenJwt.Valid {
		writeError(w, http.StatusBadRequest, errors.New("token is invalid"))

		return
	}

	type Todo struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	var t Todo
	json.NewDecoder(r.Body).Decode(&t)

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusForbidden, err)

		return
	}

	err = UpdateTodoService(r.Context(), id, t.Title, t.Description)
	if err != nil {
		writeError(w, http.StatusForbidden, err)

		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":          id,
		"title":       t.Title,
		"description": t.Description,
	})
}

func DeleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, http.StatusUnauthorized, errors.New("missing token"))

		return
	}

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

	tokenJwt, err := VerifyJwtToken(bearerToken)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("failed to verify token"))

		return
	}

	if !tokenJwt.Valid {
		writeError(w, http.StatusBadRequest, errors.New("token is invalid"))

		return
	}

	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusForbidden, err)

		return
	}

	err = DeleteTodoService(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusForbidden, err)

		return
	}

	writeMessage(w, http.StatusNoContent, "todo is deleted")
}

func GetTodosHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		writeError(w, http.StatusUnauthorized, errors.New("missing token"))

		return
	}

	bearerToken := strings.TrimPrefix(authHeader, "Bearer ")

	tokenJwt, err := VerifyJwtToken(bearerToken)
	if err != nil {
		writeError(w, http.StatusBadRequest, errors.New("failed to verify token"))

		return
	}

	if !tokenJwt.Valid {
		writeError(w, http.StatusBadRequest, errors.New("token is invalid"))

		return
	}

	q := r.URL.Query()
	page, err := strconv.Atoi(q.Get("page"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)

		return
	}
	limit, err := strconv.Atoi(q.Get("limit"))
	if err != nil {
		writeError(w, http.StatusBadRequest, err)

		return
	}

	todos, err := GetTodosService(r.Context(), page, limit)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)

		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"todos": todos,
		"page":  page,
		"limit": limit,
		"total": len(todos),
	})
}
