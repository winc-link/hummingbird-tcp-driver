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
	"github.com/winc-link/hummingbird-tcp-driver/config"
	"net"
)

type TcpDataHandlers func(deviceSn string, data []byte) (retBuff []byte, err error)

// serverConnHandler 用户可以根据项目需要自行修改此方法的业务逻辑！
func serverConnHandler(conn net.Conn) {
	switch config.GetConfig().TcpUnpackRule.RuleName {
	case config.RuleDelimiter:
		delimiterConnHandler(conn)
	case config.RuleFixedLength:
		fixLengthConnHandler(conn)
	default:
		delimiterConnHandler(conn)
	}
}
