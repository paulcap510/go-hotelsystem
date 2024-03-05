package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "errors"
)

//* stored rooms
type room struct {
	ID 			string 		`json:"id"`
	Type		string		`json:"type"`
	Vacant		bool		`json:"vacant"`
}


//* slice struct to imitate database
var rooms = []room {
	{ID: "101", Type: "Single", Vacant: true},
	{ID: "102", Type: "Single", Vacant: true},
	{ID: "103", Type: "Single", Vacant: true},
	{ID: "104", Type: "Single", Vacant: true},
	{ID: "105", Type: "Double", Vacant: true},
	{ID: "106", Type: "Double", Vacant: true},
	{ID: "107", Type: "King", Vacant: true},
	{ID: "108", Type: "King", Vacant: true},
	{ID: "109", Type: "Suite", Vacant: true},
	{ID: "110", Type: "Suite", Vacant: true},
}


//* function to return JSON version of rooms
func getRooms(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, rooms)
}


//! func to change room to vacant
func changeOccupancy(c *gin.Context) {
    roomID := c.Param("id")     
    // initializes BOOL named 'found' to 'false'
	// used to check whether an ID was found
    found := false

	// iterates over rooms and toggles vacant or not vacant 
    for i, room := range rooms {
        if room.ID == roomID {
            rooms[i].Vacant = !rooms[i].Vacant
            found = true
            c.IndentedJSON(http.StatusOK, rooms[i])
            break
        }
    }
    if !found {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Room not found"})
    }
}

// ! delete room 
func DeleteRoom(c *gin.Context) {
    roomID := c.Param("id") 

    found := false

    for i, room := range rooms {
        if room.ID == roomID {
            rooms = append(rooms[:i], rooms[i+1:]...)
            found = true
            break 
        }
    }

	if found {
        c.IndentedJSON(http.StatusOK, gin.H{"message": "Room deleted"})
    } else {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Room not found"})
    }
}



func createRoom (c *gin.Context) {
	var newRoom room

	if err := c.BindJSON(&newRoom); err != nil {
		return
	}

	rooms = append(rooms, newRoom)
	c.IndentedJSON(http.StatusCreated, newRoom)
}

func main() {
	router := gin.Default()
	router.GET("/rooms", getRooms)
	router.POST("/rooms", createRoom)
	router.PUT("/rooms/:id", changeOccupancy)
	router.DELETE("/rooms/:id", DeleteRoom)
	router.Run("localhost:8000")
}