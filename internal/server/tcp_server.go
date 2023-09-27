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
	"github.com/winc-link/hummingbird-sdk-go/service"
	"net"
	"sync"
)

var GlobalDriverService *service.DriverService

type TcpServer struct {
	ClientCons map[string]*Connect
	Lock       sync.RWMutex
}

type Connect struct {
	Conn       net.Conn
	DeviceInfo interface{}
}

var tcpServer = &TcpServer{}

func init() {
	tcpServer = &TcpServer{
		ClientCons: map[string]*Connect{},
	}
}

// Start 启动tcp服务器
func (t *TcpServer) Start(sd *service.DriverService, tdh TcpDataHandlers) {
	GlobalDriverService = sd
	server, err := net.Listen("tcp", "")
	if err != nil {
		panic(err)
	}
	defer server.Close()
	for {
		conn, err := server.Accept()
		if err != nil {
			panic(err)
		}
		go serverConnHandler(conn, tdh)
	}
}

func GetTcpServer() *TcpServer {
	return tcpServer
}

func (t *TcpServer) GetConnectByDeviceSn(deviceSn string) *Connect {
	return t.ClientCons[deviceSn]
}

func (t *TcpServer) DeleteClientByDeviceId(deviceSn string) {
	t.Lock.Lock()
	defer t.Lock.Unlock()
	if _, ok := t.ClientCons[deviceSn]; ok {
		delete(t.ClientCons, deviceSn)
	}
}
