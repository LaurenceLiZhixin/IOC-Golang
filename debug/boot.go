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

package debug

import (
	"context"
	"reflect"

	"github.com/fatih/color"
	"github.com/glory-go/monkey"

	"github.com/alibaba/ioc-golang/autowire"
	"github.com/alibaba/ioc-golang/config"
	"github.com/alibaba/ioc-golang/debug/common"
	"github.com/alibaba/ioc-golang/debug/interceptor"
)

const (
	defaultDebugPort = "1999"
)

func init() {
	autowire.RegisterMonkeyFunction(implMonkey)
}

var guardMap = make(map[string]*common.DebugMetadata)

func Load() error {
	// start debug port if enabled
	bootConfig := &Config{}
	if err := config.LoadConfigByPrefix("debug", bootConfig); err != nil || !bootConfig.Enable {
		color.Blue("[Debug] Debug mod is not enabled")
		return nil
	}
	if bootConfig.Port == "" {
		color.Blue("[Debug] Debug port is set to default :%s", defaultDebugPort)
		bootConfig.Port = defaultDebugPort
	}
	if err := interceptor.StartDebugServer(bootConfig.Port, guardMap); err != nil {
		color.Red("[Debug] Start debug server error = %s", err)
		return err
	}
	return nil
}

func implMonkey(servicePtr interface{}, tempInterfaceId string) {
	if _, ok := guardMap[tempInterfaceId]; !ok {
		guardMap[tempInterfaceId] = &common.DebugMetadata{
			GuardMap: map[string]*common.GuardInfo{},
		}
	}
	valueOf := reflect.ValueOf(servicePtr)
	typeOf := reflect.TypeOf(servicePtr)
	valueOfElem := valueOf.Elem()
	typeOfElem := valueOfElem.Type()
	if typeOfElem.Kind() != reflect.Struct {
		panic("invalid struct ptr")
	}

	numField := valueOf.NumMethod()
	for i := 0; i < numField; i++ {
		methodType := typeOf.Method(i)
		if _, ok := guardMap[tempInterfaceId].GuardMap[methodType.Name]; !ok {
			guardMap[tempInterfaceId].GuardMap[methodType.Name] = &common.GuardInfo{}
		}
		if guardMap[tempInterfaceId].GuardMap[methodType.Name].Guard == nil {
			// each method of one type should only injected once
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(servicePtr), methodType.Name,
				reflect.MakeFunc(methodType.Type, makeCallProxy(tempInterfaceId, methodType.Name, methodType.Type.IsVariadic())).Interface(),
			)
			guardMap[tempInterfaceId].GuardMap[methodType.Name].Guard = guard
		}
		continue
	}
}

func makeCallProxy(tempInterfaceId, methodName string, isVariadic bool) func(in []reflect.Value) []reflect.Value {
	return func(in []reflect.Value) []reflect.Value {
		guardMap[tempInterfaceId].GuardMap[methodName].Lock.Lock()
		guardMap[tempInterfaceId].GuardMap[methodName].Guard.Unpatch()
		defer func() {
			r := recover()
			if r != nil {
				color.Red("[Function Proxy] Recover panic %s", r)
			}
			guardMap[tempInterfaceId].GuardMap[methodName].Guard.Restore()
			guardMap[tempInterfaceId].GuardMap[methodName].Lock.Unlock()
		}()
		ctx := common.NewInterceptorContext(context.Background(), tempInterfaceId, methodName, true)
		var err error
		for _, i := range registeredInterceptors {
			in, err = i.Invoke(ctx, in)
			if err != nil {
				color.Red("[Interceptor] Interceptor %s invoke failed, error is %s, interceptor context is %s", i.Name(), err, ctx)
			}
		}

		if isVariadic {
			varParam := in[len(in)-1]
			in = in[:len(in)-1]
			for j, l := 0, varParam.Len(); j < l; j++ {
				in = append(in, varParam.Index(j))
			}
		}

		ctx = common.NewInterceptorContext(ctx.Context, tempInterfaceId, methodName, false)
		out := in[0].MethodByName(methodName).Call(in[1:])
		for _, i := range registeredInterceptors {
			out, err = i.Invoke(ctx, out)
			if err != nil {
				color.Red("[Interceptor] Interceptor %s invoke failed, error is %s, interceptor context is %s", i.Name(), err, ctx)
			}
		}
		return out
	}
}
