package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User is the  data structure for the users table
type User struct {
	ID       int `gorm:"primary-key;AUTO_INCREMENT"`
	Username string
	Password string
	Name     string
	Mobile   string
	Email    string
}

func main() {
	// listen on localhost:LISTEN_PORT
	listenPort := getEnv("LISTEN_PORT", "3333")
	serviceName := getEnv("SERVICE_NAME", "No Name")
	dbHost := getEnv("DB_HOST", "localhost")

	addr := fmt.Sprintf(":%s", listenPort)

	// database
	dbaddr := fmt.Sprintf("postgresql://tester@%s:5432/rpapoc?sslmode=disable", dbHost)
	db, err := gorm.Open("postgres", dbaddr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// http router
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// TODO: Investigate how to configure this for the reverse proxy environment
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Token-Claim-Username"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		// Save a copy of this request for debugging.
		// requestDump, err := httputil.DumpRequest(r, true)
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// fmt.Println(string(requestDump))
		username := r.Header.Get("Token-Claim-Username")
		fmt.Println("username: ", username)

		if username != "" {
			var user User
			error := db.Where("username = ?", username).First(&user).Error
			if error != nil {
				// user not found
				fmt.Println(error)
				http.Error(w, "User Not Found", http.StatusNotFound)
			} else {
				// success
				responseMap := map[string]interface{}{
					"serviceName": serviceName,
					"username":    username,
					"user":        user,
				}
				jr, err := json.Marshal(responseMap)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(jr)
			}
		} else {
			errorMsg := fmt.Sprintf("Error: service %s didn't receive a username in the header", serviceName)
			http.Error(w, errorMsg, http.StatusInternalServerError)
		}
	})

	fmt.Printf("Starting server %s on %v\n", serviceName, addr)
	http.ListenAndServe(addr, r)
}

// this is a helper func to fetch environment variables
func getEnv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}
