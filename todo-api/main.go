package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var tasks = []Task{}
var nextId = 0

func getTasks(c *gin.Context) {
	c.JSON(http.StatusOK, tasks)
}
