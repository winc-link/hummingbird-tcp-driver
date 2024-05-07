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

package driver

import (
	"context"
	"encoding/json"
	"github.com/winc-link/hummingbird-sdk-go/commons"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"github.com/winc-link/hummingbird-tcp-driver/constant"
	"github.com/winc-link/hummingbird-tcp-driver/dtos"
	"github.com/winc-link/hummingbird-tcp-driver/internal/server"
	"github.com/winc-link/hummingbird-tcp-driver/protocol"
)

type TcpProtocolDriver struct {
	sd *service.DriverService
}

// CloudPluginNotify 云插件启动/停止通知
func (dr TcpProtocolDriver) CloudPluginNotify(ctx context.Context, t commons.CloudPluginNotifyType, name string) error {
	return nil
}

// DeviceNotify 设备添加/修改/删除通知
func (dr TcpProtocolDriver) DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, deviceId string, device model.Device) error {
	return nil
}

// ProductNotify 产品添加/修改/删除通知
func (dr TcpProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	return nil
}

// Stop 驱动退出通知。
func (dr TcpProtocolDriver) Stop(ctx context.Context) error {
	for _, device := range dr.sd.GetDeviceList() {
		dr.sd.Offline(device.Id)
	}
	return nil
}

// HandlePropertySet 设备属性设置
func (dr TcpProtocolDriver) HandlePropertySet(ctx context.Context, deviceId string, data model.PropertySet) error {
	device, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constant.DeviceNotFound),
				ErrorMessage: string(constant.ErrorCodeMsgMap[constant.DeviceNotFound]),
			},
		})
		return nil

	}
	conn := server.GetTcpServer().GetConnectByDeviceId(device.Id)
	if conn == nil {
		_ = dr.sd.PropertySetResponse(deviceId, model.PropertySetResponse{
			MsgId: data.MsgId,
			Data: model.PropertySetResponseData{
				Success:      false,
				Code:         uint32(constant.ConnectionNotFoundErrorCode),
				ErrorMessage: string(constant.ErrorCodeMsgMap[constant.ConnectionNotFoundErrorCode]),
			},
		})

		return nil
	}
	var propertySet dtos.PropertySet
	propertySet.Params = make(map[string]interface{})
	propertySet.Params = data.Data
	b, _ := json.Marshal(propertySet)
	resBuff := protocol.Packet(b, data.MsgId, constant.PropertySet)
	_, err := conn.Conn.Write(resBuff)
	if err != nil {
		return err
	}
	return nil
}

// HandlePropertyGet 设备属性查询
func (dr TcpProtocolDriver) HandlePropertyGet(ctx context.Context, deviceId string, data model.PropertyGet) error {
	device, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.PropertyGetResponse(deviceId, model.PropertyGetResponse{
			MsgId: data.MsgId,
			Data:  []model.PropertyGetResponseData{},
		})
		return nil
	}
	conn := server.GetTcpServer().GetConnectByDeviceId(device.Id)
	if conn == nil {
		_ = dr.sd.PropertyGetResponse(deviceId, model.PropertyGetResponse{
			MsgId: data.MsgId,
			Data:  []model.PropertyGetResponseData{},
		})
		return nil
	}

	b, _ := json.Marshal(data.Data)
	resBuff := protocol.Packet(b, data.MsgId, constant.PropertyGet)
	_, err := conn.Conn.Write(resBuff)
	if err != nil {
		return err
	}
	return nil
}

// HandleServiceExecute 设备服务调用
func (dr TcpProtocolDriver) HandleServiceExecute(ctx context.Context, deviceId string, data model.ServiceExecuteRequest) error {
	device, ok := dr.sd.GetDeviceById(deviceId)
	if ok != true {
		_ = dr.sd.ServiceExecuteResponse(deviceId, model.ServiceExecuteResponse{
			MsgId: data.MsgId,
			Data:  model.ServiceDataOut{},
		})
		return nil
	}
	conn := server.GetTcpServer().GetConnectByDeviceId(device.Id)
	if conn == nil {
		_ = dr.sd.ServiceExecuteResponse(deviceId, model.ServiceExecuteResponse{
			MsgId: data.MsgId,
			Data:  model.ServiceDataOut{},
		})
		return nil
	}

	b, _ := json.Marshal(data.Data)
	resBuff := protocol.Packet(b, data.MsgId, constant.ServiceExecute)
	_, err := conn.Conn.Write(resBuff)
	if err != nil {
		return err
	}
	return nil
}

// NewTcpProtocolDriver Tcp协议驱动
func NewTcpProtocolDriver(sd *service.DriverService) *TcpProtocolDriver {
	go server.GetTcpServer().Start(sd)
	return &TcpProtocolDriver{
		sd: sd,
	}
}
