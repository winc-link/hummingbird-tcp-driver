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

import (
	"github.com/winc-link/hummingbird-sdk-go/model"
	"github.com/winc-link/hummingbird-tcp-driver/internal/driver"
)

type Dev interface {
	GetDeviceId() string
	GetDeviceSn() string
	Online() error
	Offline() error
	IsOnline() bool
	PropertyReport(report model.PropertyReport) error
	EventReport(report model.EventReport) error
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

func (d *Device) Online() error {
	return driver.GlobalDriverService.Online(d.DeviceId)
}

func (d *Device) Offline() error {
	return driver.GlobalDriverService.Offline(d.DeviceId)
}

func (d *Device) IsOnline() bool {
	return d.isOnline
}

func (d *Device) PropertyReport(report model.PropertyReport) error {
	_, err := driver.GlobalDriverService.PropertyReport(d.DeviceId, report)
	if err != nil {
		driver.GlobalDriverService.GetLogger().Errorf("device [%s] report property error:%s ", d.DeviceId, err.Error())
	}
	return err
}

func (d *Device) EventReport(report model.EventReport) error {
	_, err := driver.GlobalDriverService.EventReport(d.DeviceId, report)
	if err != nil {
		driver.GlobalDriverService.GetLogger().Errorf("device [%s] report event error:%s ", d.DeviceId, err.Error())
	}
	return err
}

func NewDevice(deviceId, deviceSn, ProductId string, isOnline bool) Dev {
	return &Device{
		DeviceId:  deviceId,
		DeviceSn:  deviceSn,
		ProductId: ProductId,
		isOnline:  isOnline,
	}
}
