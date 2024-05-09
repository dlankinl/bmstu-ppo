package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	DbName = "postgres"
	DbUser = "postgres"
	DbPass = "postgres"
)

type TestDatabase struct {
	DbInstance *pgxpool.Pool
	DbAddress  string
	container  testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test: ", err)
	}

	err = migrateDb(dbAddr)
	if err != nil {
		log.Fatal("failed to perform db migration: ", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	tdb.DbInstance.Close()
	_ = tdb.container.Terminate(context.Background())
}

func createContainer(ctx context.Context) (testcontainers.Container, *pgxpool.Pool, string, error) {
	var env = map[string]string{
		"POSTGRES_PASSWORD": DbPass,
		"POSTGRES_USER":     DbUser,
		"POSTGRES_DB":       DbName,
	}
	var port = "5432/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:14-alpine",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog("database system is ready to accept connections"),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("postgres container ready and running at port: ", p.Port())

	time.Sleep(time.Second)

	dbAddr := fmt.Sprintf("localhost:%s", p.Port())
	db, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName))
	if err != nil {
		return container, db, dbAddr, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, dbAddr, nil
}

func migrateDb(dbAddr string) error {
	currentDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("получение текущей директории: %w", err)
	}

	relativePathToMigrationFiles := "./migrations"
	pathToMigrationFiles := filepath.Join(currentDir, "..", "..", "..", relativePathToMigrationFiles)

	databaseURL := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", DbUser, DbPass, dbAddr, DbName)
	m, err := migrate.New(fmt.Sprintf("file:%s", pathToMigrationFiles), databaseURL)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func SeedTestData(db *pgxpool.Pool) error {
	testDataDir := "../../../sql/test_data/"
	files, err := os.ReadDir(testDataDir)
	if err != nil {
		return fmt.Errorf("ошибка при чтении директории с тестовыми данными: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(testDataDir, file.Name())
			err = executeTestDataScript(db, filePath)
			if err != nil {
				return fmt.Errorf("выполнение sql скрипта %s: %w", file.Name(), err)
			}
		}
	}

	return nil
}

func executeTestDataScript(db *pgxpool.Pool, filePath string) error {
	scriptContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("чтение sql скрипта: %w", err)
	}

	_, err = db.Exec(context.Background(), string(scriptContent))
	if err != nil {
		return fmt.Errorf("выполнение sql скрипта: %w", err)
	}

	return nil
}
