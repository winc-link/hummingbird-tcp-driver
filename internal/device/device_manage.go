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
	"errors"
	"fmt"
	"sync"
)

type DeviceManage struct {
	Devices map[string]Dev
	Lock    sync.RWMutex
}

var deviceManage = &DeviceManage{}

func init() {
	deviceManage = &DeviceManage{
		Devices: map[string]Dev{},
	}
}

func GetAllDevice() []Dev {
	var allDevice []Dev
	for _, dev := range deviceManage.Devices {
		allDevice = append(allDevice, dev)
	}
	return allDevice
}

func PutDevice(deviceSn string, device Dev) {
	deviceManage.Lock.Lock()
	defer deviceManage.Lock.Unlock()
	deviceManage.Devices[deviceSn] = device
}

func GetDevice(deviceSn string) (d Dev, err error) {
	deviceManage.Lock.RLock()
	defer deviceManage.Lock.RUnlock()
	var ok bool
	if d, ok = deviceManage.Devices[deviceSn]; !ok {
		err = errors.New(fmt.Sprint("get device err deviceSn:", deviceSn))
		return
	}
	return
}
