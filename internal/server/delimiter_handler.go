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
	"bufio"
	"github.com/winc-link/hummingbird-tcp-driver/constant"
	"github.com/winc-link/hummingbird-tcp-driver/protocol"
	"io"
	"net"
)

// delimiterConnHandler 带有固定分割符数据处理方式
func delimiterConnHandler(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		buffer, err := reader.ReadSlice('\n')
		if err != nil {
			if err != io.EOF {
				GlobalDriverService.GetLogger().Errorf(err.Error())
			} else {
				break
			}
		}

		pack, err := protocol.Depack(buffer)
		if err != nil {
			GlobalDriverService.GetLogger().Error("depack error:", err.Error())
			continue
		}
		GlobalDriverService.GetLogger().Infof("depack info %s:", pack.MessageToString())

		if pack.Version != constant.Version {
			GlobalDriverService.GetLogger().Error("version not match")
			continue
		}

		switch pack.Operation {
		case constant.AuthOperation:
			err = deviceAuthHandler(pack.Data, conn)
			if err != nil {
				return
			}
		case constant.PropertyReportOperation:
			if tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()] == nil {
				GlobalDriverService.GetLogger().Errorf("remote addr [%s] unauthorized", conn.RemoteAddr().String())
				return
			}
			device := tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()]
			err = devicePropertyReportHandler(pack.Data, device.Id, conn, pack.SequenceID)
		case constant.EventReportOperation:
			if tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()] == nil {
				GlobalDriverService.GetLogger().Errorf("remote addr [%s] unauthorized", conn.RemoteAddr().String())
				return
			}
			device := tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()]
			err = deviceEventReportHandler(pack.Data, device.Id, conn, pack.SequenceID)
		case constant.PropertySetReplyOperation:
			if tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()] == nil {
				GlobalDriverService.GetLogger().Errorf("remote addr [%s] unauthorized", conn.RemoteAddr().String())
				return
			}
			device := tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()]
			err = devicePropertySetReplyHandler(pack.Data, device.Id, pack.SequenceID)
		case constant.PropertyGetReplyOperation:
			if tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()] == nil {
				GlobalDriverService.GetLogger().Errorf("remote addr [%s] unauthorized", conn.RemoteAddr().String())
				return
			}
			device := tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()]
			err = devicePropertyGetReplyHandler(pack.Data, device.Id, pack.SequenceID)
		case constant.ServiceExecuteReplyOperation:
			if tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()] == nil {
				GlobalDriverService.GetLogger().Errorf("remote addr [%s] unauthorized", conn.RemoteAddr().String())
				return
			}
			device := tcpServer.ConnAddrMapDev[conn.RemoteAddr().String()]
			err = deviceServiceExecuteReplyHandler(pack.Data, device.Id, pack.SequenceID)
		}

		if err != nil {
			GlobalDriverService.GetLogger().Errorf("handle data error: %s", err.Error())
		}

	}
}
