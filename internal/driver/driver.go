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
	"encoding/hex"
	"fmt"
	"github.com/winc-link/hummingbird-sdk-go/commons"
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-sdk-go/service"
	"github.com/winc-link/hummingbird-tcp-driver/internal/device"
	"github.com/winc-link/hummingbird-tcp-driver/internal/server"
	"reflect"
	"strconv"
)

type TcpProtocolDriver struct {
	sd *service.DriverService
}

// CloudPluginNotify 云插件启动/停止通知
func (dr TcpProtocolDriver) CloudPluginNotify(ctx context.Context, t commons.CloudPluginNotifyType, name string) error {
	return nil
}

// DeviceNotify 设备添加/修改/删除通知
func (dr TcpProtocolDriver) DeviceNotify(ctx context.Context, t commons.DeviceNotifyType, deviceId string, dev model.Device) error {
	if t == commons.DeviceAddNotify {
		device.PutDevice(dev.DeviceSn, device.NewDevice(dev.Id, dev.DeviceSn, dev.ProductId, dev.Status == commons.DeviceOnline))
	}
	return nil
}

// ProductNotify 产品添加/修改/删除通知
func (dr TcpProtocolDriver) ProductNotify(ctx context.Context, t commons.ProductNotifyType, productId string, product model.Product) error {
	return nil
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
	devicesn := device.DeviceSn
	if data.Data["mStatus"] != nil {
		var sendData byte
		switch t := data.Data["mStatus"].(type) {
		default:
			typeOfA := reflect.TypeOf(t)
			dr.sd.GetLogger().Info(typeOfA.Name(), typeOfA.Kind())
		case float64:
			if t == 1 {
				sendData = 1
			} else if t == 0 {
				sendData = 0
			}
		case int:
			if t == 1 {
				sendData = 1
			} else if t == 0 {
				sendData = 0
			}
		case int64:
			if t == 1 {
				sendData = 1
			} else if t == 0 {
				sendData = 0
			}
		case string:
			if t == "1" {
				sendData = 1
			} else if t == "0" {
				sendData = 0
			}
		}
		//sendData = 0x01
		//请根据业务需要自行填写内容
		resBuff := make([]byte, 24)
		resBuff[0] = 65 //包头
		resBuff[1] = 76 //包头
		resBuff[2] = 77 // 包头

		resBuff[3] = 00 // 包长
		resBuff[4] = 23 // 包长

		//543204462026
		resBuff[5] = byte(Hex2Dec(devicesn[0:2]))    //mac地址 36
		resBuff[6] = byte(Hex2Dec(devicesn[2:4]))    //mac地址 20
		resBuff[7] = byte(Hex2Dec(devicesn[4:6]))    //mac地址 04
		resBuff[8] = byte(Hex2Dec(devicesn[6:8]))    //mac地址
		resBuff[9] = byte(Hex2Dec(devicesn[8:10]))   //mac地址
		resBuff[10] = byte(Hex2Dec(devicesn[10:12])) //mac地址

		resBuff[11] = 00 //设备类型
		resBuff[12] = 66 //128 ?// 设备类型

		resBuff[13] = 87 //数据头
		resBuff[14] = 67 //数据头

		resBuff[15] = 07 //总长度

		resBuff[16] = 66 //设备类型

		resBuff[17] = 145 //命令类型

		resBuff[18] = 02       //数据长度
		resBuff[19] = 03       //msgid
		resBuff[20] = sendData //继电器开关

		resBuff[21] = 00  //时间戳类型
		resBuff[22] = 138 //校验和
		resBuff[23] = 10  //结束符
		//65 76 77 0 23 102 85 68 51 34 17 0 66 87 67 7 66 145 2 3 0 0 138 10
		dr.sd.GetLogger().Info("Write data:", resBuff)
		_, err := conn.Conn.Write(resBuff)
		if err != nil {
			return err
		}
		encodedStr := hex.EncodeToString(resBuff)
		dr.sd.GetLogger().Info("Write data to string:", encodedStr)

		return nil

		//b, err := hex.DecodeString("414C4D001766554433221100425743074291020300008A0A")
		//if err != nil {
		//	return err
		//}
		//_, err = conn.Conn.Write(b)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}

// HandlePropertyGet 设备属性查询
func (dr TcpProtocolDriver) HandlePropertyGet(ctx context.Context, deviceId string, data model.PropertyGet) error {
	return nil
}

// HandleServiceExecute 设备服务调用
func (dr TcpProtocolDriver) HandleServiceExecute(ctx context.Context, deviceId string, data model.ServiceExecuteRequest) error {
	device, ok := dr.sd.GetDeviceById(deviceId)
	if !ok {
		return fmt.Errorf(" device [%s] not found", deviceId)
	}
	conn := server.GetTcpServer().GetConnectByDeviceSn(device.DeviceSn)
	if conn == nil {
		return fmt.Errorf(" device [%s] connection was not found", deviceId)
	}
	devicesn := device.DeviceSn

	switch data.Data.Code {
	case "collection_time":
		for k, v := range data.Data.InputParams {
			dr.sd.GetLogger().Info("v:", v)

			if k == "set_time" {
				//请根据业务需要自行填写内容
				resBuff := make([]byte, 25)
				resBuff[0] = 65 //包头
				resBuff[1] = 76 //包头
				resBuff[2] = 77 // 包头

				resBuff[3] = 00 // 包长
				resBuff[4] = 24 // 包长

				//543204462026
				resBuff[5] = byte(Hex2Dec(devicesn[0:2]))    //mac地址 36
				resBuff[6] = byte(Hex2Dec(devicesn[2:4]))    //mac地址 20
				resBuff[7] = byte(Hex2Dec(devicesn[4:6]))    //mac地址 04
				resBuff[8] = byte(Hex2Dec(devicesn[6:8]))    //mac地址
				resBuff[9] = byte(Hex2Dec(devicesn[8:10]))   //mac地址
				resBuff[10] = byte(Hex2Dec(devicesn[10:12])) //mac地址

				resBuff[11] = 00 //设备类型
				resBuff[12] = 66 //128 ?// 设备类型

				resBuff[13] = 87 //数据头
				resBuff[14] = 67 //数据头

				resBuff[15] = 8 //总长度

				resBuff[16] = 66 //设备类型

				resBuff[17] = 145 //命令类型

				resBuff[18] = 03 //数据长度
				resBuff[19] = 02 //计量插座参数 ID
				resBuff[20] = 00 //采集时间
				resBuff[21] = 14 //采集时间

				resBuff[22] = 00  //时间戳类型
				resBuff[23] = 138 //校验和
				resBuff[24] = 10  //结束符
				//65 76 77 0 23 102 85 68 51 34 17 0 66 87 67 7 66 145 2 3 0 0 138 10
				dr.sd.GetLogger().Info("Write data:", resBuff)
				_, err := conn.Conn.Write(resBuff)
				if err != nil {
					return err
				}
				encodedStr := hex.EncodeToString(resBuff)

				dr.sd.GetLogger().Info("Write data to string:", encodedStr)
			}
		}
	case "1":
	default:

	}
	return nil
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

func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int(n)
}
