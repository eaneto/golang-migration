package connection

// DatabaseInformation Basic information needed to stablish a database connection.
type DatabaseInformation struct {
	User     string
	Password string
	Address  string
	Port     string
	Database string
}
