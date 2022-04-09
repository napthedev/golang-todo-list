package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type Todo struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
}

type NewTodo struct {
	Value string `json:"value"`
}

func main() {
	data, err := os.ReadFile("./data.json")
	if err != nil {
		panic(err)
	}

	var todos []Todo
	json.Unmarshal(data, &todos)

	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, gin.H{
			"message": "Hello world",
		})
	})

	router.GET("/todos", func(ctx *gin.Context) {
		ctx.IndentedJSON(http.StatusOK, todos)
	})

	router.GET("/todos/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		for _, item := range todos {
			if fmt.Sprintf("%v", item.ID) == id {
				ctx.IndentedJSON(http.StatusOK, item)
				return
			}
		}
		ctx.IndentedJSON(http.StatusNotFound, gin.H{
			"notFound": true,
		})
	})

	router.POST("/todos", func(ctx *gin.Context) {
		var newTodo NewTodo

		if err := ctx.BindJSON(&newTodo); err != nil {
			ctx.IndentedJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		todos = append(todos, Todo{

			ID:    todos[len(todos)-1].ID + 1,
			Value: newTodo.Value,
		})

		marshaled, err := json.Marshal(todos)
		var indented bytes.Buffer
		json.Indent(&indented, marshaled, "", "  ")
		if err != nil {
			ctx.IndentedJSON(http.StatusInternalServerError, gin.H{
				"message": "Something went wrong",
			})
		}
		os.WriteFile("./data.json", indented.Bytes(), fs.ModeExclusive)
		ctx.IndentedJSON(http.StatusOK, todos[len(todos)-1])
	})

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "5000"
	}
	fmt.Printf("Server is listening on port %v\n", port)
	router.Run(fmt.Sprintf(":%v", port))
}
