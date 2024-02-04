package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//Нужно реализовать три теста:

// 1) Запрос сформирован корректно, сервис возвращает код ответа 200
// и тело ответа не пустое.
func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("expected status code: %d, got %d", http.StatusOK, status)
	}
}

// 2) Город, который передаётся в параметре city, не поддерживается.
// Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhenUnsupportedCity(t *testing.T) {
	req, err := http.NewRequest("GET", "/cafe?count=3&city=tula", nil)
	require.NoError(t, err, "Ошибка при создании запроса!")

	responseRecorder := httptest.NewRecorder()

	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusBadRequest, responseRecorder.Code, "Код ответа 400")
	require.Contains(t, responseRecorder.Body.String(), "wrong city value", "Тело ответа содержать ошибку wrong city value!")
}

// 3) Если в параметре count указано больше, чем есть всего,
// должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req, err := http.NewRequest("GET", "/cafe?count=7&city=moscow", nil)
	require.NoError(t, err, "Ошибка при создании запроса!")

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	require.Equal(t, http.StatusOK, responseRecorder.Code, "Код ответа  200")
	require.NotEmpty(t, responseRecorder.Body.String(), "Тело ответа не должно быть пустым")

	responseBody := responseRecorder.Body.String()
	cafes := strings.Split(responseBody, ",")

	require.Len(t, cafes, totalCount, "Количество выведеных кафе, должно быть равно 4!")
}
