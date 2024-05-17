package main

import (
	"fmt"
	"html/template"
	"net/http"
	"news/database"
	"news/models"
)

func index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/index.html", "templates/header.html")

	if err != nil {
		fmt.Fprintf(w, err.Error())
		return
	}

	t.ExecuteTemplate(w, "index", nil) // вторым параметром, мы говорим какой блок html мы будем выводить,третий парамет- параметр которые мы передаём в шаблон
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

	// Если все прошло успешно, отправляем пользователю сообщение об успешном сохранении
	fmt.Fprintf(w, "Article saved successfully!")
}

func handleFunc() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static/")))) // указываем url адреса, которые мы обрабатываем,далее мы говорим в какой папке мы будем искать соответвующий файл
	http.HandleFunc("/", index)
	http.HandleFunc("/create", create)
	http.HandleFunc("/save_article", save_article)
	http.ListenAndServe(":3000", nil)
}

func main() {
	database.ConnectDb()
	handleFunc()
}
