package database

// SQLScript Represents a SQL script with the filename and content.
type SQLScript struct {
	// The actual SQL script content.
	Content string
	// The script filename with the .sql extension.
	Name string
}
