package api

// ErrorCode is the template for error codes
type ErrorCode struct {
	code          int
	humanReadable string
}

// ErrorCodes contains all known API error codes
var ErrorCodes = map[string]ErrorCode{
	"INVALID_CREDENTIALS": ErrorCode{
		401,
		"Sorry, your credentials are invalid. Please try again.",
	},
	"NORDSYD_API_NOT_WORKING": ErrorCode{
		998,
		"The Nordsyd API is currently expericencing issues. Please try again later.",
	},
	"UNKNOWN_ERROR_CODE": ErrorCode{
		999,
		"The Nordsyd API returned an unknown error code. Are you running the latest CLI version?",
	},
}

// GetErrorCode finds and returns an error code given the ID
func GetErrorCode(key string) ErrorCode {
	value, ok := ErrorCodes[key]

	if !ok {
		return ErrorCodes["UNKNOWN_ERROR_CODE"]
	}

	return value
}
