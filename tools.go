package main

import (
	_ "golang-graphql-subscriptions/graph"
	_ "golang-graphql-subscriptions/graph/generated"
	_ "golang-graphql-subscriptions/infrastructure/router"

	_ "github.com/99designs/gqlgen"
	_ "github.com/golang-jwt/jwt"
	_ "github.com/google/uuid"
	_ "github.com/gorilla/websocket"
	_ "github.com/samber/lo"
	_ "github.com/trycourier/courier-go/v2"
	_ "gorm.io/driver/postgres"
)
