package models

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Setup PostgreSQL database for tests
func setupTestDB() (*gorm.DB, error) {
	// Use environment variables or default values
	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "localhost"
	}
	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}
	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}
	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "password"
	}
	dbname := os.Getenv("POSTGRES_DB")
	if dbname == "" {
		dbname = "test_db"
	}

	//err := dropTestDatabase(host, port, user, password, dbname)
	//if err != nil {
	//	return nil, err
	//}

	// Connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// Connect to PostgreSQL
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	db.Exec("DROP TABLE IF EXISTS test_db")

	// Auto-migrate all models
	err = db.AutoMigrate(&User{}, &AgendaInvite{}, &AgendaSource{}, &AgendaItem{}, &ProceduralAgenda{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Test User Creation
func TestUserCreation(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	user := User{
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	// Create User
	result := db.Create(&user)
	assert.NoError(t, result.Error)
	assert.NotZero(t, user.ID)

	// Check User exists in DB
	var fetchedUser User
	err = db.First(&fetchedUser, user.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, user.Email, fetchedUser.Email)
}

// Test AgendaInvite Creation and Associations
func TestAgendaInviteCreation(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	user := User{
		Email:        "testuser@example.com",
		PasswordHash: "hashedpassword",
	}
	db.Create(&user)

	invite := AgendaInvite{
		UserID:      user.ID,
		Description: "Team Meeting",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(2 * time.Hour),
	}

	// Create AgendaInvite
	result := db.Create(&invite)
	assert.NoError(t, result.Error)
	assert.NotZero(t, invite.ID)

	// Check AgendaInvite exists in DB
	var fetchedInvite AgendaInvite
	err = db.First(&fetchedInvite, "id = ?", invite.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, invite.Description, fetchedInvite.Description)
	assert.Equal(t, user.ID, fetchedInvite.UserID)
}

// Test AgendaSource and AgendaItem Relations
func TestAgendaSourceAndItemRelations(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	user := User{
		Email:        "testuser@example.com",
		PasswordHash: "hashedpassword",
	}
	db.Create(&user)

	source := AgendaSource{
		//Url:    url.URL{Scheme: "https", Host: "example.com", Path: "/agenda"},
		Url:    "https://foo.com/agenda",
		UserID: user.ID,
	}
	db.Create(&source)

	item := AgendaItem{
		StartTime:      time.Now(),
		EndTime:        time.Now().Add(1 * time.Hour),
		Description:    "Sample Agenda Item",
		AgendaSourceID: source.ID,
		UserID:         user.ID,
	}
	db.Create(&item)

	// Fetch and verify AgendaSource and AgendaItem relationship
	var fetchedSource AgendaSource
	db.Preload("AgendaItems").First(&fetchedSource, source.ID)
	assert.NoError(t, err)
	assert.Len(t, fetchedSource.AgendaItems, 1)
	assert.Equal(t, "Sample Agenda Item", fetchedSource.AgendaItems[0].Description)
}

// Test ProceduralAgenda Creation
func TestProceduralAgendaCreation(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	agenda := ProceduralAgenda{
		Descriptor:  "Routine",
		Description: "Daily Standup Agenda",
	}

	// Create ProceduralAgenda
	result := db.Create(&agenda)
	assert.NoError(t, result.Error)
	assert.NotZero(t, agenda.ID)

	// Check ProceduralAgenda exists in DB
	var fetchedAgenda ProceduralAgenda
	err = db.First(&fetchedAgenda, agenda.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, agenda.Descriptor, fetchedAgenda.Descriptor)
	assert.Equal(t, agenda.Description, fetchedAgenda.Description)
}

// Test AgendaInvite and ProceduralAgenda Many-to-Many Relationship
func TestAgendaInviteProceduralAgendaRelationship(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)

	// Create a user for the invite
	user := User{
		Email:        "testuser@example.com",
		PasswordHash: "hashedpassword",
	}
	result := db.Create(&user)
	assert.NoError(t, result.Error)
	assert.NotZero(t, user.ID)

	// Create an AgendaInvite associated with the user
	invite := AgendaInvite{
		UserID:      user.ID, // Associate with the created user
		Description: "Project Kickoff",
		ExpiresAt:   time.Now().Add(24 * time.Hour),
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(2 * time.Hour),
	}

	// Create a ProceduralAgenda
	agenda := ProceduralAgenda{
		Descriptor:  "Kickoff",
		Description: "Agenda for Project Kickoff",
	}

	// Create Invite and Agenda
	result = db.Create(&invite)
	assert.NoError(t, result.Error)
	result = db.Create(&agenda)
	assert.NoError(t, result.Error)

	// Associate ProceduralAgenda with AgendaInvite
	err = db.Model(&invite).Association("ProceduralAgendas").Append(&agenda)
	assert.NoError(t, err)

	// Verify relationship
	var fetchedInvite AgendaInvite
	err = db.Preload("ProceduralAgendas").First(&fetchedInvite, invite.ID).Error
	assert.NoError(t, err)
	assert.Len(t, fetchedInvite.ProceduralAgendas, 1)
	assert.Equal(t, agenda.Descriptor, fetchedInvite.ProceduralAgendas[0].Descriptor)
}
