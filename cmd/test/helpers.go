package main

import (
	"github.com/charmbracelet/log"
	"github.com/cooperstandard/NetZero/internal/routes"
	"net/http"
)

func createAndLoginUser(client *http.Client, email, password, name string) routes.User {
	user, err := register(client, registerParameters{
		Email:    email,
		Password: password,
		Name:     name,
	})
	if err != nil {
		log.Fatal("register user failed", "error", err)
	}
	log.Info("created user", "email", user.Email, "ID", user.ID)

	// 03) login
	user, status := login(client, loginParameters{
		Email:         email,
		Password:      password,
		ExpiresInSecs: 100,
	})
	if status != 200 {
		log.Fatal("failed to login user", "email", email, "error", err)
	}
	log.Info("login successful for", "email", user.Email, "ID", user.ID)

	return user
}
