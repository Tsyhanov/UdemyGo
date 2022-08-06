package main

import (
	"fmt"
	"testing"
	"udemygo/bookings/internal/config"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		//do nothing
	default:
		t.Error(fmt.Sprintf("type is not chi.Mux but %T", v))
	}
}
