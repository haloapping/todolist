package main

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name" fake:"{name}"`
	Email    string    `json:"email" fake:"{email}"`
	Password string    `json:"password"`
}

type Todo struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"userId"`
	Title       string    `json:"title" fake:"{noun}"`
	Description string    `json:"description" fake:"{sentence:2}"`
}
