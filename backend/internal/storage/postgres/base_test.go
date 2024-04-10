package postgres

import (
	"log"
	"os"
	"testing"
)

func fillTestTables() error {
	return nil
}

func TestMain(m *testing.M) {
	testDB := SetupTestDatabase()
	defer testDB.TearDown()
	testDbInstance = testDB.DbInstance
	err := SeedTestData(testDbInstance)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
}
