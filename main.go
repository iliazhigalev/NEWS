package main

import (
	"fmt"
	"html/template"
	"net/http"
	"news/database"
	"news/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func index(w http.ResponseWriter, r *http.Request) {
	var articles []models.Article
	result := database.DB.Db.Find(&articles)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	// Log the fetched articles
	fmt.Printf("Fetched articles: %+v\n", articles)

	t, err := template.ParseFiles("templates/index.html", "templates/header.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "index", articles)
}

func create(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/create.html", "templates/header.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "create", nil) // вторым параметром, мы говорим какой блок html мы будем выводить,третий парамет- параметр которые мы передаём в шаблон
}

func save_article(w http.ResponseWriter, r *http.Request) {

	title := r.FormValue("title") //указываем название поля
	anons := r.FormValue("anons")
	full_text := r.FormValue("full_text")
	if title == "" || anons == "" || full_text == "" {
		fmt.Fprintf(w, "Не все данные заполнены")
	} else {

		// Создаем экземпляр статьи
		article := models.Article{
			Title:     title,
			Anons:     anons,
			Full_text: full_text,
		}

		// Сохраняем статью в базе данных
		result := database.DB.Db.Create(&article)
		if result.Error != nil {
			// Если произошла ошибка при сохранении в базу данных, отправляем ошибку клиенту
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}
		// Log the saved article
		fmt.Printf("Article saved: %+v\n", article)
		// Если все прошло успешно, отправляем пользователю сообщение об успешном сохранении
		fmt.Fprintf(w, "Article saved successfully!")
	}
}
func show_post(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		http.Error(w, "Missing id in the URL", http.StatusBadRequest)
		return
	}

	var showPost models.Article
	result := database.DB.Db.Where("id = ?", id).First(&showPost)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			http.Error(w, "Article not found", http.StatusNotFound)
		} else {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		}
		return
	}

	t, err := template.ParseFiles("templates/show.html", "templates/header.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.ExecuteTemplate(w, "show", showPost); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleFunc() {
	rtr := mux.NewRouter()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))) // указываем url адреса, которые мы обрабатываем,далее мы говорим в какой папке мы будем искать соответвующий файл
	rtr.HandleFunc("/", index).Methods("GET")                                                     //указываем метод передачи данных
	rtr.HandleFunc("/create", create).Methods("GET")
	rtr.HandleFunc("/save_article", save_article).Methods("POST")
	rtr.HandleFunc("/post/{id:[0-9]+}", show_post).Methods("GET")

	http.Handle("/", rtr) // указываем что обработка всех url адресов будет происходить через объект роутер
	http.ListenAndServe(":3000", nil)
}

func main() {
	database.ConnectDb()
	handleFunc()
}
