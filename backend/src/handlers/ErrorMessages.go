package handlers

const (
	ErrFetchSwiftCodes    = "Failed to fetch SWIFT codes "
	ErrNoSwiftCodeFound   = "No SWIFT code found "
	ErrFailedToDelete     = "Could not delete a record"
	ErrFailedToInsert     = "Error inserting to database "
	ErrInvalidSwiftLength = "Invalid SWIFT code length. It must be either 8 or 11 characters long."
	ErrInvalidISO2Length  = "Invalid ISO2 code length. It must be 2 characters long."
)
