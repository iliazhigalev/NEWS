package models

import "gorm.io/gorm"

type Article struct {
	gorm.Model
	Title     string `json:"title" gorm:"text;not null;default:null"`
	Anons     string `json:"anons" gorm:"text;not null;default:null"`
	Full_text string `json:"full_text" gorm:"text;not null;default:null"`
}
