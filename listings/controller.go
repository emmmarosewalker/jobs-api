package listings

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"
)

type BaseHandler struct {
	db *sql.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *sql.DB) *BaseHandler {
	return &BaseHandler{
		db: db,
	}
}

type Listing struct {
	ID             int64       `json:"id"`
	Company        string      `json:"company"`
	Title          string      `json:"title"`
	JobType        string      `json:"jobType"` // e.g. full time, part time, casual
	JobDescription string      `json:"jobDescription"`
	Category       null.String `json:"category,omitempty"`
	StreetAddress  null.String `json:"streetAddress,omitempty"`
	City           null.String `json:"city,omitempty"`
	Country        null.String `json:"country,omitempty"`
	BeginDate      null.Time   `json:"beginDate,omitempty"`
	Compensation   null.Float  `json:"compensation,omitempty"`
}

// GetListingsByCompany responds to incoming request to query for listings by company
func (h *BaseHandler) GetListingsByCompany(c *gin.Context) {
	company := c.Param("company")

	queried_companies, err := h.queryListingsByCompany(company)

	if err == nil {
		c.IndentedJSON(http.StatusOK, queried_companies)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "listings not found"})
}

// GetAllListings responds to incoming request for all listings
func (h *BaseHandler) GetAllListings(c *gin.Context) {
	listings, err := h.queryAllListings()

	if err == nil {
		c.IndentedJSON(http.StatusOK, listings)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "get all listings: listings not found"})

}

func (h *BaseHandler) PostListing(c *gin.Context) {
	var requestBody Listing

	err := c.BindJSON(&requestBody)

	if err != nil {
		fmt.Printf("postListing: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "postListing: check request body"})
		return
	}

	h.addListing(requestBody)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "postListing: could not insert into DB"})
		return
	}

	c.IndentedJSON(http.StatusOK, requestBody)
}
