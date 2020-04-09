package main

import (
	"fmt"
	"net/http"

	"github.com/jihoon6372/hog/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"

	custommiddleware "./middleware"
	"github.com/jihoon6372/hog/config"
	"github.com/jihoon6372/hog/handler"
	"github.com/jihoon6372/hog/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func getJWTConfig(secretKey []byte) middleware.JWTConfig {
	DefaultJWTConfig := middleware.JWTConfig{
		Claims:     &config.JwtCustomClaims{},
		SigningKey: secretKey,
	}

	return DefaultJWTConfig
}

func main() {
	var cfg config.Config
	utils.ReadConfig(&cfg)

	dbinfo := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.Username, cfg.Database.Database, cfg.Database.Password, cfg.Database.SSLMode)
	db, err := gorm.Open("postgres", dbinfo)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Profile{})
	// db.AutoMigrate(&model.PlayList{})

	e := echo.New()
	jwtConfig := getJWTConfig([]byte(cfg.Config.SecretKey))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	h := &handler.Handler{DB: db}
	e.POST("/users", h.UserCreate)
	e.GET("/users", h.FindUsers)
	e.POST("/login", h.Login)
	e.GET("/tracks/:id", h.FindTrack)
	e.PATCH("/tracks/:id", h.UpdateTrack)

	r := e.Group("/users/*")
	r.Use(middleware.JWTWithConfig(jwtConfig))
	r.GET("/me", h.UserRead)
	r.PATCH("/me", h.UserUpdate)
	// e.DELETE("/user/:id", h.UserDelete)

	m := &custommiddleware.TestMiddleware{}
	t := e.Group("/tests")
	t.Use(m)
	t.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
