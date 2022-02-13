package listings

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type Listing struct {
	gorm.Model
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
	Compensation   null.String `json:"compensation,omitempty"`
}
type BaseHandler struct {
	db *gorm.DB
}

// NewBaseHandler returns a new BaseHandler
func NewBaseHandler(db *gorm.DB) *BaseHandler {
	db.AutoMigrate(&Listing{})
	return &BaseHandler{
		db: db,
	}
}

// GetListingsByCompany responds to incoming request to query for listings by company
func (h *BaseHandler) GetListingsByCompany(c *gin.Context) {
	id := c.Param("company")
	var listings []Listing

	err := h.queryListingsByCompany(id, &listings)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, listings)
}

// GetAllListings responds to incoming request for all listings
func (h *BaseHandler) GetAllListings(c *gin.Context) {
	var listings []Listing
	err := h.queryAllListings(&listings)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, listings)
}

func (h *BaseHandler) PostListing(c *gin.Context) {
	var listing Listing
	bindErr := c.BindJSON(&listing)

	if bindErr != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": bindErr})
		return
	}

	err := h.addListing(&listing)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, listing)
	return
}

// GetListingById responds to incoming request to query for a listing by id
func (h *BaseHandler) GetListingById(c *gin.Context) {
	id := c.Param("id")
	var listing Listing

	err := h.listingByID(&listing, id)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.IndentedJSON(http.StatusOK, listing)
}
