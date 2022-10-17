package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

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

type chatConnSockets struct {
	id   string
	conn *websocket.Conn
}

type typeChatSocket struct {
	id      string
	content string
}

var upgrader = websocket.Upgrader{}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	conn, _ := upgrader.Upgrade(w, r, nil)
	// conns = append(conns, conn)

	for {
		// _, p, _ := conn.ReadMessage()
		// log.Println(p)

		// var chatSocket = make(map[string][string])
		var chatSocket interface{}
		var typedChatSocket = *&typeChatSocket{
			id:      "",
			content: "",
		}

		conn.ReadJSON(&chatSocket)

		val := reflect.ValueOf(chatSocket).Elem()
		n := val.FieldByName("UserEmail").Interface().(string)
		fmt.Printf("%+v\n", n)

		log.Println(chatSocket)
		// log.Println(chatSocket.(typeChatSocket))
		log.Println(typedChatSocket)
		// if err != nil {
		// 	log.Println(err)
		// }

		// for _, _conn := range conns {
		// 	if err = _conn.WriteMessage(messageType, p); err != nil {
		// 		log.Println(err)
		// 	}
		// }

	}

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
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	router.HandleFunc("/ws", wsEndpoint)
	router.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
