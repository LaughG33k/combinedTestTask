package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/LaughG33k/notes/client/psql"
	"github.com/LaughG33k/notes/iternal"
	"github.com/LaughG33k/notes/iternal/handler"
	"github.com/LaughG33k/notes/iternal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattes/migrate/source/file"
	"gopkg.in/yaml.v2"
)

func Run() {

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	iternal.InitLogrus("./logs.json")

	router := chi.NewRouter()

	var cfg iternal.ConfigApp
	if err := initConfig("./config.yaml", &cfg); err != nil {
		iternal.Logger.Error(err)
		return
	}

	if err := startMigrations("file://migrations/psql/notes", cfg.NoteDB); err != nil {
		iternal.Logger.Error(err)
		return
	}

	tmDb, canc := context.WithTimeout(mainCtx, time.Duration(cfg.NoteDB.NewConnTimeoutInSec)*time.Second)
	defer canc()

	db, err := psql.NewClient(tmDb, cfg.NoteDB)

	if err != nil {
		iternal.Logger.Error(err)
		return
	}

	httpServer := &http.Server{
		WriteTimeout: time.Duration(cfg.WriteTimeoutInSec) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeoutInSec) * time.Second,
		IdleTimeout:  time.Duration(cfg.IdleTimeoutInSec) * time.Second,
		Addr:         cfg.Addr,
		Handler:      router,
	}

	noteRepo := &repository.Note{
		Conn: db,
	}

	createNote := &handler.CreateNote{
		NoteRepository: noteRepo,
		Ctx:            mainCtx,
		Timeout:        time.Duration(cfg.WriteTimeoutInSec*80/100) * time.Second,
	}

	getNotes := &handler.GetNotes{
		Ctx:       mainCtx,
		NotesRepo: noteRepo,
		Timeout:   time.Duration(cfg.WriteTimeoutInSec*80/100) * time.Second,
	}

	router.Post("/create", func(w http.ResponseWriter, r *http.Request) {

		ip := strings.Split(r.RemoteAddr, ":")[0]

		for _, v := range cfg.TrustedAddrs {
			if ip == v {
				createNote.Handle(w, r)
				return
			}
		}

		http.Error(w, "not trusted ip", 401)

	})
	router.Get("/get", func(w http.ResponseWriter, r *http.Request) {

		ip := strings.Split(r.RemoteAddr, ":")[0]

		for _, v := range cfg.TrustedAddrs {
			if ip == v {
				getNotes.Handle(w, r)
				return
			}
		}

		http.Error(w, "not trusted ip", 401)

	})

	go httpServer.ListenAndServe()

	<-mainCtx.Done()

	httpServer.Shutdown(context.Background())
	iternal.Logger.Info("app shutdown")

}

func initConfig(path string, dest any) error {

	bytes, err := os.ReadFile(path)

	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(bytes, dest); err != nil {
		return err
	}

	return nil

}

func startMigrations(path string, cfg iternal.DBConfig) error {
	migrations, err := migrate.New(
		path,
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB, "disable"),
	)

	if err != nil {
		return err
	}

	if err := migrations.Up(); err != nil {
		if err.Error() != "no change" {
			migrations.Drop()
			return err
		}

	}

	migrations.Close()
	return nil
}
