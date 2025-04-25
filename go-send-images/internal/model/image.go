package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Image struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`  // Mongo ID (автоматически)
	Filename string             `bson:"filename" json:"filename"` // Название файла
	URL      string             `bson:"url" json:"url"`           // Путь вида /uploads/filename.png
	Modified int64              `bson:"modified" json:"modified"` // Временная метка
}
