package go_huawei

import (
	"fmt"
)

type GoHuaweiError struct {
	IsApiError bool
	ReturnCode ReturnCode
	ReturnDesc ReturnDesc

	Context     string
	RawRequest  []byte
	RawResponse []byte

	Err error
}

func (e *GoHuaweiError) AddRawRequest(s []byte) {
	e.RawRequest = s
}

func (e *GoHuaweiError) AddRawResponse(s []byte) {
	e.RawResponse = s
}

func NewGoHuaweiError(context string, err error) *GoHuaweiError {
	return &GoHuaweiError{IsApiError: false, Context: context, Err: err}
}

func NewGoHuaweiApiError(returnCode ReturnCode, returnDesc ReturnDesc, context string) *GoHuaweiError {
	return &GoHuaweiError{IsApiError: true, ReturnCode: returnCode, ReturnDesc: returnDesc, Context: context}
}

func (e *GoHuaweiError) Error() string {
	startString := "go-huawei: "
	if e.Context != "" {
		startString += e.Context + ": "
	}

	if e.IsApiError {
		if e.ReturnCode != ReturnCodeOK && e.ReturnDesc != ReturnDescZeroResults {
			return fmt.Sprintf("%s%s - %s", startString, e.ReturnCode, e.ReturnDesc)
		}
	}

	return fmt.Sprintf("%s%s", startString, e.Err.Error())
}
