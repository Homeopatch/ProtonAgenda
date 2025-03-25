package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Durations []time.Duration

// MarshalJSON Implements Marshaller for JSON
func (d Durations) MarshalJSON() ([]byte, error) {
	// Convert []time.Duration to []string
	durations := make([]string, len(d))
	for i, duration := range d {
		durations[i] = duration.String() // Convert to string (e.g., "30m", "1h")
	}
	return json.Marshal(durations) // Marshal as JSON array of strings
}

// UnmarshalJSON Implements Unmarshaler for JSON
func (d *Durations) UnmarshalJSON(data []byte) error {
	var durations []string
	if err := json.Unmarshal(data, &durations); err != nil {
		return err
	}

	// Convert []string back to []time.Duration
	parsedDurations := make([]time.Duration, len(durations))
	for i, duration := range durations {
		parsed, err := time.ParseDuration(duration)
		if err != nil {
			return err
		}
		parsedDurations[i] = parsed
	}
	*d = parsedDurations
	return nil
}

// Value implements the driver.Valuer interface for database serialization
func (d Durations) Value() (driver.Value, error) {
	return d.MarshalJSON() // Serialize Durations to JSON
}

// Scan implements the sql.Scanner interface for database deserialization
func (d *Durations) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to scan Durations: value is not []byte")
	}
	return d.UnmarshalJSON(bytes) // Deserialize JSON to Durations
}

type AgendaInvite struct {
	gorm.Model
	ResourceID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	UserID            uint
	Description       string
	ExpiresAt         time.Time
	NotBefore         time.Time
	NotAfter          time.Time
	PaddingBefore     time.Duration
	PaddingAfter      time.Duration
	SlotSizes         Durations          `gorm:"type:json"` // Store durations as JSON array
	AgendaSources     []AgendaSource     `gorm:"many2many:invite_sources;"`
	ProceduralAgendas []ProceduralAgenda `gorm:"many2many:invite_procedural_agendas;"`
}

func (user AgendaInvite) BeforeCreate(tx *gorm.DB) error {
	field := tx.Statement.Schema.LookUpField("SlotSizes")
	if field.DataType == "json" {
		// do something
	}
	return nil
}
