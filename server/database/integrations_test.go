package database

// import (
// 	"context"
// 	"fmt"
// 	"sync"
// 	"testing"
// 	"time"

// 	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
// 	"github.com/staszkiet/DictionaryGolang/server/graph/model"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// 	"github.com/testcontainers/testcontainers-go"
// 	"github.com/testcontainers/testcontainers-go/wait"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// type DictionaryTestSuite struct {
// 	suite.Suite
// 	DB   *gorm.DB
// 	repo dictionaryRepository
// 	svc  DictionaryService
// }

// func (s *DictionaryTestSuite) SetupSuite() {
// 	ctx := context.Background()

// 	dbName := fmt.Sprintf("dict_test_%d", time.Now().UnixNano())

// 	req := testcontainers.ContainerRequest{
// 		Image:        "postgres:14",
// 		ExposedPorts: []string{"5432/tcp"},
// 		Env: map[string]string{
// 			"POSTGRES_PASSWORD": "password",
// 			"POSTGRES_USER":     "postgres",
// 			"POSTGRES_DB":       dbName,
// 		},
// 		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
// 	}

// 	var err error
// 	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
// 		ContainerRequest: req,
// 		Started:          true,
// 	})
// 	if err != nil {
// 		s.T().Fatalf("Failed to start postgres container: %v", err)
// 	}

// 	host, err := postgresContainer.Host(ctx)
// 	if err != nil {
// 		s.T().Fatalf("Failed to get container host: %v", err)
// 	}

// 	port, err := postgresContainer.MappedPort(ctx, "5432")
// 	if err != nil {
// 		s.T().Fatalf("Failed to get container port: %v", err)
// 	}

// 	dsn := fmt.Sprintf("host=%s user=postgres password=password dbname=%s port=%s sslmode=disable TimeZone=UTC",
// 		host, dbName, port.Port())

// 	s.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		s.T().Fatalf("Failed to connect to test database: %v", err)
// 	}

// 	err = s.DB.AutoMigrate(&dbmodels.Word{}, &dbmodels.Translation{}, &dbmodels.Sentence{})
// 	if err != nil {
// 		s.T().Fatalf("Failed to migrate schema: %v", err)
// 	}

// 	s.repo = dictionaryRepository{s.DB}
// 	s.svc = DictionaryService{&s.repo}
// }

// func (s *DictionaryTestSuite) SetupTest() {
// 	s.DB.Exec("TRUNCATE words CASCADE")
// }

// func TestDictionarySuite(t *testing.T) {
// 	suite.Run(t, new(DictionaryTestSuite))
// }

// func (s *DictionaryTestSuite) TestCreateWord() {

// 	baseWord := "rower"
// 	translation := model.NewTranslation{English: "bike", Sentences: []string{"I like my bike"}}

// 	_, err := s.svc.CreateWord(context.Background(), baseWord, translation)

// 	assert.NoError(s.T(), err)

// 	word, err := s.svc.SelectWord(context.Background(), baseWord)

// 	assert.NoError(s.T(), err)
// 	assert.Equal(s.T(), baseWord, word.Polish)
// }

// func (s *DictionaryTestSuite) TestCreateWordParallel() {

// 	baseWord := "równoległy"
// 	translation1 := model.NewTranslation{English: "parallel", Sentences: []string{"These lines are parallel."}}
// 	translation2 := model.NewTranslation{English: "concurrent", Sentences: []string{"We received concurrent requests."}}

// 	var wg sync.WaitGroup
// 	wg.Add(2)

// 	retChan := make(chan error, 2)

// 	go func() {
// 		defer wg.Done()
// 		_, err := s.svc.CreateWord(context.Background(), baseWord, translation1)
// 		retChan <- err
// 	}()

// 	go func() {
// 		defer wg.Done()
// 		_, err := s.svc.CreateWord(context.Background(), baseWord, translation2)
// 		retChan <- err
// 	}()

// 	wg.Wait()
// 	close(retChan)

// 	var errors []error
// 	for err := range retChan {
// 		errors = append(errors, err)
// 	}

// 	assert.Equal(s.T(), 2, len(errors))
// 	if assert.Nil(s.T(), errors[0], "First result should be nil") {
// 		assert.NotNil(s.T(), errors[1], "Second result should not be nil")
// 	} else {
// 		assert.Nil(s.T(), errors[1], "Second result should be nil")
// 		assert.NotNil(s.T(), errors[0], "First result should not be nil")
// 	}

// 	var count int64
// 	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)
// 	assert.Equal(s.T(), int64(1), count)
// }
