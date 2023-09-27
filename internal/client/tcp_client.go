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

package client

import (
	"github.com/winc-link/hummingbird-tcp-driver/internal/driver"
	"net"
)

type TcpClient struct {
}

func (t *TcpClient) Start(tdh TcpDataHandlers) {
	conn, err := net.Dial("tcp", "")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	var deviceSn string
	for {

		var buf [8]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			driver.GlobalDriverService.GetLogger().Errorf("Read from tcp cline failed,err:", err)
			closeConn(deviceSn)
			break
		}

		bytes, err := tdh(deviceSn, buf[:n])
		if err != nil {
			return
		}
		conn.Write(bytes)
	}

}
