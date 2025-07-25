package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Person struct {
	User string `uri:"user" binding:"required"`
}

type Calc struct {
	x int `form:"x"`
	y int `form:"y"`
}

func main() {
	r := gin.Default()
	r.GET("/greet/:user", greetUserHandlers)
	r.GET("/calc", calcHandlers)

	r.Run()
}

func greetUserHandlers(c *gin.Context) {
	name := c.Param("user")
	c.String(http.StatusOK, "Hello %s", name)
}

func calcHandlers(c *gin.Context) {
	xStr := c.Query("x")
	yStr := c.Query("y")

	x, errX := strconv.Atoi(xStr)
	y, errY := strconv.Atoi(yStr)

	if errX != nil || errY != nil {
		c.JSON(400, gin.H{"error": "Неправильные параметры"})
		return
	}

	sum := x + y

	log.Println(x)
	c.JSON(200, gin.H{"sum": sum})
}
