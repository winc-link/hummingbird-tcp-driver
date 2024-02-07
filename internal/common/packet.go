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

package common

import (
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
)

type DataPacket string

func (t DataPacket) IsHeartbeatPacket() bool {
	return t == "700a"
}

func (t DataPacket) IsDataPacket() bool {
	return len(t) == len("414c4d0039a0764eafbcf2004257432942032042a0764eafbcf2cd29a05835a1001de200000599a313868415e50000266186010400005b7c000a")
}

func (t DataPacket) Verify() bool {
	return t.IsDataPacket() || t.IsDataPacket()
}

// DatePacketPrefix 数据前缀
func (t DataPacket) DatePacketPrefix() (data string) {
	return string(t[0:6])
}

// TotalLength 数据总长度
func (t DataPacket) TotalLength() (data string) {
	return string(t[6:10])
}

// DeviceSn 设备ID
func (t DataPacket) DeviceSn() (data string) {
	return string(t[10:22])
}

// DeviceType 设备类型
func (t DataPacket) DeviceType() (data string) {
	return string(t[22:26])
}

// DataHeader 数据头
func (t DataPacket) DataHeader() (data string) {
	return string(t[26:30])
}

// DataLength 数据总长度
func (t DataPacket) DataLength() (data string) {
	return string(t[30:32])
}

// CommandType 命令类型
func (t DataPacket) CommandType() (data string) {
	return string(t[30:32])
}

// RSSI RSSI
func (t DataPacket) RSSI() (data string) {
	return string(t[52:54])
}

// FlagId 标识id
func (t DataPacket) FlagId() (data string) {
	return string(t[54:56])
}

// Voltage 电压
func (t DataPacket) Voltage() (data float64) {
	return Decimal(float64(Hex2Dec(string(t[58:62]))) / 100.00)
}

// Current 电流
func (t DataPacket) Current() (data float64) {
	//fmt.Println(string(t[62:68]))
	return Decimal(float64(Hex2Dec(string(t[64:68]))) / 100.00)
}

// Power 功率
func (t DataPacket) Power() (data float64) {
	//fmt.Println(string(t[68:78]))
	return Decimal(float64(Hex2Dec(string(t[70:78]))) / 100.00)
}

// Frequency 电压频率
func (t DataPacket) Frequency() (data float64) {
	//fmt.Println(string(t[78:84]))
	return Decimal(float64(Hex2Dec(string(t[80:84]))) / 100.00)
}

// PowerFactor 功率因数
func (t DataPacket) PowerFactor() (data float64) {
	//fmt.Println(string(t[84:88]))
	return Decimal(float64(Hex2Dec(string(t[86:88]))) / 100.00)
}

// PowerConsumption 累计用电量
func (t DataPacket) PowerConsumption() (data float64) {
	return Decimal(float64(Hex2Dec(string(t[90:98]))) / 100.00)
}

// MStatus 开关状态
func (t DataPacket) MStatus() (data int) {
	return Hex2Dec(string(t[100:102]))
}

func Hex2Dec(val string) int {
	n, err := strconv.ParseUint(val, 16, 32)
	if err != nil {
		fmt.Println(err)
	}
	return int(n)
}

func Decimal(value float64) float64 {
	v1, _ := decimal.NewFromFloat(value).Round(2).Float64()
	return v1
}
