package listings

// listingsByCompany queries for albums that have the specified artist name.
func (h *BaseHandler) queryListingsByCompany(name string, Listings *[]Listing) error {
	var listings []Listing
	err := h.db.Where("company = ?", name).Find(&listings).Error

	if err != nil {
		return err
	}

	return nil
}

// getAllListings queries for all listings
func (h *BaseHandler) queryAllListings(Listings *[]Listing) error {
	err := h.db.Find(Listings).Error
	if err != nil {
		return err
	}
	return nil
}

// listingByID queries for the listing with the specified ID.
func (h *BaseHandler) listingByID(Listing *Listing, id int64) error {
	// An listing to hold data from the returned row.
	err := h.db.Where("id = ?", id).First(Listing).Error
	if err != nil {
		return err
	}
	return nil
}

// addListing adds the specified listing to the database
func (h *BaseHandler) addListing(Listing *Listing) error {
	err := h.db.Create(Listing).Error
	if err != nil {
		return err
	}
	return nil
}
