package gosqlexec

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/gocraft/dbr/v2"
)

// Gosqlexec struct
type Gosqlexec struct {
	Sess        *dbr.Session
	CustomQuery string
	DropQuery   string
	AlterQuery  string
	Schemas     []string
}

// IGosqlexec interface
type IGosqlexec interface {
	MigrateSchemas()
	DropTablesIfExists()
	AlterTables()
	CustomQueryExecutor()
}

func lineByLineReader(path string) string {
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

func queryExecutor(
	sess *dbr.Session,
	path string,
) {
	query := lineByLineReader(path)
	fmt.Println(fmt.Sprintf("\n[PROC] Executing slq : %s", path))
	fmt.Print(query)
	result, err := sess.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	result.Close()
	fmt.Println("[OK] query executed successfully")
}

// CustomQueryExecutor custom alter tables sql
func (qexec Gosqlexec) CustomQueryExecutor() {
	queryExecutor(qexec.Sess, qexec.CustomQuery)
}

// AlterTables custom alter tables sql
func (qexec Gosqlexec) AlterTables() {
	queryExecutor(qexec.Sess, qexec.AlterQuery)
}

// DropTablesIfExists drop all schemas
func (qexec Gosqlexec) DropTablesIfExists() {
	queryExecutor(qexec.Sess, qexec.DropQuery)
}

// MigrateSchemas make qexec
func (qexec Gosqlexec) MigrateSchemas() {
	for _, path := range qexec.Schemas {
		queryExecutor(qexec.Sess, path)
	}
}
