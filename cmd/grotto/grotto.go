package grotto

import (
	"flag"

	"github.com/eaneto/grotto/pkg/connection"
	"github.com/eaneto/grotto/pkg/processor"
)

// Run runs the main grotto program. All options are captured with te
// command line and are used to create a new migration processor.
func Run() {
	user := flag.String("user", "", "Database user's name")
	password := flag.String("password", "", "Database user's password")
	database := flag.String("database", "", "Name of the database")
	address := flag.String("addresss", "localhost", "Database server address")
	port := flag.String("port", "5432", "Database server port")
	migrationDirectory := flag.String("dir", "", "The migration directory containing the scripts to be executed")

	flag.Parse()

	migrationProcessor := processor.New(connection.DatabaseInformation{
		User:     *user,
		Password: *password,
		Database: *database,
		Address:  *address,
		Port:     *port,
	}, *migrationDirectory)
	migrationProcessor.ProcessMigration()
}
