package handlers

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/martin-flower/roboz-web/database"
	"go.uber.org/zap"
)

// List
// @Summary return list of cleanings
// @Description list of cleanings - one row per cleaning - in ascending order
// @Tags list
// @Produce json
// @Param offset query int false "how many cleanings to skip - defaults to 0"
// @Param limit query int false "maximum number of cleanings - defaults to 10 - maximum is 20"
// @Success 200 {object} handlers.ListResponse true "list of instructions"
// @Failure 400 {string}
// @Failure 500 {string}
// @Router /list [get]
func List(c *fiber.Ctx) (err error) {
	zap.S().Info("roboz list")

	offset := 0
	offsetString := c.Query("offset")
	if len(offsetString) > 0 {
		offset, err = strconv.Atoi(offsetString)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid query parameter offset " + err.Error())
		}
		if offset < 0 {
			offset = 0
		}
	}

	const maxcleanings = 20
	limit := maxcleanings
	limitString := c.Query("limit")
	if len(limitString) > 0 {
		limit, err = strconv.Atoi(limitString)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("invalid query parameter limit " + err.Error())
		}
		if limit > maxcleanings {
			limit = maxcleanings
		}
		if limit < 1 {
			limit = 1
		}
	}

	var rows []database.Row
	rows, err = database.List(offset, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("server error") // + err.Error()) // do not expose the server error
	}

	listResponse := []EnterResponse{}
	for _, row := range rows {
		listResponse = append(listResponse, EnterResponse{ID: row.ID, Timestamp: row.Timestamp, Commands: row.Commands, Result: row.Result, Duration: row.Duration})
	}

	return c.Status(fiber.StatusOK).JSON(listResponse)
}

type ListResponse struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Commands  int       `json:"commands"` // note specification is for commmands, not commands
	Result    int       `json:"result"`
	Duration  float64   `json:"duration"`
}
