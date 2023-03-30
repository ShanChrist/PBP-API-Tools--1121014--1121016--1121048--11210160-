package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/redis/go-redis/v9"
	"gopkg.in/gomail.v2"
)

var redisClient *redis.Client
var ctx = context.Background()

func Login(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	err := r.ParseForm()
	if err != nil {
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	rows, err := db.Query("SELECT * FROM user WHERE email = ? AND password = ?",
		email,
		password,
	)

	var user User
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err)
			return
		}
	}

	if user.ID == 0 {
		log.Println(err)
		SendErrorResponse(w, r, "email atau password salah")
		return
	} else {
		generateToken(w, user.ID, user.Name, user.UserType)
		SendRespondDoang(w, r, "berhasil login, selamat datang kembali "+user.Name)

		connectRedis(ctx)
		idString := strconv.Itoa(user.ID)
		setToRedis(ctx, "id", idString)
		setToRedis(ctx, "name", user.Name)
		setToRedis(ctx, "email", user.Email)
	}
}

func Logout(w http.ResponseWriter, r *http.Request) {
	resetUserToken(w)
	SendRespondDoang(w, r, "berhasil logout")
	setToRedis(ctx, "id", "")
	setToRedis(ctx, "name", "")
	setToRedis(ctx, "email", "")
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	email := getFromRedis(ctx, "email")
	if email == "" {
		SendRespondDoang(w, r, "Silahkan Login dulu")
		return
	}
	SendEmailForgotPassword(w, r, email)
	SendRespondDoang(w, r, "Silahkan cek email untuk reset password")
}

func connectRedis(ctx context.Context) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	redisClient = client
}

func setToRedis(ctx context.Context, key, val string) {
	err := redisClient.Set(ctx, key, val, 0).Err()
	if err != nil {
		fmt.Println(err)
	}
}

func getFromRedis(ctx context.Context, key string) string {
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return ""
	}
	return val
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM user"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var user User
	var users []User

	found := false
	for rows.Next() {
		found = true
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err)
			return
		} else {
			users = append(users, user)
		}
	}
	if !found {
		SendErrorResponse(w, r, "Data Not Found")
		return
	}

	if len(users) < 1 {
		SendErrorResponse(w, r, "Error Array Size Not Correct")
	} else {
		userResponse(w, users)
	}
}

func Cron(w http.ResponseWriter, r *http.Request) {
	db := connect()
	defer db.Close()

	SendRespondDoang(w, r, "berhasil")

	c := gocron.NewScheduler(time.UTC)
	c.Every(1).Month().Do(SendEmailSubsMonthly)
	// c.Every(10).Second().Do(SendEmailSubsMonthly)
	c.StartAsync()
}

func SendEmailForgotPassword(w http.ResponseWriter, r *http.Request, email string) {
	db := connect()
	defer db.Close()

	stmt, err := db.Prepare("SELECT name FROM user WHERE email = ?")
	if err != nil {
		SendErrorResponse(w, r, "Something Went Wrong")
		return
	}
	defer stmt.Close()

	var name string
	err = stmt.QueryRow(email).Scan(&name)
	if err != nil {
		SendErrorResponse(w, r, "Something Went Wrong ")
		return
	}

	m := gomail.NewDialer("smtp-mail.outlook.com", 587, "hohohihehooyah@outlook.com", "Hohohooyah")
	mail := gomail.NewMessage()

	mail.SetHeader("From", "hohohihehooyah@outlook.com")
	mail.SetHeader("To", email)
	mail.SetHeader("Subject", "Reset Password")
	mail.SetBody("text/html", "Hello, <b>"+name+"</b><br> Silahkan klik link dibawah ini untuk reset password <br><a href='chess.com'>klik link dibawah sini</a>")

	if err := m.DialAndSend(mail); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

func SendEmailSubsMonthly() {
	db := connect()
	defer db.Close()

	query := "SELECT * FROM user"

	rows, err := db.Query(query)
	if err != nil {
		log.Println(err)
		return
	}

	var user User
	var names []string
	var emails []string
	for rows.Next() {
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.UserType); err != nil {
			log.Println(err)
			return
		} else {
			emails = append(emails, user.Email)
			names = append(names, user.Name)
		}
	}

	log.Println(emails)
	emailSender := "ganjarhooyah@outlook.com"
	m := gomail.NewDialer("smtp-mail.outlook.com", 587, emailSender, "Hohohooyah")

	var wg sync.WaitGroup
	for i := 0; i < len(emails); i++ {
		wg.Add(1)
		go func(email string, name string) {
			defer wg.Done()

			mail := gomail.NewMessage()
			mail.SetHeader("From", emailSender)
			mail.SetHeader("To", email)
			mail.SetHeader("Subject", "Subscription")
			mail.SetBody("text/html", "Hi "+name+", its your lucky day, we are offering you a special price of $9.99 Subscription for the next 3 months. <br>Check the link below for more info <br><a href='chess.com'>click link here</a><br> Sincerely, <br> hohohihehooyah")
			if err := m.DialAndSend(mail); err != nil {
				fmt.Println(err)
				return
			}
			log.Println(email)
		}(emails[i], names[i])
	}
	wg.Wait()
}
