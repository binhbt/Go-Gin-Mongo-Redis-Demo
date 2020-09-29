package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func handleGetTasks(c *gin.Context) {
	var loadedTasks, err = GetAllTasks()
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"block_items": loadedTasks})
}

func handleGetTask(c *gin.Context) {
	var task BlockItem
	if err := c.BindUri(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	var loadedTask, err = GetTaskByID(task.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": loadedTask.ID, "user_id": loadedTask.UserID, "token": loadedTask.Token, "expired_time": loadedTask.ExpiredTime})
}

func handleCreateTask(c *gin.Context) {
	var task BlockItem
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	id, err := Create(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func handleUpdateTask(c *gin.Context) {
	var task BlockItem
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	savedTask, err := Update(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"block_item": savedTask})
}
func handleAddtoBlacklist(c *gin.Context) {
	var task BlockItem
	if err := c.ShouldBindJSON(&task); err != nil {
		log.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"msg": err})
		return
	}
	redisClient := getRedisClient()
	var key1 = ""
	if task.BlockType != "user" {
		key1 = fmt.Sprintf("block_token_%s", task.Token)
	} else {
		key1 = fmt.Sprintf("block_user_%d", task.UserID)
	}
	log.Printf("Key1: %s", key1)
	value1 := task
	err := redisClient.setKey(key1, value1, time.Duration(task.ExpiredTime))
	if err != nil {
		log.Printf("Error: %v", err.Error())
	}

	id, err := Create(&task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"msg": err})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}
func handleCheckBlacklist(c *gin.Context) {
	key := c.Param("key")
	blocktype := c.Param("blocktype")
	log.Printf("blocktype: %s", blocktype)
	var key1 = ""
	if blocktype != "user" {
		key1 = fmt.Sprintf("block_token_%s", key)
	} else {
		key1 = fmt.Sprintf("block_user_%s", key)
	}
	log.Printf("Key: %s", key1)
	redisClient := getRedisClient()
	value2 := &BlockItem{}
	var err = redisClient.getKey(key1, value2)
	if err != nil {
		var loadedTask, err = GetBlockByKey(key, blocktype)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"msg": err})
			return
		}
		log.Printf("loadedTask: %v", loadedTask)
		c.JSON(http.StatusOK, gin.H{"id": loadedTask.ID, "user_id": loadedTask.UserID, "token": loadedTask.Token, "expired_time": loadedTask.ExpiredTime})
	} else {
		// c.JSON(http.StatusNotFound, gin.H{"msg": "ok"})
		c.JSON(http.StatusOK, gin.H{"id": value2.ID, "user_id": value2.UserID, "token": value2.Token, "expired_time": value2.ExpiredTime})
	}

}

type valueEx struct {
	Name  string
	Email string
}

func main() {
	redisClient := getRedisClient()
	redisClient1 := getRedisClient()
	log.Printf("Expired: %d", time.Minute*1)
	key1 := "sampleKey"
	value1 := &valueEx{Name: "someName", Email: "someemail@abc.com"}
	err := redisClient.setKey(key1, value1, time.Minute*1)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	value2 := &valueEx{}
	err = redisClient1.getKey(key1, value2)
	if err != nil {
		log.Fatalf("Error: %v", err.Error())
	}

	log.Printf("Name: %s", value2.Name)
	log.Printf("Email: %s", value2.Email)

	r := gin.Default()
	r.GET("/blocks/:id", handleGetTask)
	r.GET("/blocks", handleGetTasks)
	r.PUT("/blocks", handleCreateTask)
	r.POST("/blocks", handleUpdateTask)

	r.POST("/blacklist", handleAddtoBlacklist)
	r.GET("/blacklist/check/:key/:blocktype", handleCheckBlacklist)
	r.Run() // listen and serve on 0.0.0.0:8080
}
