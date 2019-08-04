package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/chuabingquan/snippets/bcrypt"
	"github.com/chuabingquan/snippets/http"
	"github.com/chuabingquan/snippets/http/jwt"
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

	jwtAuthenticator := jwt.Authenticator{
		SigningKey: []byte(config["AUTH_SECRET"]),
		ExpiryTime: time.Duration(toInt(config["AUTH_EXPIRY"])) * time.Minute,
	}
	hu := bcrypt.Utilities{HashCost: toInt(config["HASH_COST"])}

	us := postgres.UserService{DB: db, HashUtilities: hu}
	ss := postgres.SnippetService{DB: db}
	as := postgres.AuthenticationService{DB: db, HashUtilities: hu}

	userHandler := http.NewUserHandler(us, jwtAuthenticator)
	snippetHandler := http.NewSnippetHandler(ss, jwtAuthenticator)
	authHandler := http.NewAuthHandler(as, us, jwtAuthenticator)

	handler := http.Handler{
		UserHandler:    userHandler,
		SnippetHandler: snippetHandler,
		AuthHandler:    authHandler,
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

func toInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal("Error parsing \"" + str + "\" to int: " + err.Error())
	}
	return val
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
		"PORT", "HASH_COST", "AUTH_SECRET", "AUTH_EXPIRY"}
	for _, name := range envNames {
		val, ok := os.LookupEnv(name)
		if !ok {
			log.Fatal(name + " environment variable is required but not set")
		}
		config[name] = val
	}
	return config
}
