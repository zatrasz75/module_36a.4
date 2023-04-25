package main

import (
	api "GoNews/pkg/api"
	"GoNews/pkg/rss"
	storage "GoNews/pkg/storage"
	db "GoNews/pkg/storage/db"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// сервер GoNews
type server struct {
	db  storage.Interface
	api *api.API
}

// Конфигурация приложения
type config struct {
	Period  int      `json:"request_period"`
	LinkArr []string `json:"rss"`
}

func main() {

	// объект сервера
	var srv server

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	// объект базы данных postgresql
	db, err := db.New(ctx, "postgres://postgres:rootroot@localhost:5432/aggregator")
	if err != nil {
		log.Fatal(err)
	}

	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	// чтение и раскодированное файла конфигурации
	b, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	// каналы для обработки новостей и ошибок
	chanPosts := make(chan []storage.Post)
	chanErrs := make(chan error)

	// получение и парсинг ссылок из конфига
	myLinks := getRss("config.json", chanErrs)
	for i := range myLinks.LinkArr {
		go parseNews(myLinks.LinkArr[i], chanErrs, chanPosts, config.Period)
	}

	// обработка потока новостей
	go func() {
		for posts := range chanPosts {
			for i := range posts {
				db.AddPost(posts[i])
			}
		}
	}()

	// обработка потока ошибок
	go func() {
		for err := range chanErrs {
			log.Println("ошибка:", err)
		}
	}()

	// запуск веб-сервера с API и приложением
	err = http.ListenAndServe(":80", srv.api.Router())
	if err != nil {
		log.Fatal(err)
	}

}

// Получающая отдельные ссылки из конфига, ошибки направляются в поток ошибок
func getRss(fileName string, errors chan<- error) config {
	jsonFile, err := os.Open(fileName)
	if err != nil {
		errors <- err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var links config

	json.Unmarshal(byteValue, &links)

	return links
}

// Получение новостей по ссылкам и отправка новостей и ошибок в соответствующие каналы
func parseNews(link string, errs chan<- error, posts chan<- []storage.Post, period int) {
	for {
		newsPosts, err := rss.RssToStruct(link)
		if err != nil {
			errs <- err
			continue
		}
		posts <- newsPosts
		time.Sleep(time.Minute * time.Duration(period))
	}
}
