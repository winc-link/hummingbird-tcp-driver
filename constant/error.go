/*******************************************************************************
 * Copyright 2017.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/

package constant

type ErrorMessage string

const (
	DefaultSuccess ErrorMessage = "success"

	SystemError               ErrorMessage = "system error"
	RpcRequestError           ErrorMessage = "rpc request error"
	ManyRequestsError         ErrorMessage = "too many requests"
	FormatError               ErrorMessage = "the format of result is error"
	DeviceNotFoundError       ErrorMessage = "device not found"
	ProductNotFoundError      ErrorMessage = "product not found"
	ReportDataRangeError      ErrorMessage = "data size is not within the defined range"
	PropertyReportTypeError   ErrorMessage = "data report type error"
	PropertyCodeNotFoundError ErrorMessage = "property code not found"
	EventCodeNotFoundError    ErrorMessage = "event code not found"
	ReportDataLengthError     ErrorMessage = "data length is greater than the defined"
	InvalidParameterError     ErrorMessage = "invalid Parameter"
	AuthPermissionDenyError   ErrorMessage = "auth PermissionDeny"
	InvalidTokenError         ErrorMessage = "invalid token"
	ConnectionNotFoundError   ErrorMessage = "device tcp connection not found"
)

type ErrorCode int

const (
	DefaultSuccessCode ErrorCode = 200

	SystemErrorCode             ErrorCode = 10001
	RpcRequestErrorCode         ErrorCode = 10002
	ManyRequestsErrorCode       ErrorCode = 10003
	FormatErrorCode             ErrorCode = 10004
	DeviceNotFound              ErrorCode = 20001
	ProductNotFound             ErrorCode = 30001
	ReportDataRangeErrorCode    ErrorCode = 40001
	PropertyReportTypeErrorCode ErrorCode = 40002
	PropertyCodeNotFound        ErrorCode = 40003
	EventCodeNotFound           ErrorCode = 40004
	ReportDataLengthErrorCode   ErrorCode = 40005
	InvalidParameterErrorCode   ErrorCode = 40006
	AuthPermissionDenyErrorCode ErrorCode = 40007
	InvalidTokenErrorCode       ErrorCode = 40008
	ConnectionNotFoundErrorCode ErrorCode = 40009
)

var ErrorCodeMsgMap = map[ErrorCode]ErrorMessage{
	DefaultSuccessCode:          DefaultSuccess,
	SystemErrorCode:             SystemError,
	RpcRequestErrorCode:         RpcRequestError,
	ManyRequestsErrorCode:       ManyRequestsError,
	FormatErrorCode:             FormatError,
	DeviceNotFound:              DeviceNotFoundError,
	ProductNotFound:             ProductNotFoundError,
	ReportDataRangeErrorCode:    ReportDataRangeError,
	PropertyReportTypeErrorCode: PropertyReportTypeError,
	PropertyCodeNotFound:        PropertyCodeNotFoundError,
	EventCodeNotFound:           EventCodeNotFoundError,
	ReportDataLengthErrorCode:   ReportDataLengthError,
	InvalidParameterErrorCode:   InvalidParameterError,
	AuthPermissionDenyErrorCode: AuthPermissionDenyError,
	InvalidTokenErrorCode:       InvalidTokenError,
	ConnectionNotFoundErrorCode: ConnectionNotFoundError,
}
