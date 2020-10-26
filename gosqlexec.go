package gosqlexec

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gocraft/dbr/v2"
)

// GoSQLExec struct
type GoSQLExec struct {
	Sess        *dbr.Session
	CustomQuery string
	DropQuery   string
	AlterQuery  string
	Schemas     []string
}

// IGoSQLExec interface
type IGoSQLExec interface {
	MigrateSchemas()
	DropTablesIfExists()
	AlterTables()
	CustomQueryExecutor()
}

// LineByLineReader for reading sql file syntax
func LineByLineReader(path string) string {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var codes []string
	for scanner.Scan() {
		codes = append(codes, scanner.Text())
	}

	var stringCode string
	for _, eachLine := range codes {
		stringCode += eachLine + "\n"
	}
	return stringCode
}

// QueryExecutor execute query from line by line func
func QueryExecutor(
	sess *dbr.Session,
	path string,
) error {
	query := LineByLineReader(path)
	if query == "" {
		return fmt.Errorf("No SQL to exec")
	}
	fmt.Println("[OK] Query Loaded")
	result, err := sess.Query(query)
	if err != nil {
		return err
	}
	result.Close()
	fmt.Println("[OK] Query Executed")
	return nil
}

// CustomQueryExecutor custom alter tables sql
func (qexec GoSQLExec) CustomQueryExecutor() error {
	return QueryExecutor(qexec.Sess, qexec.CustomQuery)
}

// AlterTables custom alter tables sql
func (qexec GoSQLExec) AlterTables() error {
	return QueryExecutor(qexec.Sess, qexec.AlterQuery)
}

// DropTablesIfExists drop all schemas
func (qexec GoSQLExec) DropTablesIfExists() error {
	return QueryExecutor(qexec.Sess, qexec.DropQuery)
}

// MigrateSchemas make qexec
func (qexec GoSQLExec) MigrateSchemas() error {
	for _, path := range qexec.Schemas {
		err := QueryExecutor(qexec.Sess, path)
		if err != nil {
			return err
		}
	}
	return nil
}
