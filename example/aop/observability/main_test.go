/*
 * Copyright (c) 2022, Alibaba Group;
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"strings"
	"testing"
	"time"

	"github.com/alibaba/ioc-golang/aop/common"
	"github.com/alibaba/ioc-golang/config"

	"github.com/alibaba/ioc-golang"

	"github.com/stretchr/testify/assert"

	"github.com/alibaba/ioc-golang/test/iocli_command"
)

func TestObservability(t *testing.T) {
	testAppName := "testAppName"
	assert.Nil(t, ioc.Load(config.AddProperty(common.IOCGolangAOPConfigPrefix+".app-name", testAppName)))
	app, err := GetAppIOCInterfaceSingleton()
	assert.Nil(t, err)
	go func() {
		app.Run()
	}()
	time.Sleep(time.Second * 1)
	output, err := iocli_command.Run([]string{"list"}, time.Second)
	assert.Nil(t, err)
	assert.Equal(t, `appName: `+testAppName+`
github.com/alibaba/ioc-golang/example/aop/observability.App
[Run]

github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1
[GetHelloString]

github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl2
[GetHelloString]

`, output)

	output, err = iocli_command.Run([]string{"monitor", "-i", "3"}, time.Second*4)
	assert.Nil(t, err)
	t.Log(output)
	assert.True(t, strings.Contains(output, `github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1.GetHelloString()
Total: 1, Success: 1, Fail: 0, AvgRT: `))
	assert.True(t, strings.Contains(output, `us, FailRate: 0.00%
github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl2.GetHelloString()
Total: 1, Success: 1, Fail: 0, AvgRT: `))

	output, err = iocli_command.Run([]string{"watch", "github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1", "GetHelloString"}, time.Second*6)
	assert.Nil(t, err)
	assert.True(t, strings.Contains(output, `========== On Call ==========
github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1.GetHelloString()
Param 1: (string) (len=8) "laurence"

========== On Response ==========
github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1.GetHelloString()
Response 1: (string) (len=36) "This is ServiceImpl2, hello laurence"`))

	output, err = iocli_command.Run([]string{"monitor", "-i", "3"}, time.Second*4)
	assert.Nil(t, err)
	t.Log(output)
	assert.True(t, strings.Contains(output, `github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1.GetHelloString()
Total: 1, Success: 1, Fail: 0, AvgRT: `))
	assert.True(t, strings.Contains(output, `us, FailRate: 0.00%
github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl2.GetHelloString()
Total: 1, Success: 1, Fail: 0, AvgRT: `))

	output, err = iocli_command.Run([]string{"trace", "github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl1", "GetHelloString"}, time.Second*6)
	assert.Nil(t, err)
	t.Log(output)
	assert.True(t, strings.Contains(output, `OperationName: github.com/alibaba/ioc-golang/example/aop/observability.(*serviceImpl2_).GetHelloString, StartTime: `))
	assert.True(t, strings.Contains(output, `OperationName: github.com/alibaba/ioc-golang/example/aop/observability.(*serviceImpl1_).GetHelloString, StartTime: `))
	assert.True(t, strings.Contains(output, `ReferenceSpans: [{TraceID:`))

	output, err = iocli_command.Run([]string{"trace", "github.com/alibaba/ioc-golang/example/aop/observability.ServiceImpl2", "GetHelloString"}, time.Second*6)
	assert.Nil(t, err)
	t.Log(output)
	assert.True(t, strings.Contains(output, `OperationName: github.com/alibaba/ioc-golang/example/aop/observability.(*serviceImpl2_).GetHelloString, StartTime: `))
	assert.True(t, !strings.Contains(output, `OperationName: github.com/alibaba/ioc-golang/example/aop/observability.(*serviceImpl1_).GetHelloString, StartTime: `))
	assert.True(t, strings.Contains(output, `ReferenceSpans: [{TraceID:`))

}
