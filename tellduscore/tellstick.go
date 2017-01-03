package tellduscore

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
