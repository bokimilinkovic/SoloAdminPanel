package dto

// SearchByName represents dto object for name search
type SearchByName struct {
	Firstname string `json:"firstname"`
}

// SearchByEmail represents dto object for email search
type SearchByEmail struct {
	Email string `json:"email"`
}

// SearchByGUID represents dto object for guid search
type SearchByGUID struct {
	Guid string `json:"guid"`
}

// Search is used for advanced search, for not specific parameter.
type Search struct {
	Query string
}

// SearchByUserID is used for searching loans by user_id (which is user's guid value)
type SearchByUserID struct {
	UserId string `json:"user_id"`
}
