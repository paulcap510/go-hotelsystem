package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	// "errors"
)


type reservation struct {
	ID			string	`json:"id"`
	StartDate	string	`json:"startDate"`
	EndDate		string	`json:"endDate"`
	Name		string	`json:"name"`
	Guests   	int	    `json:"guests"`
}

//* stored rooms
type room struct {
	ID 				string 			`json:"id"`
	Type			string			`json:"type"`
	Reservations    []reservation    `json:"reservations"`
}


//* slice struct to imitate database
var rooms = []room {
	{ID: "101", Type: "Single"},
	{ID: "102", Type: "Single"},
	{ID: "103", Type: "Single"},
	{ID: "104", Type: "Single"},
	{ID: "105", Type: "Double"},
	{ID: "106", Type: "Double"},
	{ID: "107", Type: "King"},
	{ID: "108", Type: "King"},
	{ID: "109", Type: "Suite"},
	{ID: "110", Type: "Suite"},
}



func getReservations(c *gin.Context) {
    var allReservations []reservation 

	//? this loops through the rooms and appends to the allReservations slice each of the room.Reservations
    for _, room := range rooms {
        allReservations = append(allReservations, room.Reservations...)
    }
    c.IndentedJSON(http.StatusOK, allReservations)
}

//* function to add a reservation to a room
func addReservation(c *gin.Context) {
    var newReservation reservation
	roomID := c.Param("roomId")

    if err := c.BindJSON(&newReservation); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    for i, room := range rooms {
        if room.ID == roomID {
            for _, r := range room.Reservations {
                if (newReservation.StartDate <= r.EndDate) && (newReservation.EndDate >= r.StartDate) {
                    c.JSON(http.StatusBadRequest, gin.H{"error": "Date conflict"})
                    return
                }
            }
            rooms[i].Reservations = append(rooms[i].Reservations, newReservation)
            c.JSON(http.StatusOK, gin.H{"message": "Reservation added"})
            return
        }
    }
    c.JSON(http.StatusNotFound, gin.H{"message": "Room not found"})
}


//* function to return JSON version of rooms
func getRooms(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, rooms)
}


func DeleteReservation(c *gin.Context) {
    roomID := c.Param("roomId") 
    resID := c.Param("reservationId") 

    foundRoom := false
    foundRes := false
    for i, room := range rooms {
        if room.ID == roomID {
            foundRoom = true
            for j, reservation := range room.Reservations {
                if reservation.ID == resID {
                    rooms[i].Reservations = append(room.Reservations[:j], room.Reservations[j+1:]...)
                    foundRes = true
                    break
                }
            }
            if foundRes {
                break 
            }
        }
    }

    if foundRoom && foundRes {
        c.IndentedJSON(http.StatusOK, gin.H{"message": "Reservation deleted"})
    } else if !foundRoom {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Room not found"})
    } else {
        c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Reservation not found"})
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


func editReservation(c *gin.Context) {
    roomId := c.Param("roomId")          
    reservationId := c.Param("reservationId")  

    var updatedReservation reservation
    if err := c.BindJSON(&updatedReservation); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    foundRoom := false
    foundReservation := false
    for i, room := range rooms {
        if room.ID == roomId {                
            foundRoom = true
            for j, res := range room.Reservations {
                if res.ID == reservationId {  
                    rooms[i].Reservations[j] = updatedReservation 
                    foundReservation = true
                    break
                }
            }
            if foundReservation {
                break 
            }
        }
    }

    if !foundRoom {
        c.JSON(http.StatusNotFound, gin.H{"message": "Room not found"})
    } else if !foundReservation {
        c.JSON(http.StatusNotFound, gin.H{"message": "Reservation not found"})
    } else {
        c.JSON(http.StatusOK, gin.H{"message": "Reservation updated"})
    }
}


func main() {
    router := gin.Default()
    router.GET("/rooms", getRooms)
    router.POST("/rooms", createRoom)
    router.DELETE("/rooms/:roomId", DeleteRoom)
    router.POST("/rooms/:roomId/reservations", addReservation)
    router.GET("/reservations", getReservations)
	router.PATCH("/rooms/:roomId/reservations/:reservationId", editReservation)
    router.DELETE("/rooms/:roomId/reservations/:reservationId", DeleteReservation)
    router.Run("localhost:8000")
}


