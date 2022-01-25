package listings

import (
	"database/sql"
	"fmt"
)

// listingsByCompany queries for albums that have the specified artist name.
func (h *BaseHandler) queryListingsByCompany(name string) ([]Listing, error) {
	var listings []Listing

	rows, err := h.db.Query("SELECT * FROM listings WHERE company = ?", name)
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

// getAllListings queries for all listings
func (h *BaseHandler) queryAllListings() ([]Listing, error) {
	var listings []Listing

	rows, err := h.db.Query("SELECT * FROM listings")
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

// listingByID queries for the listing with the specified ID.
func (h *BaseHandler) listingByID(id int64) (Listing, error) {
	// An listing to hold data from the returned row.
	var listing Listing

	row := h.db.QueryRow("SELECT * FROM listings WHERE id = ?", id)
	if err := row.Scan(&listing.ID, &listing.Company, &listing.Title, &listing.JobDescription, &listing.JobDescription, &listing.Category, &listing.StreetAddress, &listing.City, &listing.Country, &listing.BeginDate, &listing.Compensation); err != nil {
		if err == sql.ErrNoRows {
			return listing, fmt.Errorf("listingByID %d: no such listing", id)
		}
		return listing, fmt.Errorf("listingByID %d: %v", id, err)
	}
	return listing, nil
}

// addListing adds the specified listing to the database
func (h *BaseHandler) addListing(listing Listing) (int64, error) {
	result, err := h.db.Exec("INSERT INTO listings (company, title, job_type, job_description, category, street_address, city, country, begin_date, compensation) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", listing.Company, listing.Title, listing.JobType, listing.JobDescription, nil, nil, nil, nil, nil, nil)
	if err != nil {
		return 0, fmt.Errorf("addListing: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addListing: %v", err)
	}
	return id, nil
}
