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
	"fmt"
	"github.com/winc-link/hummingbird-sdk-go/commons"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"github.com/winc-link/hummingbird-tcp-driver/internal/device"
	"github.com/winc-link/hummingbird-tcp-driver/internal/server"
)

type TcpProtocolDriver struct {
	sd *service.DriverService
}

// CloudPluginNotify 云插件启动/停止通知
func (dr TcpProtocolDriver) CloudPluginNotify(ctx context.Context, t commons.CloudPluginNotifyType, name string) error {
	//TODO implement me
	panic("implement me")
}

// DeviceNotify 设备添加/修改/删除通知
func (dr TcpProtocolDriver) DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, deviceId string, device model.Device) error {
	//TODO implement me
	panic("implement me")
}

// ProductNotify 产品添加/修改/删除通知
func (dr TcpProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	//TODO implement me
	panic("implement me")
}

// Stop 驱动退出通知。
func (dr TcpProtocolDriver) Stop(ctx context.Context) error {
	for _, dev := range device.GetAllDevice() {
		dr.sd.Offline(dev.GetDeviceId())
	}
	return nil
}

// HandlePropertySet 设备属性设置
func (dr TcpProtocolDriver) HandlePropertySet(ctx context.Context, deviceId string, data model.PropertySet) error {
	device, ok := dr.sd.GetDeviceById(deviceId)
	if !ok {
		return fmt.Errorf(" device [%s] not found", deviceId)
	}
	conn := server.GetTcpServer().GetConnectByDeviceSn(device.DeviceSn)
	if conn == nil {
		return fmt.Errorf(" device [%s] connection was not found", deviceId)
	}
	//请根据业务需要自行填写内容
	resBuff := make([]byte, 3)
	resBuff[0] = 0x01
	resBuff[1] = 0x02
	resBuff[2] = 0x03
	_, err := conn.Conn.Write(resBuff)
	if err != nil {
		return err
	}
	return nil
}

// HandlePropertyGet 设备属性查询
func (dr TcpProtocolDriver) HandlePropertyGet(ctx context.Context, deviceId string, data model.PropertyGet) error {
	//TODO implement me
	panic("implement me")
}

// HandleServiceExecute 设备服务调用
func (dr TcpProtocolDriver) HandleServiceExecute(ctx context.Context, deviceId string, data model.ServiceExecuteRequest) error {
	//TODO implement me
	panic("implement me")
}

// NewTcpProtocolDriver Tcp协议驱动
func NewTcpProtocolDriver(sd *service.DriverService) *TcpProtocolDriver {
	loadDevices(sd)
	go server.GetTcpServer().Start(sd, server.TcpDataHandler)
	return &TcpProtocolDriver{
		sd: sd,
	}
}

// loadDevices 获取所有已经创建成功的设备，保存在内存中。
func loadDevices(sd *service.DriverService) {
	for _, dev := range sd.GetDeviceList() {
		device.PutDevice(dev.DeviceSn, device.NewDevice(dev.Id, dev.DeviceSn, dev.ProductId, dev.Status == commons.DeviceOnline))
	}
}
