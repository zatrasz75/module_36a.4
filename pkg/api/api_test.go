package api

import (
	storage "GoNews/pkg/storage"
	db "GoNews/pkg/storage/db"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestAPI_postsHandler(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// Создаём чистый объект API для теста.
	dbase, _ := db.New(ctx, "postgres://postgres:rootroot@localhost:5432/aggregator")
	dbase.AddPost(storage.Post{})

	api := New(dbase)

	// Создаём HTTP-запрос.
	req := httptest.NewRequest(http.MethodGet, "/news/10", nil)

	// Создаём объект для записи ответа обработчика.
	rr := httptest.NewRecorder()

	// Вызываем маршрутизатор. Маршрутизатор для пути и метода запроса
	// вызовет обработчик. Обработчик запишет ответ в созданный объект.
	api.r.ServeHTTP(rr, req)
	// Проверяем код ответа.
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен: получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	// Читаем тело ответа.
	b, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Раскодируем JSON в массив заказов.
	var data []storage.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать ответ сервера: %v", err)
	}
	// Проверяем, что в массиве ровно один элемент.
	const wantLen = 10
	if len(data) != wantLen {
		t.Fatalf("получено %d записей, ожидалось %d", len(data), wantLen)
	}
	// Также можно проверить совпадение заказов в результате
	// с добавленными в БД для теста.
}
