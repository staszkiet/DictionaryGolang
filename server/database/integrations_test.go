package database

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	dbmodels "github.com/staszkiet/DictionaryGolang/server/database/models"
	"github.com/staszkiet/DictionaryGolang/server/graph/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DictionaryTestSuite struct {
	suite.Suite
	DB   *gorm.DB
	repo dictionaryRepository
	svc  DictionaryService
}

func (s *DictionaryTestSuite) SetupSuite() {
	ctx := context.Background()

	dbName := fmt.Sprintf("dict_test_%d", time.Now().UnixNano())

	req := testcontainers.ContainerRequest{
		Image:        "postgres:14",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_USER":     "postgres",
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
	}

	var err error
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		s.T().Fatalf("Failed to start postgres container: %v", err)
	}

	host, err := postgresContainer.Host(ctx)
	if err != nil {
		s.T().Fatalf("Failed to get container host: %v", err)
	}

	port, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		s.T().Fatalf("Failed to get container port: %v", err)
	}

	dsn := fmt.Sprintf("host=%s user=postgres password=password dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, dbName, port.Port())

	s.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		s.T().Fatalf("Failed to connect to test database: %v", err)
	}

	err = s.DB.AutoMigrate(&dbmodels.Word{}, &dbmodels.Translation{}, &dbmodels.Sentence{})
	if err != nil {
		s.T().Fatalf("Failed to migrate schema: %v", err)
	}

	s.repo = dictionaryRepository{s.DB}
	s.svc = DictionaryService{&s.repo}
}

func (s *DictionaryTestSuite) SetupTest() {

	s.DB.Exec("TRUNCATE words CASCADE")
	s.DB.Exec("TRUNCATE translations CASCADE")
	s.DB.Exec("TRUNCATE sentences CASCADE")

}

func TestDictionarySuite(t *testing.T) {
	suite.Run(t, new(DictionaryTestSuite))
}

func (s *DictionaryTestSuite) TestCreateWord() {

	baseWord := "rower"
	translation := model.NewTranslation{English: "bike", Sentences: []string{"I like my bike"}}

	_, err := s.svc.CreateWordOrAddTranslationOrSentence(baseWord, translation)

	assert.NoError(s.T(), err)

	word, err := s.svc.SelectWord(baseWord)

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), baseWord, word.Polish)
}

func (s *DictionaryTestSuite) TestCreateWordParallel() {

	var count int64

	baseWord := "równoległy"
	translations := make([]model.NewTranslation, 0)
	for i := 0; i < 20; i++ {
		translations = append(translations, model.NewTranslation{English: strconv.Itoa(i), Sentences: []string{"a"}})
	}

	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)
	assert.Equal(s.T(), int64(0), count)

	var wg sync.WaitGroup
	wg.Add(20)

	retChan := make(chan error, 20)

	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			_, err := s.svc.CreateWordOrAddTranslationOrSentence(baseWord, translations[i])
			retChan <- err
		}()
	}

	wg.Wait()
	close(retChan)

	var errors []error
	for err := range retChan {
		errors = append(errors, err)
	}

	assert.Equal(s.T(), 20, len(errors))
	for i := 0; i < 20; i++ {
		assert.Nil(s.T(), errors[i])
	}

	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)

	assert.Equal(s.T(), int64(1), count)

	s.DB.Model(&dbmodels.Translation{}).Count(&count)

	assert.Equal(s.T(), int64(20), count)
}

func (s *DictionaryTestSuite) TestAddTranslationParallel() {

	var count int64

	baseWord := "równoległy"
	translation := model.NewTranslation{English: "parallel", Sentences: []string{"These lines are parallel."}}
	s.svc.CreateWordOrAddTranslationOrSentence(baseWord, translation)
	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)
	assert.Equal(s.T(), int64(1), count)

	translations := make([]model.NewTranslation, 0)
	for i := 0; i < 20; i++ {
		translations = append(translations, model.NewTranslation{English: "parallel", Sentences: []string{strconv.Itoa(i)}})
	}

	var wg sync.WaitGroup
	wg.Add(20)

	retChan := make(chan error, 20)

	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			_, err := s.svc.CreateWordOrAddTranslationOrSentence(baseWord, translations[i])
			retChan <- err
		}()
	}

	wg.Wait()
	close(retChan)

	var errors []error
	for err := range retChan {
		errors = append(errors, err)
	}

	assert.Equal(s.T(), 20, len(errors))
	for i := 0; i < 20; i++ {
		assert.Nil(s.T(), errors[i])
	}

	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)
	assert.Equal(s.T(), int64(1), count)
	s.DB.Model(&dbmodels.Translation{}).Count(&count)
	assert.Equal(s.T(), int64(1), count)
	s.DB.Model(&dbmodels.Sentence{}).Count(&count)
	assert.Equal(s.T(), int64(21), count)
}

func (s *DictionaryTestSuite) TestDeleteWordParallel() {

	var count int64

	baseWord := "równoległy"
	translation := model.NewTranslation{English: "parallel", Sentences: []string{"These lines are parallel."}}
	s.svc.CreateWordOrAddTranslationOrSentence(baseWord, translation)
	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)
	assert.Equal(s.T(), int64(1), count)

	var wg sync.WaitGroup
	wg.Add(2)

	retChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		_, err := s.svc.DeleteWord(baseWord)
		retChan <- err
	}()

	go func() {
		defer wg.Done()
		_, err := s.svc.DeleteWord(baseWord)
		retChan <- err
	}()

	wg.Wait()
	close(retChan)

	var errors []error
	for err := range retChan {
		errors = append(errors, err)
	}

	assert.Equal(s.T(), 2, len(errors))
	assert.Nil(s.T(), errors[0])
	assert.Nil(s.T(), errors[1])

	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)
	assert.Equal(s.T(), int64(0), count)
}

func (s *DictionaryTestSuite) TestDeleteTranslationParallel_WhenItIsTheOnlyTranslation_ShouldAlsoDeleteWord() {

	var count int64

	baseWord := "równoległy"
	englishWord := "parallel"
	translation := model.NewTranslation{English: englishWord, Sentences: []string{"These lines are parallel."}}
	s.svc.CreateWordOrAddTranslationOrSentence(baseWord, translation)
	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)
	assert.Equal(s.T(), int64(1), count)
	s.svc.DeleteTranslation(baseWord, englishWord)

	var wg sync.WaitGroup
	wg.Add(2)

	retChan := make(chan error, 2)

	go func() {
		defer wg.Done()
		_, err := s.svc.DeleteTranslation(baseWord, englishWord)
		retChan <- err
	}()

	go func() {
		defer wg.Done()
		_, err := s.svc.DeleteTranslation(baseWord, englishWord)
		retChan <- err
	}()

	wg.Wait()
	close(retChan)

	var errors []error
	for err := range retChan {
		errors = append(errors, err)
	}

	assert.Equal(s.T(), 2, len(errors))
	assert.Nil(s.T(), errors[0])
	assert.Nil(s.T(), errors[1])

	s.DB.Model(&dbmodels.Word{}).Where("polish = ?", baseWord).Count(&count)
	assert.Equal(s.T(), int64(0), count)
}
