package go_huawei

type ReturnDesc string

const (
	ReturnDescOK             ReturnDesc = "OK"
	ReturnDescInvalidRequest ReturnDesc = "INVALID_REQUEST"
	ReturnDescUnknownError   ReturnDesc = "UNKNOWN_ERROR"
	ReturnDescNotFound       ReturnDesc = "NOT_FOUND"
	ReturnDescZeroResults    ReturnDesc = "ZERO_RESULTS"
	ReturnDescRequestDenied  ReturnDesc = "REQUEST_DENIED"
	ReturnDescOverQueryLimit ReturnDesc = "OVER_QUERY_LIMIT"
)
