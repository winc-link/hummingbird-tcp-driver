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

package device

type Dev interface {
	GetDeviceId() string
	GetDeviceSn() string
	IsOnline() bool
}

type Device struct {
	ProductId string
	DeviceId  string
	DeviceSn  string
	isOnline  bool
}

func (d *Device) GetDeviceId() string {
	return d.DeviceId
}

func (d *Device) GetDeviceSn() string {
	return d.DeviceSn
}

func (d *Device) IsOnline() bool {
	return d.isOnline
}

func NewDevice(deviceId, deviceSn, ProductId string, isOnline bool) Dev {
	return &Device{
		DeviceId:  deviceId,
		DeviceSn:  deviceSn,
		ProductId: ProductId,
		isOnline:  isOnline,
	}
}
