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

var GlobalDriverService *service.DriverService

type TcpProtocolDriver struct{}

// CloudPluginNotify 云插件启动/停止通知
func (t2 TcpProtocolDriver) CloudPluginNotify(ctx context.Context, t commons.CloudPluginNotifyType, name string) error {
	//TODO implement me
	panic("implement me")
}

// DeviceNotify 设备添加/修改/删除通知
func (t2 TcpProtocolDriver) DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, deviceId string, device model.Device) error {
	//TODO implement me
	panic("implement me")
}

// ProductNotify 产品添加/修改/删除通知
func (t2 TcpProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	//TODO implement me
	panic("implement me")
}

// Stop 蜂鸟物联网平台通知
func (t2 TcpProtocolDriver) Stop(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

// HandlePropertySet 设备属性设置
func (t2 TcpProtocolDriver) HandlePropertySet(ctx context.Context, deviceId string, data model.PropertySet) error {
	device, ok := GlobalDriverService.GetDeviceById(deviceId)
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
func (t2 TcpProtocolDriver) HandlePropertyGet(ctx context.Context, deviceId string, data model.PropertyGet) error {
	//TODO implement me
	panic("implement me")
}

// HandleServiceExecute 设备服务调用
func (t2 TcpProtocolDriver) HandleServiceExecute(ctx context.Context, deviceId string, data model.ServiceExecuteRequest) error {
	//TODO implement me
	panic("implement me")
}

// NewTcpProtocolDriver Tcp协议驱动
func NewTcpProtocolDriver(ctx context.Context, sd *service.DriverService) *TcpProtocolDriver {
	GlobalDriverService = sd
	loadDevices()
	go server.GetTcpServer().Start(server.TcpDataHandler)
	go cancel(ctx)
	return &TcpProtocolDriver{}
}

// loadDevices 获取所有已经创建成功的设备，保存在内存中。
func loadDevices() {
	for _, dev := range GlobalDriverService.GetDeviceList() {
		device.NewDevice(dev.Id, dev.DeviceSn, dev.ProductId, dev.Status == commons.DeviceOnline)
	}
}

// cancel 监听驱动退出，如果驱动退出则把此驱动关联的设备设置成离线
func cancel(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			for _, dev := range device.GetAllDevice() {
				_ = dev.Offline()
			}
		}
	}
}
