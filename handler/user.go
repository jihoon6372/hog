package handler

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jihoon6372/hog/config"
	"github.com/jihoon6372/hog/model"
	"github.com/jihoon6372/hog/utils"
	"github.com/labstack/echo"
	"github.com/lib/pq"
)

// ResultMessage ...
type ResultMessage struct {
	Message string
}

// ResultToken ...
type ResultToken struct {
	Token string
}

// Content ...
type Content struct {
	Id        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserRead ...
func (h *Handler) UserRead(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*config.JwtCustomClaims)
	email := claims.Email

	u := &model.AuthUser{}
	h.DB.Where(model.AuthUser{Email: email}).Find(&u)

	return c.JSON(http.StatusOK, &Content{
		Id:        u.ID,
		Email:     u.Email,
		Username:  u.Username,
		CreatedAt: u.CreatedAt,
	})
}

// UserCreate ...
func (h *Handler) UserCreate(c echo.Context) error {
	user := &model.AuthUser{}
	email := c.FormValue("email")
	username := c.FormValue("username")
	password := c.FormValue("password")

	if email == "" {
		resultError := ResultMessage{Message: "email is require"}
		c.Logger().Error(resultError)
		return c.JSON(http.StatusBadRequest, resultError)
	}

	if password == "" {
		resultError := ResultMessage{Message: "password is require"}
		c.Logger().Error(resultError)
		return c.JSON(http.StatusBadRequest, resultError)
	}

	// 비밀번호 암호화
	hashPassword, _ := utils.HashPassword(password)

	user.Username = username
	user.Email = email
	user.Password = hashPassword

	result := h.DB.Create(&user)
	if err, ok := result.Error.(*pq.Error); ok && err.Code.Name() == "unique_violation" {
		resultError := ResultMessage{Message: "Is Already User"}
		c.Logger().Error(resultError)
		return c.JSON(http.StatusBadRequest, resultError)
	}

	if result.Error != nil {
		resultError := ResultMessage{Message: "Error"}
		c.Logger().Error(resultError)
		return c.JSON(http.StatusBadRequest, resultError)
	}

	return c.JSON(http.StatusCreated, user)
}

// UserUpdate ...
func (h *Handler) UserUpdate(c echo.Context) error {
	requestUser := c.Get("user").(*jwt.Token)
	claims := requestUser.Claims.(*config.JwtCustomClaims)

	username := c.FormValue("username")
	password := c.FormValue("password")
	confirmPassword := c.FormValue("confirm_password")

	user := &model.AuthUser{}
	h.DB.Where(model.AuthUser{Email: claims.Email}).Find(&user)

	// 이름 변경 요청
	if username != "" {
		user.Username = username
	}

	// 비밀번호 변경 요청
	if password != "" && confirmPassword != "" {
		// 비밀번호가 서로 다를때
		if password != confirmPassword {
			return c.JSON(http.StatusBadRequest, ResultMessage{Message: "password is not matched"})
		}

		hashPassword, _ := utils.HashPasswordPbkdf2Sha256(password)
		user.Password = hashPassword
	}

	h.DB.Model(&user).Updates(&user)

	return c.JSON(http.StatusOK, &Content{
		Id:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
	})
}

// UserDelete ...
func (h *Handler) UserDelete(c echo.Context) error {
	return c.JSON(http.StatusNoContent, nil)
}

// Login ...
func (h *Handler) Login(c echo.Context) error {
	user := &model.AuthUser{}
	email := c.FormValue("email")
	password := c.FormValue("password")

	// select
	h.DB.Where(model.AuthUser{Email: email}).Find(&user)

	// password check
	match := utils.ComparePassword(password, user.Password)
	if match != true {
		return c.JSON(http.StatusUnauthorized, ResultMessage{Message: "Invalid Password"})
	}

	// Set custom claims
	claims := &config.JwtCustomClaims{
		Email:    user.Email,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	var cfg config.Config
	utils.ReadConfig(&cfg)

	t, err := token.SignedString([]byte(cfg.Config.SecretKey))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, ResultToken{Token: t})
}
