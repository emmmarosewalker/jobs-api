package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"

	"gopkg.in/guregu/null.v4"
)

var db *sql.DB

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

func main() {
	// Capture connection properties
	cfg := mysql.Config{
		User:                 os.Getenv("DBUSER"),
		Passwd:               os.Getenv("DBPASS"),
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "jobs",
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	listings, err := queryListingsByCompany("Intygrate")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listings found: %v\n", listings)

	// Hard-code ID 2 here to test the query.
	listing, err := listingByID(1)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", listing)

	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/listings", getAllListings)
	router.POST("/listings", postListing)
	router.GET("/listings/:company", getListingsByCompany)

	router.Run("localhost:7000")
}

// listingsByCompany queries for albums that have the specified artist name.
func queryListingsByCompany(name string) ([]Listing, error) {
	var listings []Listing

	rows, err := db.Query("SELECT * FROM listings WHERE company = ?", name)
	if err != nil {
		return nil, fmt.Errorf("listingsByCompany %q: %v", name, err)
	}
	defer rows.Close()
	for rows.Next() {
		var listing Listing
		if err := rows.Scan(&listing.ID, &listing.Company, &listing.Title, &listing.JobDescription, &listing.JobDescription, &listing.Category, &listing.StreetAddress, &listing.City, &listing.Country, &listing.BeginDate, &listing.Compensation); err != nil {
			return nil, fmt.Errorf("listingsByCompany %q: %v", name, err)
		}
		listings = append(listings, listing)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("listingsByCompany %q: %v", name, err)
	}
	return listings, nil
}

// getListingsByCompany responds to incoming request to query for listings by company
func getListingsByCompany(c *gin.Context) {
	company := c.Param("company")

	queried_companies, err := queryListingsByCompany(company)

	if err == nil {
		c.IndentedJSON(http.StatusOK, queried_companies)
		return
	}

	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "listings not found"})
}

// getAllListings queries for all listings
func queryAllListings() ([]Listing, error) {
	var listings []Listing

	rows, err := db.Query("SELECT * FROM listings")
	if err != nil {
		return nil, fmt.Errorf("queryAllListings %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var listing Listing
		if err := rows.Scan(&listing.ID, &listing.Company, &listing.Title, &listing.JobType, &listing.JobDescription, &listing.Category, &listing.StreetAddress, &listing.City, &listing.Country, &listing.BeginDate, &listing.Compensation); err != nil {
			return nil, fmt.Errorf("queryAllListings %v", err)
		}
		listings = append(listings, listing)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("queryAllListings %v", err)
	}
	return listings, nil
}

// getAllListings responds to incoming request for all listings
func getAllListings(c *gin.Context) {
	listings, err := queryAllListings()

	if err == nil {
		c.IndentedJSON(http.StatusOK, listings)
		return
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "get all listings: listings not found"})

}

// listingByID queries for the listing with the specified ID.
func listingByID(id int64) (Listing, error) {
	// An listing to hold data from the returned row.
	var listing Listing

	row := db.QueryRow("SELECT * FROM listings WHERE id = ?", id)
	if err := row.Scan(&listing.ID, &listing.Company, &listing.Title, &listing.JobDescription, &listing.JobDescription, &listing.Category, &listing.StreetAddress, &listing.City, &listing.Country, &listing.BeginDate, &listing.Compensation); err != nil {
		if err == sql.ErrNoRows {
			return listing, fmt.Errorf("listingByID %d: no such listing", id)
		}
		return listing, fmt.Errorf("listingByID %d: %v", id, err)
	}
	return listing, nil
}

// addListing adds the specified album to the database,
// returning the album ID of the new entry
func addListing(listing Listing) (int64, error) {
	result, err := db.Exec("INSERT INTO listings (company, title, job_type, job_description, category, street_address, city, country, begin_date, compensation) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", listing.Company, listing.Title, listing.JobType, listing.JobDescription, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("addListing: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addListing: %v", err)
	}
	return id, nil
}

func postListing(c *gin.Context) {
	var requestBody Listing

	err := c.BindJSON(&requestBody)

	if err != nil {
		fmt.Printf("postListing: %v", err)
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "postListing: check request body"})
		return
	}

	addListing(requestBody)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "postListing: could not insert into DB"})
		return
	}

	c.IndentedJSON(http.StatusOK, requestBody)
}
