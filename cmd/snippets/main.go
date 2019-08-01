package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

	"github.com/chuabingquan/snippets/bcrypt"
	"github.com/chuabingquan/snippets/http"
	"github.com/chuabingquan/snippets/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	loadEnvironmentVariables()
	config := getConfig()

	dbURL := postgres.DBUrl{
		Protocol: config["DB_PROTOCOL"],
		User:     config["DB_USER"],
		Password: config["DB_PASSWORD"],
		Host:     config["DB_HOST"],
		Port:     config["DB_PORT"],
		Name:     config["DB_NAME"],
		Sslmode:  config["DB_SSLMODE"],
	}.GetURL()

	db, err := postgres.Open(dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	hashCost, err := strconv.Atoi(config["HASH_COST"])
	if err != nil {
		log.Fatal("Failed to parse hash cost environment variable: " + err.Error())
	}

	hu := bcrypt.Utilities{HashCost: hashCost}
	us := postgres.UserService{DB: db, HashUtilities: hu}
	userHandler := http.NewUserHandler(us)

	ss := postgres.SnippetService{DB: db}
	snippetHandler := http.NewSnippetHandler(ss)

	handler := http.Handler{
		UserHandler:    userHandler,
		SnippetHandler: snippetHandler,
	}

	server := http.Server{Handler: &handler, Addr: ":" + config["PORT"]}
	err = server.Open()
	if err != nil {
		log.Fatal("Failed to start server:", err.Error())
	} else {
		log.Println("Server is running")
	}

	// Block until an interrupt signal is received to keep server alive
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	fmt.Println("Got signal:", s)
}

func loadEnvironmentVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getConfig() map[string]string {
	config := make(map[string]string)
	envNames := []string{"DB_PROTOCOL", "DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME", "DB_SSLMODE",
		"PORT", "HASH_COST"}
	for _, name := range envNames {
		val, ok := os.LookupEnv(name)
		if !ok {
			log.Fatal(name + " environment variable is required but not set")
		}
		config[name] = val
	}
	return config
}
