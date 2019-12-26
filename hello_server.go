package main
import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"gopkg.in/natefinch/lumberjack.v2"
)

func handler(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query()
	name := query.Get("name")
	if name == ""{
		name = "Guest"
	}
	log.Printf("Received request for %s\n", name)
	w.Write([]byte(fmt.Sprintf("Hello, %s\n", name)))
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", handler)

	srv := &http.Server{
		Handler: r,
		Addr: ":8080",
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	LogFileLocation := os.Getenv("LOG_FILE_LOCATION")
	if LogFileLocation != "" {
		log.SetOutput(&lumberjack.Logger {
			Filename: LogFileLocation,
			MaxSize: 500,
			MaxBackups: 3,
			MaxAge: 28,
			Compress: true,
		})
	}

	// start server
	go func() {
		log.Println("Starting server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	waitForShutdown(srv)
}

func waitForShutdown(srv *http.Server) {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	
	// Block until we receive a signal
	<-interruptChan

	// Create a deadline to waitfor
	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 10)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting Down")
	os.Exit(0)
}