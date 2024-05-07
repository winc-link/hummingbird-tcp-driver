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

package constant

const (
	AuthOperation                = iota + 1 //鉴权
	PropertyReportOperation                 //属性上报
	PropertyReportReplyOperation            //属性上报响应
	EventReportOperation                    //事件上报
	EventReportReplyOperation               //事件上报响应
	PropertySet                             //属性设置
	PropertySetReplyOperation               //属性设置响应
	PropertyGet                             //属性获取
	PropertyGetReplyOperation               //属性获取响应
	ServiceExecute                          //服务调用
	ServiceExecuteReplyOperation            //服务调用响应
)
