package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var db *gorm.DB

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"unique"`
	Password string
}

type Game struct {
	ID             uint `gorm:"primary_key"`
	TargetNumber   int
	Attempts       int
	CorrectGuesses int
}

func initDatabase() {
	var err error
	db, err = gorm.Open("mysql", "crolord:RyV3MGZ$@Q5rJ3i^-=@tcp(rm-j6cn3wen02w6f5b94ho.mysql.rds.aliyuncs.com:3306)/crolord?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("failed to connect to database")
	}

	// 创建数据库表
	db.AutoMigrate(&User{}, &Game{})
}
func (Game) TableName() string {
	return "game"
}

func getOrCreateGame(user *User) (*Game, error) {
	var game Game
	if err := db.Where("id = ?", user.ID).First(&game).Error; err != nil {
		log.Println("No game record found for user:", user.ID)

		if gorm.IsRecordNotFoundError(err) {
			game.ID = user.ID // 修改这一行
			game.TargetNumber = generateTargetNumber()
			game.Attempts = 0
			if err := db.Create(&game).Error; err != nil {
				return nil, err
			}
		} else {
			log.Println("Error querying game record:", err)

			return nil, err
		}
	}
	return &game, nil
}

func incrementAttempts(game *Game) {
	game.Attempts++
	db.Save(game)
}

func getUserFromAuthToken(authToken string, userID uint) (User, error) {
	userIDURL := fmt.Sprintf("http://localhost:8083/user?authToken=%s&userID=%d", authToken, userID)
	resp, err := http.Get(userIDURL)
	if err != nil {
		return User{}, fmt.Errorf("error sending request to login service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("login service returned status %d", resp.StatusCode)
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		return User{}, fmt.Errorf("error decoding user JSON: %w", err)
	}
	return user, nil
}

func getUserIDFromLoginService(authToken string) (uint, error) {
	loginServiceURL := "http://localhost:8083"
	requestURL := fmt.Sprintf("%s/user?authToken=%s", loginServiceURL, authToken)
	fmt.Printf("Requesting user ID with URL: %s\n", requestURL)

	resp, err := http.Get(requestURL)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Response body from login service: %s\n", string(respBody))
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("login service returned status %d", resp.StatusCode)
	}

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return uint(data["id"].(float64)), nil
}

func generateTargetNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100)
}
func closeDatabase() {
	db.Close()
}
