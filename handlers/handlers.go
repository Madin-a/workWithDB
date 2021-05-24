package handlers

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"
	"workWithDB/DataBase"
	"workWithDB/models"
)

var db *gorm.DB

func init() {
	db = DataBase.GetDB()
}

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	Login string `json:"login"`
}

func GenerateToken(userLogin string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		userLogin,
	})
	signedToken, err := token.SignedString([]byte(signingKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseToken(accessToken string) (*tokenClaims, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return &tokenClaims{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return &tokenClaims{}, errors.New("token claims are not type of *tokenClaims")
	}

	return claims, nil
}

func ErrorHandler(c *gin.Context, err error) {
	fmt.Println("[Response]: ", err.Error())
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
}

func CreateUser(c *gin.Context) {
	var users models.User
	var newUser models.User
	var validate *validator.Validate

	err := c.ShouldBindJSON(&newUser)
	if err != nil {
		ErrorHandler(c, err)
		return
	}

	fmt.Println("[Request]: ", newUser)

	validate = validator.New()

	err = validate.Struct(&newUser)
	if err != nil {

		ErrorHandler(c, err)
		return

	}

	err = verifyPassword(newUser.Password)
	if err != nil {
		ErrorHandler(c, err)
		return
	}
	newUser.Password, err = HashPassword(newUser.Password)
	if err != nil {
		log.Println("Can't hash the password:", err)
		ErrorHandler(c, err)
		return
	}

	login := newUser.Login
	fmt.Println("---------------------------------------------------------------------------------------------")
	count := db.First(&users, "login = ?", login).RowsAffected

	if count > 0 {
		ErrorHandler(c, errors.New("this login is already exist"))
		return
	} else {
		db.Create(&newUser)
		fmt.Println("[Response]: user added")
		c.JSON(http.StatusOK, gin.H{"message": "user added"})
		return
	}

}

func DeleteUser(c *gin.Context) {
	var user models.User
	id := c.Params.ByName("id")
	if id == "" {
		ErrorHandler(c, errors.New("error"))
	}
	i, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err)
	}

	db.Delete(&user, i)
}

func Entry(c *gin.Context) {
	var userFromDb models.User
	var RequestUser models.LogPas
	err := c.ShouldBindJSON(&RequestUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	login := RequestUser.Login

	count := db.First(&userFromDb, "login=?", login).RowsAffected

	ok := CheckPasswordHash(RequestUser.Password, userFromDb.Password)
	fmt.Println("count: ", count, "ok = ", ok)
	if count == 1 && ok == true {
		token, err := GenerateToken(login)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"token": token})
		return

	} else {
		ErrorHandler(c, errors.New("wrong login or password"))
	}

}

func Something(c *gin.Context) {
	var users models.User

	header := c.GetHeader("token")
	if header == "" {
		c.JSON(http.StatusBadRequest, "empty header")
		return
	}

	token, err := ParseToken(header)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}

	count := db.First(&users, "login=?", token.Login).RowsAffected
	if count == 1 {
		c.JSON(http.StatusOK, users)
		return
	}

	fmt.Println(token.Login)

}

func verifyPassword(password string) error {
	var uppercasePresent bool
	var lowercasePresent bool
	var numberPresent bool
	var specialCharPresent bool
	const minPassLength = 8
	const maxPassLength = 64
	var passLen int
	var errorString string

	for _, ch := range password {
		switch {
		case unicode.IsNumber(ch):
			numberPresent = true
			passLen++
		case unicode.IsUpper(ch):
			uppercasePresent = true
			passLen++
		case unicode.IsLower(ch):
			lowercasePresent = true
			passLen++
		case unicode.IsPunct(ch) || unicode.IsSymbol(ch):
			specialCharPresent = true
			passLen++
		case ch == ' ':
			passLen++
		}
	}
	appendError := func(err string) {
		if len(strings.TrimSpace(errorString)) != 0 {
			errorString += ", " + err
		} else {
			errorString = err
		}
	}
	msg := "password: "
	if !lowercasePresent {
		appendError(msg + "lowercase letter missing")
	}
	if !uppercasePresent {
		appendError(msg + "uppercase letter missing")
	}
	if !numberPresent {
		appendError(msg + "at least one numeric character required")
	}
	if !specialCharPresent {
		appendError(msg + "special character missing")
	}
	if !(minPassLength <= passLen && passLen <= maxPassLength) {
		appendError(fmt.Sprintf("password length must be between %d to %d characters long", minPassLength, maxPassLength))
	}

	if len(errorString) != 0 {
		return fmt.Errorf(errorString)
	}
	return nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
