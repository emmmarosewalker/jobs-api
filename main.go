package main

import (
	"github.com/emmmarosewalker/jobs-api/db"
	"github.com/emmmarosewalker/jobs-api/listings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	database := db.InitDb()

	h := listings.NewBaseHandler(database)

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/listings", h.GetAllListings)
	router.POST("/listings", h.PostListing)
	router.GET("/listings/:company", h.GetListingsByCompany)

	router.Run("localhost:7000")
}
