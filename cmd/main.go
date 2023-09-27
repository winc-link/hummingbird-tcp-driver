/*******************************************************************************
 * Copyright 2023 Hummingbird.
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

package main

import (
	"context"
	"github.com/winc-link/hummingbird-tcp-driver/internal/driver"
	"os"
	"os/signal"
	"syscall"

	"github.com/winc-link/hummingbird-sdk-go/commons"
	"github.com/winc-link/hummingbird-sdk-go/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	driverService := service.NewDriverService("hummingbird-tcp-driver", commons.HummingbirdIot)
	tcpDriver := driver.NewTcpProtocolDriver(ctx, driverService)
	go func() {
		if err := driverService.Start(tcpDriver); err != nil {
			driverService.GetLogger().Error("driver service start error: %s", err)
			return
		}
	}()
	waitForSignal(cancel)
}

func waitForSignal(cancel context.CancelFunc) os.Signal {
	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	s := <-signalChan
	cancel()
	signal.Stop(signalChan)
	return s
}