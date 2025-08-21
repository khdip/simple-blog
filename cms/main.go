package main

import (
	"blog/cms/handler"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"google.golang.org/grpc"

	cpb "blog/gunk/v1/category"
	ppb "blog/gunk/v1/post"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("cms/env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Printf("Error loading configuration: %v", err)
	}

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	store := sessions.NewCookieStore([]byte(config.GetString("session.secret")))
	conn, err := grpc.Dial(
		fmt.Sprintf("%s:%s", config.GetString("blog.host"), config.GetString("blog.port")),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("Connection failed", err)
	}

	pc := ppb.NewPostServiceClient(conn)
	cc := cpb.NewCategoryServiceClient(conn)
	r := handler.GetHandler(decoder, store, pc, cc)

	host, port := config.GetString("server.host"), config.GetString("server.port")
	log.Println("Server  starting...")
	if err := http.ListenAndServe(fmt.Sprintf("%s:%s", host, port), r); err != nil {
		log.Fatal(err)
	}
}
