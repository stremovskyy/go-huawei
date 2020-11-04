package go_huawei

type ReturnCode string

const (
	// Success
	ReturnCodeOK ReturnCode = "0"
	// The authentication service is abnormal
	ReturnCodeInvalidAuth ReturnCode = "010001"
	// Invalid parameter
	ReturnCodeInvalidRequest ReturnCode = "010010"
)
