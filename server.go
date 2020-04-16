package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jihoon6372/hog/utils"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/middleware"

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
		Skipper: func(c echo.Context) bool {
			parts := strings.Split("header:"+echo.HeaderAuthorization, ":")
			auth := c.Request().Header.Get(parts[1])
			authScheme := "Bearer"
			l := len(authScheme)
			if len(auth) > l+1 && auth[:l] == authScheme {
				return false
			}

			return true
		},
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

	// var testA = config.ReadOrAuthenticatedJWTConfig
	// fmt.Println("testA", testA.SigningKey)

	r := e.Group("/users")
	r.Use(middleware.JWTWithConfig(jwtConfig))
	r.GET("/me", h.UserRead)
	r.PATCH("/me", h.UserUpdate)
	r.GET("/test", h.TestSkipper)
	// e.DELETE("/user/:id", h.UserDelete)

	// m := &custommiddleware.TestMiddleware{}
	t := e.Group("/tests")
	t.Use(testMiddleware)
	t.GET("/:id", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "hello world",
		})
	})

	s := NewStats()
	e.Use(s.Process)
	e.GET("/stats", s.Handle)

	e.GET("/lists", h.FindListTest)

	e.Logger.Fatal(e.Start(":1323"))
}

func testMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		token := new(jwt.Token)

		fmt.Println("token", token)

		// context.WithValue
		// err := next(c)
		id := c.Param("id")

		fmt.Println("id", id)
		if id != "20" {
			return echo.ErrUnauthorized
		}

		return next(c)
	}
}

type (
	// Stats ...
	Stats struct {
		Uptime       time.Time      `json:"uptime"`
		RequestCount uint64         `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		mutex        sync.RWMutex
	}
)

// NewStats ...
func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: map[string]int{},
	}
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("user", "hello?")

		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++
		return nil
	}
}

// Handle is the endpoint to get stats.
func (s *Stats) Handle(c echo.Context) error {
	user := c.Get("user")
	fmt.Println("user", user)
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}
