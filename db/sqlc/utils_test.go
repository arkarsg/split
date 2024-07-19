package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"time"

	u "github.com/arkarsg/splitapp/utils"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // used by migrator
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type TestDatabase struct {
	DbInstance *sql.DB
	DbAddress  string
	container  testcontainers.Container
}

func createContainer(c context.Context) (testcontainers.Container, *sql.DB, string, error) {
	config := u.GetDevDbEnvs()
	var env = map[string]string{
		"POSTGRES_PASSWORD": config.DbPassword,
		"POSTGRES_USER":     config.DbUser,
		"POSTGRES_DB":       config.DbName,
	}
	var port = nat.Port(config.DbPort)
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:12-alpine",
			ExposedPorts: []string{port.Port()},
			Env:          env,
			WaitingFor:   wait.ForListeningPort(port),
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(c, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err.Error())
	}

	p, err := container.MappedPort(c, "5432")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container exposed port: %v", err.Error())
	}

	dbAddr := fmt.Sprintf(
		"%s://%s:%s@localhost:%s/%s?sslmode=disable",
		config.DbDriver,
		config.DbUser,
		config.DbPassword,
		p.Port(),
		config.DbName,
	)

	fmt.Println(dbAddr)
	db, err := sql.Open(config.DbDriver, dbAddr)
	if err != nil {
		return container, db, dbAddr, fmt.Errorf("failed to esablish connection to %v : %v", dbAddr, err.Error())
	}
	return container, db, dbAddr, nil
}

func migrateDb(dbAddr string) error {
	migrations, _ := filepath.Abs("../migrations")
	fmt.Println(migrations)
	m, err := migrate.New(
		fmt.Sprintf("file:%s", migrations),
		dbAddr,
	)
	if err != nil {
		return err
	}
	defer m.Close()

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	log.Println("Migrated successfully")
	return nil
}

func SetUpTestDatabase() *TestDatabase {
	c, cancel := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstace, dbAddr, err := createContainer(c)
	if err != nil {
		log.Fatal("Failed to set up test database", err)
	}

	err = migrateDb(dbAddr)
	if err != nil {
		log.Fatal("Failed to perform migration: ", err)
	}
	cancel()

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstace,
		DbAddress:  dbAddr,
	}
}

func (t *TestDatabase) TearDown() {
	t.DbInstance.Close()
	_ = t.container.Terminate(context.Background())
}

func createRandomUser(test_input CreateUserParams) (User, error) {
	user, err := testQueries.CreateUser(context.Background(), test_input)
	return user, err
}

func createRandomTransaction() Transaction {
	var user User
	var err error
	user, err = testQueries.GetUserById(
		context.Background(),
		1,
	)

	if err != nil {
		user, _ = testQueries.CreateUser(
			context.Background(),
			CreateUserParams{
				Username: u.RandomUser(),
				Email:    u.RandomEmail(),
			})
	}

	txnParams := CreateTransactionParams{
		Amount:   u.RandomAmount(),
		Currency: CurrencySGD,
		Title:    u.RandomString(10),
		PayerID:  user.ID,
	}
	txn, _ := testQueries.CreateTransaction(
		context.Background(),
		txnParams,
	)
	return txn
}

func createRandomDebt() Debt {
	txn := createRandomTransaction()
	debt, _ := testQueries.CreateDebt(
		context.Background(),
		txn.ID,
	)
	return debt
}

func createRandomDebtDebtor() DebtDebtor {
	debtor, _ := createRandomUser(CreateUserParams{
		Username: u.RandomUser(),
		Email:    u.RandomEmail(),
	})
	debt := createRandomDebt()

	createDebtDebtorsParams := CreateDebtDebtorsParams{
		DebtID:   debt.ID,
		DebtorID: debtor.ID,
		Amount:   u.RandomAmount(),
		Currency: CurrencySGD,
	}
	dd, _ := testQueries.CreateDebtDebtors(
		context.Background(),
		createDebtDebtorsParams,
	)

	return dd
}

func createRandomPayment() Payment {
	dd := createRandomDebtDebtor()
	p, _ := testQueries.CreatePayment(
		context.Background(),
		CreatePaymentParams{
			DebtID:   dd.DebtID,
			DebtorID: dd.DebtorID,
			Amount:   u.RandomAmount(),
			Currency: CurrencySGD,
		},
	)

	return p
}
