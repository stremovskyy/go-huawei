package go_huawei

type ReturnCode string

const (
	// Success
	ReturnCodeOK ReturnCode = "0"
	// Invalid API key or token
	ReturnCodeInvalidAPI ReturnCode = "6"
	// The requested URL is incorrect
	ReturnCodeRequestedURLIncorrect ReturnCode = "5"

	// The authentication service is abnormal
	ReturnCodeInvalidAuth ReturnCode = "010001"
	// Route data does not exist
	ReturnCodeRouteDataDoesNotExist ReturnCode = "010008"
	// Route data does not exist
	ReturnCodeAPICallQuotaUsedUp ReturnCode = "010006"
	// Route data does not exist
	ReturnCodeAuthenticationFailed ReturnCode = "010005"
	// Invalid parameter
	ReturnCodeInvalidRequest ReturnCode = "010010"
	// Invalid parameter
	ReturnCodeUnauthorizedAPICall ReturnCode = "010003"
	// The linear distance between the departure  place and destination exceeds the upper limit.
	ReturnCodeDistanceBetweenExceedsLimit ReturnCode = "010020"

	ReturnCodeInternalServiceError ReturnCode = "110"
)
