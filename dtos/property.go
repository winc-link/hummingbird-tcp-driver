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

package dtos

import (
	"encoding/json"
	"github.com/winc-link/hummingbird-sdk-go/model"
)

// PropertyPost 属性上报
type PropertyPost struct {
	Sys    Sys                           `json:"sys"`
	Params map[string]model.PropertyData `json:"params"`
}

func BytesToPropertyStruct(b []byte) (PropertyPost, error) {
	var property PropertyPost
	var err error
	err = json.Unmarshal(b, &property)
	return property, err
}

// PropertySetReply 设置设备属性响应
type PropertySetReply struct {
	Params model.PropertySetResponseData `json:"params"`
}

func BytesToPropertySetReplyStruct(b []byte) (PropertySetReply, error) {
	var propertySet PropertySetReply
	var err error
	err = json.Unmarshal(b, &propertySet)
	return propertySet, err
}

// PropertyGetReply 设备属性查询
type PropertyGetReply struct {
	Params []model.PropertyGetResponseData `json:"params"`
}

func BytesToPropertyGetReplyStruct(b []byte) (PropertyGetReply, error) {
	var propertyGet PropertyGetReply
	var err error
	err = json.Unmarshal(b, &propertyGet)
	return propertyGet, err
}

type ServiceExecuteReply struct {
	Params model.ServiceDataOut `json:"params"`
}

func BytesToServiceExecuteReplyStruct(b []byte) (ServiceExecuteReply, error) {
	var serviceExecuteReply ServiceExecuteReply
	var err error
	err = json.Unmarshal(b, &serviceExecuteReply)
	return serviceExecuteReply, err
}

type PropertySet struct {
	Params map[string]interface{} `json:"params"`
}
