package main

import (
	"example/data-access/db"
	"example/data-access/listings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database := db.ConnectDB()

	h := listings.NewBaseHandler(database)

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/listings", h.GetAllListings)
	router.POST("/listings", h.PostListing)
	router.GET("/listings/:company", h.GetListingsByCompany)

	router.Run("localhost:7000")
}
