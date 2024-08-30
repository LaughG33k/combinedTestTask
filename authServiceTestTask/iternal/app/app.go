package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LaughG33k/authServiceTestTask/client/psql"
	"github.com/LaughG33k/authServiceTestTask/iternal"
	"github.com/LaughG33k/authServiceTestTask/iternal/handler"
	"github.com/LaughG33k/authServiceTestTask/iternal/repository"
	"github.com/LaughG33k/authServiceTestTask/pkg"
	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt"
	"gopkg.in/yaml.v2"

	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/mattes/migrate/source/file"
)

func Run() {

	iternal.InitLogrus("./logs.json")

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	appCfg := iternal.ConfigApp{}

	if err := initConfig("./cfg.yaml", &appCfg); err != nil {
		iternal.Logger.Error(err)
		return
	}

	if err := startMigrations("file://migrations/psql/auth", appCfg.AuthDB); err != nil {
		iternal.Logger.Error(err)
		return
	}

	jwtGenerator, err := initJwtWorker("./keys/rsaprivkey.pem")

	if err != nil {
		iternal.Logger.Error(err)
		return
	}

	jwtGenerator.TokenTimelife = time.Duration(appCfg.JwtTimelifeInSec) * time.Second

	jwtParser, err := initJwtWorker("./keys/public.pem")

	if err != nil {
		iternal.Logger.Error(err)
		return
	}

	router := chi.NewRouter()

	httpSrv := &http.Server{
		Addr:         appCfg.Addr,
		ReadTimeout:  time.Duration(appCfg.ReadTimeoutInSec) * time.Second,
		WriteTimeout: time.Duration(appCfg.ReadTimeoutInSec) * time.Second,
		IdleTimeout:  time.Duration(appCfg.IdleTimeoutInSec) * time.Second,
		Handler:      router,
	}

	dbTm, canc := context.WithTimeout(mainCtx, time.Duration(appCfg.AuthDB.NewConnTimeoutInSec))

	defer canc()

	dbConn, err := psql.NewClient(dbTm, appCfg.AuthDB)

	if err != nil {
		iternal.Logger.Error(err)
		return
	}

	defer dbConn.Close()

	userRepo := &repository.UserRepository{Conn: dbConn}
	rtRepo := &repository.RefreshSessionRepository{Conn: dbConn}

	login := &handler.LoginHandler{
		Ctx:          mainCtx,
		UserRepo:     userRepo,
		RefreshRepo:  rtRepo,
		Timeout:      time.Duration(appCfg.WriteTimeoutInSec*80/100) * time.Second,
		JwtGenerator: jwtGenerator,
	}

	register := &handler.RegHandler{
		UserRepo: userRepo,
		Ctx:      mainCtx,
		Timeout:  time.Duration(appCfg.WriteTimeoutInSec*80/100) * time.Second,
	}

	updateToken := &handler.UpdateTokensHandler{
		RerfreshRepo: rtRepo,
		Ctx:          mainCtx,
		Timeout:      time.Duration(appCfg.WriteTimeoutInSec*80/100) * time.Second,
		JwtGenerator: jwtGenerator,
		JwtParser:    jwtParser,
	}

	sharePublicKey := &handler.SharePublicKey{
		Path:               "./keys/public.pem",
		KeyTimelifeInCache: 10 * time.Minute,
		Method:             jwt.SigningMethodRS512.Alg(),
	}

	logout := &handler.Logout{
		RefreshRepo: rtRepo,
		Ctx:         mainCtx,
		JwtParser:   jwtParser,
		Timeout:     time.Duration(appCfg.WriteTimeoutInSec*80/100) * time.Second,
	}

	router.Post("/login", login.Handle)
	router.Post("/register", register.Handle)
	router.Post("/update", updateToken.Handle)
	router.Get("/publicKey", sharePublicKey.Handle)
	router.Post("/logout", logout.Handle)

	go httpSrv.ListenAndServe()

	<-mainCtx.Done()

	httpSrv.Shutdown(context.Background())

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

func initJwtWorker(path string) (*pkg.JwtWorker, error) {

	bytes, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	worker, err := pkg.NewJwtWorker(bytes, jwt.SigningMethodRS512)

	if err != nil {
		return nil, err
	}

	return worker, err

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
