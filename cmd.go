package gosqlexec

import (
	"fmt"
	"log"
	"os"

	"github.com/gocraft/dbr/v2"
	"github.com/joho/godotenv"
	"github.com/manifoldco/promptui"
	"github.com/urfave/cli"
	"golang.org/x/crypto/bcrypt"
)

var message = "Are you sure want to execute the query? please check, this process cant't be canceled!"
var enterPasswordMessage = "Enter password"

var gosqlexecFlags = []cli.Flag{
	cli.StringFlag{
		Name:  "loadenv",
		Value: ".env or C:/your/path/to/env/file",
		Usage: "Load env file with path or file name",
	},
}

func verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func password(fileName, question string) bool {
	secret := os.Getenv("SECRET")
	if len(secret) == 0 {
		return true
	}
	prompt := promptui.Prompt{
		Label: question,
	}
	result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	err = godotenv.Load(fileName)
	if err != nil {
		log.Fatalf("Fail load env %v\n", err)
	}
	err = verifyPassword("$2a$10$"+os.Getenv("SECRET"), result)
	if err != nil {
		log.Fatalf("Incorrect password fuck off!!!%v\n", err)
		return false
	}
	return true
}

func confirmation(question string) bool {
	prompt := promptui.Select{
		Label: question,
		Items: []string{"Yes", "No"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}
	return result == "Yes"
}

func connectToDB(fileName string) (*dbr.Connection, *dbr.Session, error) {
	err := godotenv.Load(fileName)
	if err != nil {
		return nil, nil, err
	}
	conn, err := dbr.Open("mysql", fmt.Sprintf(
		"%s:%s@(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	), nil)
	if err != nil {
		return nil, nil, err
	}
	sess := conn.NewSession(nil)
	sess.Begin()
	return conn, sess, nil
}

// MigrateCommand function
func MigrateCommand(gosqlexec GoSQLExec) cli.Command {
	return cli.Command{
		Name:  "migrate",
		Usage: fmt.Sprintf("Migrate schemas defined %s", gosqlexec.Schemas),
		Flags: gosqlexecFlags,
		Action: func(c *cli.Context) error {
			env := c.String("loadenv")
			_, sess, err := connectToDB(env)
			if err != nil {
				return err
			}
			gosqlexec.Sess = sess
			if confirmation(message) {
				if password(env, enterPasswordMessage) {
					err := gosqlexec.MigrateSchemas()
					if err != nil {
						log.Fatal(err)
					}
				}
				os.Exit(1)
			}
			os.Exit(1)
			return nil
		},
	}
}

// DropTablesCommand function
func DropTablesCommand(gosqlexec GoSQLExec) cli.Command {
	return cli.Command{
		Name:  "drop-tables",
		Usage: fmt.Sprintf("Drop any tables defined inside %s", gosqlexec.DropQuery),
		Flags: gosqlexecFlags,
		Action: func(c *cli.Context) error {
			env := c.String("loadenv")
			_, sess, err := connectToDB(env)
			if err != nil {
				return err
			}
			gosqlexec.Sess = sess
			if confirmation(message) {
				if password(env, enterPasswordMessage) {
					err := gosqlexec.DropTablesIfExists()
					if err != nil {
						log.Fatal(err)
					}
				}
				os.Exit(1)
			}
			os.Exit(1)
			return nil
		},
	}
}

// AlterTablesCommand function
func AlterTablesCommand(gosqlexec GoSQLExec) cli.Command {
	return cli.Command{
		Name:  "alter-tables",
		Usage: fmt.Sprintf("Alter any tables defined inside %s", gosqlexec.AlterQuery),
		Flags: gosqlexecFlags,
		Action: func(c *cli.Context) error {
			env := c.String("loadenv")
			_, sess, err := connectToDB(env)
			if err != nil {
				return err
			}
			gosqlexec.Sess = sess
			if confirmation(message) {
				if password(env, enterPasswordMessage) {
					err := gosqlexec.AlterTables()
					if err != nil {
						log.Fatal(err)
					}
				}
				os.Exit(1)
			}
			os.Exit(1)
			return nil
		},
	}
}

// CustomQueryExecCommand function
func CustomQueryExecCommand(gosqlexec GoSQLExec) cli.Command {
	return cli.Command{
		Name:  "custom-query-exec",
		Usage: fmt.Sprintf("Execute any queries defined inside %s", gosqlexec.CustomQuery),
		Flags: gosqlexecFlags,
		Action: func(c *cli.Context) error {
			env := c.String("loadenv")
			_, sess, err := connectToDB(env)
			if err != nil {
				return err
			}
			gosqlexec.Sess = sess
			if confirmation(message) {
				if password(env, enterPasswordMessage) {
					err := gosqlexec.CustomQueryExecutor()
					if err != nil {
						log.Fatal(err)
					}
				}
				os.Exit(1)
			}
			os.Exit(1)
			return nil
		},
	}
}
