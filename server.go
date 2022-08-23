package main

import (
	// "gorm.io/driver/mysql"

	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	auth "MZ221-TPA-Web-Back/auth"
	"MZ221-TPA-Web-Back/graph"
	"MZ221-TPA-Web-Back/graph/generated"
	"MZ221-TPA-Web-Back/graph/model"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(auth.AuthMiddleware)

	// Add CORS middleware around every request
	// See https://github.com/rs/cors for full option listing
	router.Use(cors.New(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}).Handler)

	// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	// dsn := os.Getenv("DATABASE_URL")
	dsn := "root:@tcp(127.0.0.1:3306)/tpaweb?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&model.ConnectRequest{},
		&model.Connection{},
		&model.Activation{},
		&model.Education{},
		&model.Experience{},

		&model.User{},
		&model.Activity{},
		&model.Reset{},
		&model.UserFollow{},
		&model.UserVisit{},
		&model.UserEducation{},
		&model.UserExperience{},
		&model.Message{},

		&model.Post{},
		&model.PostLike{},
		&model.Comment{},
		&model.CommentLike{},
	)

	c := generated.Config{Resolvers: &graph.Resolver{
		DB: db,
	}}
	c.Directives.Auth = auth.Auth

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(c))
	srv.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// Check against your desired domains here
				return r.Host == "example.org"
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	router.Handle("/",
		// middleware(
		playground.Handler("GraphQL playground", "/query"))
	// )
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
