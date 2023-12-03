package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
)

type Artist struct {
	ID    string   `json:"id"`    // id коллектива
	Name  string   `json:"name"`  // название группы
	Born  string   `json:"born"`  // год основания группы
	Genre string   `json:"genre"` // жанр
	Songs []string `json:"songs"` // популярные песни, это слайс строк, так как песен может быть несколько
}

var artists = map[string]Artist{
	"1": {
		ID:    "1",
		Name:  "30 Seconds To Mars",
		Born:  "1998",
		Genre: "alternative",
		Songs: []string{
			"The Kill",
			"A Beautiful Lie",
			"Attack",
			"Live Like A Dream",
		},
	},
	"2": {
		ID:    "2",
		Name:  "Garbage",
		Born:  "1994",
		Genre: "alternative",
		Songs: []string{
			"Queer",
			"Shut Your Mouth",
			"Cup of Coffee",
			"Til the Day I Die",
		},
	},
}

/*
	Первый эндпоинт

Теперь можно писать код для реализации первого эндпоинта. Начать следует с обработчика для возврата всех элементов — их у вас два. Для этого клиент должен сделать
запрос `GET /artists`. И тогда вы вернете всех артистов в формате JSON.
Чтобы написать обработчик, нужно создать функцию, которая принимает в качестве параметров запрос `r *http.Request` и ответ `w http.ResponseWriter`. В `r *http.Request`
содержатся данные запроса клиента, а в `w http.ResponseWriter` записываются данные ответа, который отправляется клиенту:
*/
func getArtists(w http.ResponseWriter, r *http.Request) {
	// сериализуем данные из слайса artists
	resp, err := json.Marshal(artists)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// в заголовок записываем тип контента, у нас это данные в формате JSON
	w.Header().Set("Content-Type", "application/json")
	// так как все успешно, то статус OK
	w.WriteHeader(http.StatusOK)
	// записываем сериализованные в JSON данные в тело ответа
	w.Write(resp)
}

func postArtist(w http.ResponseWriter, r *http.Request) {
	var artist Artist
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &artist); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	//artists = append(artists, artist)
	artists[artist.ID] = artist

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getArtist(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	artist, ok := artists[id]
	if !ok {
		http.Error(w, "Артист не найден", http.StatusNoContent)
		return
	}

	resp, err := json.Marshal(artist)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func main() {
	// создаём новый роутер
	r := chi.NewRouter()

	// регистрируем в роутере эндпоинт `/artists` с методом GET, для которого используется обработчик `getArtists`
	r.Get("/artists", getArtists)
	// регистрируем в роутере эндпоинт `/artists` с методом POST, для которого используется обработчик `postArtist`
	r.Post("/artists", postArtist)
	// регистрируем в роутере эндпоинт `/artist/{id}` с методом GET, для которого используется обработчик `getArtist`
	r.Get("/artist/{id}", getArtist)

	// запускаем сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
