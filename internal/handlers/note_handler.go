package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"pocketapi/internal/services"
	"pocketapi/internal/utils"
	"pocketapi/internal/validators"
)

type NoteHandler struct {
	service *services.NoteService
}

func NewNoteHandler() *NoteHandler {
	return &NoteHandler{
		service: services.NewNoteService(),
	}
}

type NoteBody struct {
	Title   string `json:"title" validate:"required,min=3,max=200"`
	Content string `json:"content" validate:"required"`
}

// POST /notes
func (handler *NoteHandler) Create(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	var body NoteBody

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "Invalid JSON body")
	}

	if err := validators.Validate.Struct(body); err != nil {
		return utils.Error(c, 400, err.Error())
	}

	note, err := handler.service.Create(
		userID,
		services.CreateNoteRequest{
			Title:   body.Title,
			Content: body.Content,
		},
	)

	if err != nil {
		return utils.Error(c, 500, err.Error())
	}

	return utils.Success(c, 201, "Note created successfully", note)
}

// GET /notes
func (handler *NoteHandler) GetAll(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	search := c.Query("search", "")

	notes, total, err := handler.service.GetAll(
		userID,
		page,
		limit,
		search,
	)

	if err != nil {
		return utils.Error(c, 500, err.Error())
	}

	return utils.Success(
		c,
		200,
		"Notes fetched successfully",
		fiber.Map{
			"notes": notes,
			"page":  page,
			"limit": limit,
			"total": total,
		},
	)
}

// GET /notes/:id
func (handler *NoteHandler) GetOne(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	noteID := c.Params("id")

	note, err := handler.service.GetByID(noteID, userID)

	if err != nil {
		return utils.Error(c, 404, err.Error())
	}

	return utils.Success(c, 200, "Note fetched successfully", note)
}

// PUT /notes/:id
func (handler *NoteHandler) Update(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	noteID := c.Params("id")

	var body NoteBody

	if err := c.BodyParser(&body); err != nil {
		return utils.Error(c, 400, "Invalid JSON body")
	}

	note, err := handler.service.Update(
		noteID,
		userID,
		services.UpdateNoteRequest{
			Title:   body.Title,
			Content: body.Content,
		},
	)

	if err != nil {
		return utils.Error(
			c,
			400,
			validators.FormatValidationError(err),
		)
	}

	return utils.Success(c, 200, "Note updated successfully", note)
}

// DELETE /notes/:id
func (handler *NoteHandler) Delete(c *fiber.Ctx) error {

	userID := c.Locals("userID").(string)

	noteID := c.Params("id")

	err := handler.service.Delete(noteID, userID)

	if err != nil {
		return utils.Error(
			c,
			400,
			validators.FormatValidationError(err),
		)
	}

	return utils.Success(c, 200, "Note deleted successfully", nil)
}
