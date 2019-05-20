package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"log"
	"net/http"
)


import "github.com/fvbock/endless"

type User struct {
	User string `db:"user" json:"user"`
}

func main() {
	e := echo.New()
	logger := e.Logger
	db := MustNewMysql()
	//e.Use(middleware.Logger())
	//e.Use(middleware.Recover())
	//e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.Router().Add("GET", "/", func(context echo.Context) error {

		rows, err := db.Queryx("select user from user")
		if err != nil {
			logger.Error(err.Error())
			return context.JSON(http.StatusBadRequest, NewErrorCommonResp(502, err.Error()))
		}
		defer rows.Close()
		var users []*User
		for rows.Next() {
			user := new(User)
			err = rows.StructScan(user)
			if err != nil {
				logger.Error(err.Error())
				return context.JSON(http.StatusBadRequest, NewErrorCommonResp(5001, err.Error()))
			}
			users = append(users, user)
		}
		return context.JSON(http.StatusOK, NewSuccessCommonResp("ok", users))
	})
	e.HTTPErrorHandler = func(err error, c echo.Context) {
		e.Logger.Warn(err)
		_ = c.JSON(http.StatusOK, NewErrorCommonResp(5000, "service err"))
	}
	if err := endless.ListenAndServe("0.0.0.0:8080", e); err != nil {
		log.Fatal(err)
	}
}

func MustNewMysql() *sqlx.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/mysql?charset=utf8&parseTime=True&loc=Local"
	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(50)
	return db
}

type CommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func NewErrorCommonResp(code int, msg string) *CommonResp {
	return &CommonResp{
		Code: code,
		Msg:  msg,
	}
}

func NewSuccessCommonResp(msg string, data interface{}) *CommonResp {
	return &CommonResp{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
}
