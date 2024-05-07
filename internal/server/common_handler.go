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

package server

import (
	"encoding/json"
	"fmt"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-tcp-driver/constant"
	"github.com/winc-link/hummingbird-tcp-driver/dtos"
	"github.com/winc-link/hummingbird-tcp-driver/protocol"
	"github.com/winc-link/hummingbird-tcp-driver/tool"
	"net"
	"time"
)

// deviceAuthHandler 设备连接鉴权
func deviceAuthHandler(data []byte, conn net.Conn) error {
	authDataInfo, err := dtos.BytesToAuthStruct(data)
	if err != nil {
		return err
	}
	device, ok := GlobalDriverService.GetDeviceById(authDataInfo.DeviceId)
	if !ok {
		GlobalDriverService.GetLogger().Errorf("device [%s] auth error device not found ", authDataInfo.DeviceId)
		return fmt.Errorf("device [%s] auth error device not found ", authDataInfo.DeviceId)
	}
	product, ok := GlobalDriverService.GetProductById(authDataInfo.ProductId)
	if !ok {
		GlobalDriverService.GetLogger().Errorf("device [%s] auth error product not found ", authDataInfo.DeviceId)
		return fmt.Errorf("device [%s] auth error product not found ", authDataInfo.DeviceId)
	}

	if device.ProductId != product.Id {
		GlobalDriverService.GetLogger().Errorf("device [%s] auth error device product relationship error", authDataInfo.DeviceId)
		return fmt.Errorf("device [%s] auth error device product relationship error", authDataInfo.DeviceId)
	}

	if authDataInfo.Token != tool.HmacMd5(device.Secret, device.Id+"&"+product.Key) {
		GlobalDriverService.GetLogger().Errorf("device [%s] auth error auth permission deny", authDataInfo.DeviceId)
		return fmt.Errorf("device [%s] auth error auth permission deny", authDataInfo.DeviceId)
	}

	err = GlobalDriverService.Online(authDataInfo.DeviceId)
	if err != nil {
		GlobalDriverService.GetLogger().Errorf("device [%s] auth online error: ", authDataInfo.DeviceId, err.Error())
		return fmt.Errorf("device [%s] auth online error: %s", authDataInfo.DeviceId, err.Error())
	}
	GlobalDriverService.GetLogger().Infof("device [%s] auth success", authDataInfo.DeviceId)
	GlobalDriverService.GetLogger().Infof("device auth remote addr:[%s]", conn.RemoteAddr().String())

	if tcpServer.ClientCons[authDataInfo.DeviceId] == nil {
		tcpServer.ClientCons[authDataInfo.DeviceId] = &Connect{
			Conn: conn,
		}
	}
	if tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()] == nil {
		tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()] = &device
	}
	return nil
}

// devicePropertyReportHandler 设备属性上报
func devicePropertyReportHandler(data []byte, deviceId string, conn net.Conn, sequenceID string) error {
	propertyPost, err := dtos.BytesToPropertyStruct(data)
	if err != nil {
		return err
	}
	var response dtos.Response
	device, ok := GlobalDriverService.GetDeviceById(deviceId)
	if !ok {
		GlobalDriverService.GetLogger().Errorf("device [%s] device not found ", deviceId)
		if propertyPost.Sys.Ack {
			response.Code = int(constant.DeviceNotFound)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.DeviceNotFound])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.PropertyReportReplyOperation))
		}
		return fmt.Errorf("device [%s] device not found ", deviceId)
	}

	var delPropertyCode []string
	for code, param := range propertyPost.Params {
		if param.Time == 0 {
			propertyPost.Params[code] = model.PropertyData{
				Time:  time.Now().UnixMilli(),
				Value: param.Value,
			}

		}
	}
	filterPropertyPost := propertyPost.Params
	for _, code := range delPropertyCode {
		delete(filterPropertyPost, code)
	}
	propertyPost.Params = filterPropertyPost
	if len(propertyPost.Params) == 0 {
		if propertyPost.Sys.Ack {
			response.Code = int(constant.PropertyCodeNotFound)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.PropertyCodeNotFound])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.PropertyReportReplyOperation))
		}
		return fmt.Errorf("device [%s] device not found ", deviceId)
	}
	_, err = GlobalDriverService.PropertyReport(device.Id, model.NewPropertyReport(propertyPost.Sys.Ack, propertyPost.Params))

	if propertyPost.Sys.Ack {
		if err == nil {
			response.Code = int(constant.DefaultSuccessCode)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.DefaultSuccessCode])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.PropertyReportReplyOperation))
		} else {
			GlobalDriverService.GetLogger().Errorf("device [%s] device report property err: %s", deviceId, err.Error())
			response.Code = int(constant.SystemErrorCode)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.SystemErrorCode])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.PropertyReportReplyOperation))
		}
	}
	return nil
}

// deviceEventReportHandler 设备事件上报
func deviceEventReportHandler(data []byte, deviceId string, conn net.Conn, sequenceID string) error {
	var response dtos.Response
	eventReport, err := dtos.BytesToEventStruct(data)
	if err != nil {
		if eventReport.Sys.Ack {
			response.Code = int(constant.FormatErrorCode)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.FormatErrorCode])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.EventReportReplyOperation))
		}
		return err
	}

	//检查设备是否存在
	device, ok := GlobalDriverService.GetDeviceById(deviceId)
	if !ok {
		if eventReport.Sys.Ack {
			response.Code = int(constant.DeviceNotFound)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.DeviceNotFound])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.EventReportReplyOperation))
		}
		return fmt.Errorf("device [%s] device not found ", deviceId)
	}

	// 检查产品是否存在
	if _, ok := GlobalDriverService.GetProductById(device.ProductId); !ok {
		if eventReport.Sys.Ack {
			response.Code = int(constant.ProductNotFound)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.ProductNotFound])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.EventReportReplyOperation))
		}
		return fmt.Errorf("device [%s] product not found", device.Id)
	}
	// 检测上报事件标识符
	_, ok = GlobalDriverService.GetProductEventByCode(device.ProductId, eventReport.Params.EventCode)
	if !ok {
		if eventReport.Sys.Ack {
			response.Code = int(constant.EventCodeNotFound)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.EventCodeNotFound])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.EventReportReplyOperation))
		}
		return fmt.Errorf("device [%s] event code not found", device.Id)
	}

	if eventReport.Params.EventTime == 0 {
		eventReport.Params.EventTime = time.Now().UnixMilli()
	}
	_, err = GlobalDriverService.EventReport(deviceId, model.NewEventReport(eventReport.Sys.Ack, eventReport.Params))

	if eventReport.Sys.Ack {
		if err == nil {
			response.Code = int(constant.DefaultSuccessCode)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.DefaultSuccessCode])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.EventReportReplyOperation))
		} else {
			GlobalDriverService.GetLogger().Errorf("device [%s] device report event err: %s", deviceId, err.Error())
			response.Code = int(constant.SystemErrorCode)
			response.Msg = string(constant.ErrorCodeMsgMap[constant.SystemErrorCode])
			b, _ := json.Marshal(response)
			_, _ = conn.Write(protocol.Packet(b, sequenceID, constant.EventReportReplyOperation))
		}
	}
	return nil
}

// devicePropertySetReply 设备属性设置响应
func devicePropertySetReplyHandler(data []byte, deviceId string, sequenceID string) error {
	propertySet, err := dtos.BytesToPropertySetReplyStruct(data)
	if err != nil {
		return err
	}
	if err = GlobalDriverService.PropertySetResponse(deviceId, model.PropertySetResponse{
		MsgId: sequenceID,
		Data:  propertySet.Params,
	}); err != nil {
		return err
	}
	return nil
}

// devicePropertyGetReply 设备属性查询响应
func devicePropertyGetReplyHandler(data []byte, deviceId string, sequenceID string) error {
	propertyGet, err := dtos.BytesToPropertyGetReplyStruct(data)
	if err != nil {
		return err
	}
	if err = GlobalDriverService.PropertyGetResponse(deviceId, model.PropertyGetResponse{
		MsgId: sequenceID,
		Data:  propertyGet.Params,
	}); err != nil {
		return err
	}
	return nil
}

// deviceServiceExecuteReplyHandler 设备服务调用响应
func deviceServiceExecuteReplyHandler(data []byte, deviceId string, sequenceID string) error {
	serviceExecute, err := dtos.BytesToServiceExecuteReplyStruct(data)
	if err != nil {
		return err
	}
	if err = GlobalDriverService.ServiceExecuteResponse(deviceId, model.ServiceExecuteResponse{
		MsgId: sequenceID,
		Data:  serviceExecute.Params,
	}); err != nil {
		return err
	}
	return nil
}
