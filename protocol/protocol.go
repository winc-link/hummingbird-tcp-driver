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

package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/winc-link/hummingbird-tcp-driver/constant"
	"github.com/winc-link/hummingbird-tcp-driver/dtos"
)

// 使用了length field based frame decoder
const (
	PackageLengthBytes = 4
	HeaderLengthBytes  = 2
	VersionBytes       = 2
	OperationBytes     = 2
	SequenceIDBytes    = 36

	HeaderLength = PackageLengthBytes + HeaderLengthBytes + VersionBytes + OperationBytes + SequenceIDBytes
)

// Depack 解码器
func Depack(buffer []byte) (dtos.Packet, error) {
	length := len(buffer)

	var pack dtos.Packet
	var err error

	if length == 0 {
		return pack, errors.New("invalid Parameter")
	}
	for i := 0; i < length; i++ {
		if length < i+HeaderLength {
			return pack, errors.New("data length less then header length")
		}
		messageLength := ByteToInt(buffer[i : i+PackageLengthBytes])
		if length < i+HeaderLength+messageLength {
			return pack, errors.New("data length less then header length and message length")
		}
		pack.PackageLength = messageLength

		site := i + PackageLengthBytes
		headerLength := ByteToInt16(buffer[site : site+HeaderLengthBytes])
		pack.HeaderLength = headerLength
		site += HeaderLengthBytes

		protocolVersion := ByteToInt16(buffer[site : site+VersionBytes])
		pack.Version = protocolVersion
		site += VersionBytes

		operation := ByteToInt16(buffer[site : site+OperationBytes])
		pack.Operation = operation

		site += OperationBytes

		//sequenceID := ByteToInt(buffer[site : site+SequenceIDBytes])
		sequenceID := string(buffer[site : site+SequenceIDBytes])
		pack.SequenceID = sequenceID
		site += SequenceIDBytes
		//fmt.Printf("packageLength: %d, headerLength: %d , protocolVersion: %d, operation: %d, sequenceID: %d \n", messageLength, headerLength, protocolVersion, operation, sequenceID)
		data := buffer[i+HeaderLength : i+HeaderLength+messageLength]
		pack.Data = data
		break
	}
	return pack, err
}

// ByteToInt 字节转换成整形
func ByteToInt(n []byte) int {
	bytesbuffer := bytes.NewBuffer(n)
	var x int32
	binary.Read(bytesbuffer, binary.BigEndian, &x)

	return int(x)
}

func ByteToInt16(n []byte) int {
	bytebuffer := bytes.NewBuffer(n)
	var x int16
	binary.Read(bytebuffer, binary.BigEndian, &x)

	return int(x)
}

// Packet 封包
func Packet(message []byte, sequenceID string, operation int) []byte {
	b := append(Int32ToBytes(len(message)), Int16ToBytes(len(Int16ToBytes(constant.Version))+len(Int16ToBytes(operation))+len([]byte(sequenceID)))...) //包长度4 //消息头 2
	b = append(b, Int16ToBytes(constant.Version)...)                                                                                                   //版本 2
	b = append(b, Int16ToBytes(operation)...)                                                                                                          //操作 2
	b = append(b, []byte(sequenceID)...)                                                                                                               // 唯一消息 4
	b = append(b, message...)
	return b
}

// Int32ToBytes 整数转换成字节
func Int32ToBytes(n int) []byte {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

// Int16ToBytes 整数转换成字节
func Int16ToBytes(n int) []byte {
	x := int16(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}
