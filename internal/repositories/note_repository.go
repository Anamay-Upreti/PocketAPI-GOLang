package repositories

import (
	"github.com/google/uuid"

	"pocketapi/internal/database"
	"pocketapi/internal/models"
)

type NoteRepository struct{}

func NewNoteRepository() *NoteRepository {
	return &NoteRepository{}
}

// Create note.
func (repo *NoteRepository) Create(note *models.Note) error {

	return database.DB.Create(note).Error
}

// Find note by ID.
func (repo *NoteRepository) FindByID(id string) (*models.Note, error) {

	var note models.Note

	err := database.DB.
		Where("id = ?", id).
		First(&note).Error

	if err != nil {
		return nil, err
	}

	return &note, nil
}

// Get all notes of a user.
func (repo *NoteRepository) FindAll(
	userID string,
	page int,
	limit int,
	search string,
) ([]models.Note, int64, error) {

	var notes []models.Note
	var total int64

	query := database.DB.Model(&models.Note{}).
		Where("user_id = ?", userID)

	if search != "" {
		query = query.Where(
			"title ILIKE ? OR content ILIKE ?",
			"%"+search+"%",
			"%"+search+"%",
		)
	}

	query.Count(&total)

	offset := (page - 1) * limit

	err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&notes).Error

	return notes, total, err
}

// Update note.
func (repo *NoteRepository) Update(note *models.Note) error {

	return database.DB.Save(note).Error
}

// Delete note.
func (repo *NoteRepository) Delete(id uuid.UUID) error {

	return database.DB.Delete(&models.Note{}, id).Error
}