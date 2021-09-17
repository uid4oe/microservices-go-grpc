package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uid4oe/microservices-go-grpc/bff/client"
)

var (
	timeout       = time.Second
	user_client   client.UserClient
	advice_client client.AdviceClient
)

func UserRegister(router *gin.RouterGroup) {
	router.GET("/:id", GetUserDetails)
	router.PUT("/:id", UpdateUser)
	router.GET("/", GetUsers)
	router.POST("/", CreateUser)
}

func AdviceRegister(router *gin.RouterGroup) {
	router.PUT("/", UpdateAdvice)
}

func GetUsers(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	data, err := user_client.GetUsers(&ctx)
	response(c, data, err)
}

func GetUserDetails(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	param, err := getParam(c, "id")
	if err != nil {
		response(c, nil, err)
		return
	}

	data, err := user_client.GetUserDetails(param, &ctx)
	if err != nil {
		response(c, nil, err)
		return
	}
	advice_data, err := advice_client.GetAdvice(param, &ctx)
	if err != nil {
		response(c, nil, err)
		return
	}
	data.Advice = *advice_data
	response(c, data, err)
}

func UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var user client.UserWithDetails
	err := c.BindJSON(&user)
	if err != nil {
		response(c, nil, err)
		return
	}

	_, err = user_client.CreateUpdateUser(user, &ctx, c.Request.Method)
	response(c, nil, err)
}

func CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var user client.UserWithDetails
	err := c.BindJSON(&user)
	if err != nil {
		response(c, nil, err)
		return
	}

	var user_id string
	user_id, err = user_client.CreateUpdateUser(user, &ctx, c.Request.Method)
	if err != nil {
		response(c, nil, err)
		return
	}

	err = advice_client.CreateUpdateAdvice(client.Advice{UserId: user_id, Advice: user.Advice.Advice}, &ctx, c.Request.Method)
	if err != nil {
		response(c, nil, err)
		return
	}
	response(c, nil, err)
}

func UpdateAdvice(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c, timeout)
	defer cancel()

	var req client.Advice
	err := c.BindJSON(&req)
	if err != nil {
		response(c, nil, err)
		return
	}

	err = advice_client.CreateUpdateAdvice(req, &ctx, c.Request.Method)
	response(c, nil, err)
}

func response(c *gin.Context, data interface{}, err error) {
	statusCode := http.StatusOK
	var errorMessage string
	if err != nil {
		log.Println("Server Error Occured:", err)
		errorMessage = strings.Title(err.Error())
		statusCode = http.StatusInternalServerError
	}
	c.JSON(statusCode, gin.H{"data": data, "error": errorMessage})
}

func getParam(c *gin.Context, param string) (string, error) {
	p := c.Param(param)
	if len(p) == 0 {
		return "", errors.New("invalid parameter: " + p)
	}
	return p, nil
}

func main() {
	log.Println("Bff Service")

	r := gin.Default()
	r.Use(cors.Default())

	api := r.Group("/api")
	UserRegister(api.Group("/users"))

	AdviceRegister(api.Group("/advices"))

	r.Run()
}
