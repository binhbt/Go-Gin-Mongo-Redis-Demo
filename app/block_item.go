package main

import "go.mongodb.org/mongo-driver/bson/primitive"

// Task - Model of a basic task
type BlockItem struct {
	ID          primitive.ObjectID
	UserID      int
	Token       string
	PayLoad     string
	BlockType   string
	ExpiredTime int
}
