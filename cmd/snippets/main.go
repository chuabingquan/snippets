package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/chuabingquan/snippets/postgres"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	loadEnv()
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

	// Block until an interrupt signal is received to keep server alive
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	s := <-c
	fmt.Println("Got signal:", s)
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func getConfig() map[string]string {
	config := make(map[string]string)
	envNames := []string{"DB_PROTOCOL", "DB_USER", "DB_PASSWORD", "DB_HOST", "DB_PORT", "DB_NAME", "DB_SSLMODE"}
	for _, name := range envNames {
		val, ok := os.LookupEnv(name)
		if !ok {
			log.Fatal(name + " environment variable is required but not set")
		}
		config[name] = val
	}
	return config
}
