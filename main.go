package main


import (
	"net/http"

    "github.com/gin-gonic/gin"
	"github.com/pion/webrtc/v4"
)

var sessions = []webrtc.SessionDescription{};

func getSessions(c *gin.Context){
	// return with list of all sessions 
	c.IndentedJSON(http.StatusOK, sessions)
}

func postSession(c *gin.Context){
	var newSession webrtc.SessionDescription

	// bind json to newSession
	if err := c.BindJSON(&newSession); err != nil {
		return
	}

	// add new session to the slice
	sessions = append(sessions, newSession)
	c.IndentedJSON(http.StatusCreated, newSession)
}

func main(){
	router := gin.Default()
	router.GET("/sessions", getSessions)
	router.POST("/sessions", postSession)

	router.Run("localhost:8080")
}