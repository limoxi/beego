package mns

import (
	"net/http"
	"errors"
	"fmt"
)

func send(client MNSClient, decoder MNSDecoder, method Method, headers map[string]string, message interface{}, resource string, v interface{}) (statusCode int, err error) {
	var resp *http.Response
	if resp, err = client.Send(method, headers, message, resource); err != nil {
		return
	}

	if resp != nil {
		defer resp.Body.Close()
		statusCode = resp.StatusCode

		if resp.StatusCode != http.StatusCreated &&
			resp.StatusCode != http.StatusOK &&
			resp.StatusCode != http.StatusNoContent {

			errResp := ErrorMessageResponse{}
			if e := decoder.Decode(resp.Body, &errResp); e != nil {
				err = e
				return
			}
			err = ParseError(errResp, resource)
			return
		}

		if v != nil {
			if e := decoder.Decode(resp.Body, v); e != nil {
				err = err
				return
			}
		}
	}

	return
}

func ParseError(resp ErrorMessageResponse, resource string) (err error) {
	errMsg := fmt.Sprintf("ali_mns response status error,code: %s, message: %s, resource: %s, request id: %s, host id: %s", resp.Code, resp.Message, resource, resp.RequestId, resp.HostId)
	err = errors.New(errMsg)
	return
}
