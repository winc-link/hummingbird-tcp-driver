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

package config

import (
	"encoding/json"
	"github.com/winc-link/hummingbird-sdk-go/service"
)

const (
	RuleNone        = "none"        //不处理
	RuleDelimiter   = "delimiter"   //分隔符
	RuleFixedLength = "fixedLength" //固定长度
	RuleFieldLength = "fieldLength" //长度字段
)

var baseConfig *BaseConfig

type BaseConfig struct {
	//用户自行定义结构体中信息。
	TcpUnpackRule struct {
		RuleName string `json:"rule_name"`
		BytesLen int    `json:"bytes_len"`
	} `json:"tcp_unpack_rule"`
}

func InitConfig(sd *service.DriverService) {
	customParam := sd.GetCustomParam()
	baseConfig = &BaseConfig{}
	if customParam != "" {
		err := json.Unmarshal([]byte(customParam), &baseConfig)
		if err != nil {
			sd.GetLogger().Error(err)
		}
	}
}

func GetConfig() *BaseConfig {
	return baseConfig
}
