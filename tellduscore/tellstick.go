package tellduscore

import (
	"fmt"
	"strconv"
	"strings"
)

// TellstickDim constant for dimming
const TellstickDim = 1

// TellstickTurnon constant for turn on
const TellstickTurnon = 2

// TellstickTurnoff constant for turn off
const TellstickTurnoff = 3

// TellstickDimString constant for dimming
const TellstickDimString = "dim"

// TellstickTurnonString constant for turn on
const TellstickTurnonString = "turnon"

// TellstickTurnoffString constant for turn off
const TellstickTurnoffString = "turnoff"

// TellstickLearnString constant for learn
const TellstickLearnString = "learn"

// GetResultMessage returns a tellstick result message
func GetResultMessage(tellResult int) string {
	resultType := "UNKNOWN ERROR"
	switch tellResult {
	case 0:
		resultType = "TELLSTICK_SUCCESS"
	case -1:
		resultType = "TELLSTICK_ERROR_NOT_FOUND"
		break
	case -2:
		resultType = "TELLSTICK_ERROR_PERMISSION_DENIED"
		break
	case -3:
		resultType = "TELLSTICK_ERROR_DEVICE_NOT_FOUND"
		break
	case -4:
		resultType = "TELLSTICK_ERROR_METHOD_NOT_SUPPORTED"
		break
	case -5:
		resultType = "TELLSTICK_ERROR_COMMUNICATION"
		break
	case -6:
		resultType = "TELLSTICK_ERROR_CONNECTING_SERVICE"
		break
	case -7:
		resultType = "TELLSTICK_ERROR_UNKNOWN_RESPONSE"
		break
	case -8:
		resultType = "TELLSTICK_ERROR_SYNTAX"
		break
	case -9:
		resultType = "TELLSTICK_ERROR_BROKEN_PIPE"
		break
	case -10:
		resultType = "TELLSTICK_ERROR_COMMUNICATING_SERVICE"
		break
	case -11:
		resultType = "TELLSTICK_ERROR_CONFIG_SYNTAX"
		break
	case -99:
		resultType = "TELLSTICK_ERROR_UNKNOWN"
		break
	}

	return resultType
}

// GetTellstickMessageLevel returns a formatted string for Tellstick commands
func GetTellstickMessageLevel(message string, id int, level int) string {
	return fmt.Sprintf("%d:%si%dsi%ds", len(message), message, id, level)
}

// GetTellstickMessage returns a formatted string for Tellstick commands
func GetTellstickMessage(message string, id int) string {
	return fmt.Sprintf("%d:%si%ds", len(message), message, id)
}

// GetIntFromResult parses i%ds-messages and returns the integer
func GetIntFromResult(result string) int {
	i, err := strconv.Atoi(result[1:strings.Index(result, "s")])
	if err != nil {
		return -1
	}

	return i
}
