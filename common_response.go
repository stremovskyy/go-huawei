package go_huawei

import "fmt"

type CommonResponse struct {
	ReturnCode ReturnCode `json:"returnCode"`
	ReturnDesc ReturnDesc `json:"returnDesc"`
}

func (c *CommonResponse) StatusError() error {
	if c.ReturnCode != ReturnCodeOK && c.ReturnDesc != ReturnDescZeroResults {
		return fmt.Errorf("map-kit: %s - %s", c.ReturnCode, c.ReturnDesc)
	}

	return nil
}
