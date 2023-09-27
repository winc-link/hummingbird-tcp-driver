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
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-tcp-driver/internal/device"
	"net"
	"strconv"
)

type TcpDataHandlers func(deviceSn string, data []byte) (retBuff []byte, err error)

// serverConnHandler 用户可以根据项目需要自行修改此方法的业务逻辑！
func serverConnHandler(conn net.Conn, tdh TcpDataHandlers) {
	defer conn.Close()

	var deviceSn string
	for {
		var buf [8]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			GlobalDriverService.GetLogger().Errorf("Read from tcp server failed,err:", err)
			closeConn(deviceSn)
			break
		}
		deviceSn = strconv.Itoa(int(buf[0]))
		if tcpServer.ClientCons[deviceSn] == nil {
			tcpServer.ClientCons[deviceSn] = &Connect{
				Conn: conn,
			}
		}

		bytes, err := tdh(deviceSn, buf[:n])
		if err != nil {
			return
		}
		conn.Write(bytes)
	}
}

func closeConn(deviceSn string) {
	dev, err := device.GetDevice(deviceSn)
	if err != nil {
		GlobalDriverService.Offline(dev.GetDeviceId())
	}
	tcpServer.DeleteClientByDeviceId(deviceSn)
}

// TcpDataHandler tcp数据处理
func TcpDataHandler(deviceSn string, data []byte) (retBuff []byte, err error) {
	dev, err := device.GetDevice(deviceSn)
	if err != nil {
		//新设备，做创建设备并且上线的业务逻辑。
		var (
			deviceName  = ""
			productId   = ""
			description = ""
			external    = map[string]string{}
		)

		newDevice, err := GlobalDriverService.CreateDevice(model.NewAddDevice(deviceName, productId, deviceSn, description, external))
		if err != nil {
			GlobalDriverService.GetLogger().Errorf("Create device [%s] err:", deviceSn, err.Error())
			return nil, err
		}
		//设备上线
		if err = GlobalDriverService.Online(newDevice.Id); err != nil {
			GlobalDriverService.GetLogger().Errorf("Device online [%s] err:", deviceSn, err.Error())
		}
		//把设备注册到device manage中。
		dev = device.NewDevice(newDevice.Id, deviceSn, newDevice.ProductId, true)
		device.PutDevice(deviceSn, dev)

	} else {
		if !dev.IsOnline() {
			GlobalDriverService.Online(dev.GetDeviceId())
		}
	}

	//数据处理，假如data[0]为设备标识、data[1]代表温度、data[2]代表湿度
	_, err = GlobalDriverService.PropertyReport(dev.GetDeviceId(), model.NewPropertyReport(false, map[string]model.PropertyData{
		"temperature": model.NewPropertyData(string(data[1])),
		"humidity":    model.NewPropertyData(string(data[2])),
	}))
	if err != nil {
		GlobalDriverService.GetLogger().Errorf("Device [%s] report data err:", deviceSn, err.Error())
	}
	//如果需要做tcp消息回复，请按照业务逻辑编写相应的resBuff
	resBuff := make([]byte, 3)
	resBuff[0] = 0x41               //包头
	resBuff[1] = byte(len(resBuff)) //长度
	resBuff[2] = 0xa3               //命令码
	return resBuff, nil
}
