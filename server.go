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
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const defaultPort = "8080"

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		if err = conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
		}
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, _ := upgrader.Upgrade(w, r, nil)

	reader(ws)

}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()
	router.Use(auth.AuthMiddleware)

	router.Use(cors.New(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders: []string{"Link"},
	}).Handler)

	// dsn := "root:@tcp(127.0.0.1:3306)/tpaweb?charset=utf8mb4&parseTime=True&loc=Local"
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	dsn := "host=localhost user=postgres password=pw dbname=tpaweb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// dsn := "host=localhost user=postgres password=mysecretpassword dbname=tpaweb port=5433 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

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
		&model.Block{},
		&model.Message{},

		&model.Post{},
		&model.PostLike{},
		&model.Comment{},
		&model.CommentLike{},
		&model.Job{},
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
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})

	ws := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	ws.AddTransport(&transport.Websocket{
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				log.Println("ws")
				// Check against your desired domains here
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	})
	// router.Handle("/ws", wsEndpoint)

	// e := router.NewRouter(echo.New(), ws)
	// e.Logger.Fatal(e.Start(":8080"))

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// router.Handle("/ws", wsEndpoint)
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))

}
