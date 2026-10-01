package services

import (
	"errors"

	"github.com/google/uuid"

	"pocketapi/internal/models"
	"pocketapi/internal/repositories"
)

type NoteService struct {
	repository *repositories.NoteRepository
}

func NewNoteService() *NoteService {
	return &NoteService{
		repository: repositories.NewNoteRepository(),
	}
}

type CreateNoteRequest struct {
	Title   string
	Content string
}

type UpdateNoteRequest struct {
	Title   string
	Content string
}

// Create note.
func (service *NoteService) Create(
	userID string,
	req CreateNoteRequest,
) (*models.Note, error) {

	userUUID, err := uuid.Parse(userID)

	if err != nil {
		return nil, err
	}

	note := models.Note{
		ID:      uuid.New(),
		Title:   req.Title,
		Content: req.Content,
		UserID:  userUUID,
	}

	err = service.repository.Create(&note)

	if err != nil {
		return nil, err
	}

	return &note, nil
}

// Get single note.
func (service *NoteService) GetByID(
	noteID string,
	userID string,
) (*models.Note, error) {

	note, err := service.repository.FindByID(noteID)

	if err != nil {
		return nil, err
	}

	if note.UserID.String() != userID {
		return nil, errors.New("you cannot access this note")
	}

	return note, nil
}

// Get all notes.
func (service *NoteService) GetAll(
	userID string,
	page int,
	limit int,
	search string,
) ([]models.Note, int64, error) {

	return service.repository.FindAll(
		userID,
		page,
		limit,
		search,
	)
}

// Update note.
func (service *NoteService) Update(
	noteID string,
	userID string,
	req UpdateNoteRequest,
) (*models.Note, error) {

	note, err := service.GetByID(noteID, userID)

	if err != nil {
		return nil, err
	}

	note.Title = req.Title
	note.Content = req.Content

	err = service.repository.Update(note)

	if err != nil {
		return nil, err
	}

	return note, nil
}

// Delete note.
func (service *NoteService) Delete(
	noteID string,
	userID string,
) error {

	note, err := service.GetByID(noteID, userID)

	if err != nil {
		return err
	}

	return service.repository.Delete(note.ID)
}