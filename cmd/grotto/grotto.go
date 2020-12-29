package grotto

import (
	"flag"

	"github.com/eaneto/grotto/pkg/connection"
	"github.com/eaneto/grotto/pkg/processor"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// DATABASE_URL Basic postgres connection string.  All options are
// replaced with command line arguments.
const DATABASE_URL = "postgres://%s:%s@%s:%s/%s"

func Run() {
	user := flag.String("user", "", "Database user's name")
	password := flag.String("password", "", "Database user's password")
	database := flag.String("database", "", "Name of the database")
	address := flag.String("addresss", "localhost", "Database server address")
	port := flag.String("port", "5432", "Database server port")
	migrationDirectory := flag.String("dir", "", "The migration directory containing the scripts to be executed")

	flag.Parse()

	migrationProcessor := processor.CreateProcessor(connection.DatabaseInformation{
		User:     *user,
		Password: *password,
		Database: *database,
		Address:  *address,
		Port:     *port,
	})
	migrationProcessor.ProcessMigration(*migrationDirectory)
}
