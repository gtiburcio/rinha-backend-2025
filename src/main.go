package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"rinha-backend-2025-gtiburcio/src/client"
	"rinha-backend-2025-gtiburcio/src/config"
	"rinha-backend-2025-gtiburcio/src/handler"
	"rinha-backend-2025-gtiburcio/src/job"
	"rinha-backend-2025-gtiburcio/src/repository"
	"rinha-backend-2025-gtiburcio/src/usecase"

	"github.com/joho/godotenv"
)

func main() {
	loadEnvs()

	dbConfig := config.NewDBConfig()

	c := client.NewClient()
	r := repository.NewRepository(dbConfig.DBConn)
	concurrentSlice := repository.NewConcurrentSlice(r)
	u := usecase.NewUseCase(c, r, concurrentSlice)

	j := job.NewJob(u)
	j.Run()

	h := handler.NewHandler(j, u)

	http.HandleFunc("/payments", h.HandleSavePayment)
	http.HandleFunc("/payments-summary", h.HandlePaymentSummary)

	fmt.Println("Server starting on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func loadEnvs() {
	if len(os.Args) > 1 && os.Args[1] == "local" {
		log.Default().Print("Running local")
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
}
