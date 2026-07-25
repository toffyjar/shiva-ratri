package model

import (
	"time"

	"github.com/google/uuid"
)

type Kundli struct {
	Name       string `json:"name" binding:"required" bson:"name"`
	Day        string `json:"day" binding:"required" bson:"day"`
	Month      string `json:"month" binding:"required" bson:"month"`
	Year       string `json:"year" binding:"required" bson:"year"`
	Hour       string `json:"hour" binding:"required" bson:"hour"`
	Minute     string `json:"minute" binding:"required" bson:"minute"`
	BirthPlace string `json:"birthPlace" binding:"required" bson:"birth_place"`
	Latitude   string `json:"latitude" binding:"required" bson:"latitude"`
	Longitude  string `json:"longitude" binding:"required" bson:"longitude"`
	TimeZone   string `json:"timeZone" binding:"required" bson:"time_zone"`
	IsFemale   string `json:"isFemale" binding:"required" bson:"is_female"`
}

type KundaliSave struct {
	ID         uint      `json:"id" bson:"_id"`
	UserID     uuid.UUID `json:"user_id" bson:"user_id"`
	Name       string    `json:"name" bson:"name"`
	Day        string    `json:"day" bson:"day"`
	Month      string    `json:"month" bson:"month"`
	Year       string    `json:"year" bson:"year"`
	Hour       string    `json:"hour" bson:"hour"`
	Minute     string    `json:"minute" bson:"minute"`
	BirthPlace string    `json:"birthPlace" bson:"birth_place"`
	Latitude   string    `json:"latitude" bson:"latitude"`
	Longitude  string    `json:"longitude" bson:"longitude"`
	TimeZone   string    `json:"timeZone" bson:"time_zone"`
	IsFemale   string    `json:"isFemale" bson:"is_female"`
	CreatedAt  time.Time `json:"created_at" bson:"created_at"`
}
