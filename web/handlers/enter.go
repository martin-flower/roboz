package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/martin-flower/roboz-web/database"
	"github.com/martin-flower/roboz-web/service"
	"github.com/martin-flower/roboz-web/service/clean/intersection"
	"github.com/martin-flower/roboz-web/service/clean/intmap"
	"github.com/martin-flower/roboz-web/service/clean/simplest"
	"github.com/martin-flower/roboz-web/service/clean/sortedset"
	"github.com/martin-flower/roboz-web/service/direction"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Enter
// @Summary post cleaning instructions
// @Description post instructions as robot path, return 200 and some interesting results
// @Tags enter
// @Accept json
// @Param message body handlers.EnterPost true "instructions"
// @Produce json
// @Success 200 {object} handlers.EnterResponse true "result of instructions"
// @Failure 400 {string}
// @Failure 500 {string}
// @Router /developer-test/enter-path [post]
func Enter(ctx *fiber.Ctx) (err error) {

	zap.S().Debugf("roboz enter - %d", time.Now().UnixMilli())

	request := EnterPost{}
	err = ctx.BodyParser(&request)
	if err != nil {
		zap.S().Infof("body parse failed %+v", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request " + err.Error())
	}

	zap.S().Debugf("request:%+v", request)

	err = validate(request)
	if err != nil {
		zap.S().Infof("validate failed %+v", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request " + err.Error())
	}

	zap.S().Debugf("valid request:%+v", request)

	// accept input from both commands and commmands (with 3 mmm)
	commands := []service.Command{}
	for _, c := range request.Commands {
		if !direction.Valid(c.Direction) {
			return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid direction %s", c.Direction))
		}
		commands = append(commands, service.Command{Steps: c.Steps, Direction: direction.FromString(c.Direction)})
	}
	for _, c := range request.Commmands {
		if !direction.Valid(c.Direction) {
			return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("invalid direction %s", c.Direction))
		}
		commands = append(commands, service.Command{Steps: c.Steps, Direction: direction.FromString(c.Direction)})
	}

	startTime := time.Now()

	var cleaner service.Cleaner

	// algorithm is set in roboz.yaml
	switch algorithm := viper.GetString("algorithm"); algorithm {
	case "intersection":
		zap.S().Debugf("algorithm:intersection")
		if len(commands) > 10000 {
			return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("intersection algorithm cannot handle more than 10000 commmands (you sent %d)", len(commands)))
		}
		cleaner = intersection.Cleaner{}
	case "intmap":
		zap.S().Debugf("algorithm:intmap")
		if len(commands) > 10000 {
			return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("intmap algorithm cannot handle more than 10000 commmands (you sent %d)", len(commands)))
		}
		cleaner = intmap.Cleaner{}
	case "simplest":
		zap.S().Debugf("algorithm:simplest")
		if len(commands) > 3 {
			return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("simplest algorithm cannot handle more than 3 commmands (you sent %d)", len(commands)))
		}
		cleaner = simplest.Cleaner{}
	case "sortedset":
		zap.S().Debugf("algorithm:sortedset")
		if len(commands) > 500 {
			return ctx.Status(fiber.StatusBadRequest).SendString(fmt.Sprintf("sortedset algorithm cannot handle more than 500 commmands (you sent %d)", len(commands)))
		}
		cleaner = sortedset.Cleaner{}
	default:
		zap.S().Fatalf("unrecognized algorithm %s", algorithm)
	}

	cleaned := cleaner.Clean(service.Coordinate(request.Start), commands)
	duration := time.Since(startTime)

	ID, timestamp, err := database.Store(len(commands), cleaned, duration.Seconds())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("server error") // do not expose the actual error
	}

	response := EnterResponse{}
	response.ID = ID
	response.Timestamp = timestamp
	response.Commands = len(commands)
	response.Result = cleaned
	response.Duration = duration.Seconds()

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func validate(request EnterPost) (err error) {
	// with a bit of imagination we could add some more validation in here
	if len(request.Commands) < 1 && len(request.Commmands) < 1 {
		err = fmt.Errorf("commands missing from request")
	}
	return err
}

type Coordinate struct {
	X int
	Y int
}

type Command struct {
	Direction string
	Steps     int
}

type EnterPost struct {
	Start     Coordinate
	Commands  []Command
	Commmands []Command // according to the specification
}

type EnterResponse struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Commands  int       `json:"commands"` // note specification is for commmands, not commands
	Result    int       `json:"result"`
	Duration  float64   `json:"duration"`
}
