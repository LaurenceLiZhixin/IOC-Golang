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

package autowire

import (
	"fmt"
	"reflect"
	"sync"

	perrors "github.com/pkg/errors"

	"github.com/alibaba/ioc-golang/logger"

	"github.com/alibaba/ioc-golang/autowire/util"
)

type WrapperAutowire interface {
	Autowire

	ImplWithoutParam(sdID string, withProxy, force bool) (interface{}, error)
	ImplWithParam(sdID string, param interface{}, withProxy, force bool) (interface{}, error)

	implWithField(info *FieldInfo) (interface{}, error)
}

func getWrappedAutowire(autowire Autowire, allAutowires map[string]WrapperAutowire) WrapperAutowire {
	return &WrapperAutowireImpl{
		Autowire:           autowire,
		allAutowires:       allAutowires,
		singletonImpledMap: make(map[string]*SingletonCache),
	}
}

type SingletonCache struct {
	RawPtr   interface{}
	ProxyPtr interface{}
}

type WrapperAutowireImpl struct {
	Autowire
	singletonImpledMap     map[string]*SingletonCache
	singletonImpledMapLock sync.RWMutex
	allAutowires           map[string]WrapperAutowire
}

// ImplWithParam is used to get impled struct with param
func (w *WrapperAutowireImpl) ImplWithParam(sdID string, param interface{}, withProxy, force bool) (interface{}, error) {
	var rawPtr interface{}
	var err error
	// 1. check singleton
	w.singletonImpledMapLock.RLock()
	if singletonCache, ok := w.singletonImpledMap[sdID]; !force && w.Autowire.IsSingleton() && ok {
		if !withProxy {
			// singletonCache.RawPtr must be cached before
			w.singletonImpledMapLock.RUnlock()
			return singletonCache.RawPtr, nil
		}
		// want proxy ptr, check if cached it
		if singletonCache.ProxyPtr != nil {
			// proxy ptr found
			w.singletonImpledMapLock.RUnlock()
			return singletonCache.ProxyPtr, nil
		}

		// proxy ptr not found, set raw ptr to rawPtr to re-used
		rawPtr = singletonCache.RawPtr
	}
	w.singletonImpledMapLock.RUnlock()

	// 2. factory
	if rawPtr == nil {
		rawPtr, err = w.Autowire.Factory(sdID)
		if err != nil {
			return nil, err
		}
	}

	if w.Autowire.InjectPosition() == AfterFactoryCalled {
		if err := w.inject(rawPtr, sdID); err != nil {
			return nil, err
		}
	}

	// 3. construct field
	rawPtr, err = w.Autowire.Construct(sdID, rawPtr, param)
	if err != nil {
		errMsg := fmt.Sprintf("Construct struct %s failed with error = %s, param = %s", sdID, err.Error(), param)
		logger.Red(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	if rawPtr == nil {
		errMsg := fmt.Sprintf("Construct struct %s failed, constructed ptr is nil, param = %s", sdID, param)
		logger.Red(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if w.Autowire.InjectPosition() == AfterConstructorCalled {
		if err := w.inject(rawPtr, sdID); err != nil {
			return nil, err
		}
	}

	// 4. try to wrap proxy
	var finalPtr = rawPtr
	var proxyPtr interface{}
	if withProxy {
		// if field is interface, try to inject proxy wrapped pointer
		finalPtr = GetProxyFunction()(rawPtr)
		proxyPtr = finalPtr
	}

	// 5. record singleton ptr
	if w.Autowire.IsSingleton() && !force {
		w.singletonImpledMapLock.Lock()
		w.singletonImpledMap[sdID] = &SingletonCache{
			RawPtr:   rawPtr,
			ProxyPtr: proxyPtr,
		}
		w.singletonImpledMapLock.Unlock()
	}
	return finalPtr, nil
}

// ImplWithoutParam is used to create param from field without param
func (w *WrapperAutowireImpl) ImplWithoutParam(sdID string, withProxy, force bool) (interface{}, error) {
	param, err := w.ParseParam(sdID, nil)
	if err != nil {
		if w.Autowire.IsSingleton() {
			// FIXME: ignore parse param error, because of singleton with empty param also try to find property from config file
			logger.Blue("[Wrapper Autowire] Parse param from config file with sdid %s failed, error: %s, continue with nil param.", sdID, err)
			return w.ImplWithParam(sdID, param, withProxy, force)
		} else {
			return nil, err
		}
	}
	return w.ImplWithParam(sdID, param, withProxy, force)
}

// ImplWithField is used to create param from field and call ImplWithParam
func (w *WrapperAutowireImpl) implWithField(fi *FieldInfo) (interface{}, error) {
	sdID, err := w.ParseSDID(fi)
	if err != nil {
		logger.Red("[Wrapper Autowire] Parse sdid from field %+v failed, error is %s", fi, err)
		return nil, err
	}
	sd := GetStructDescriptor(sdID)
	implWithProxy := fi.FieldReflectValue.Kind() == reflect.Interface
	if implWithProxy {
		if sd == nil {
			err = fmt.Errorf("[Wrapper Autowire] sdID %s is invalid when injecting %s type %s, please check", sdID, fi.FieldName, fi.FieldType)
			logger.Red(err.Error())
			return nil, err
		}
		implWithProxy = !sd.DisableProxy
	}
	param, err := w.ParseParam(sdID, fi)
	if err != nil {
		if w.Autowire.IsSingleton() {
			// ignore parse param error, because of singleton with empty param also try to find property from config file
			return w.ImplWithParam(sdID, param, implWithProxy, false)
		} else {
			return nil, err
		}
	}
	return w.ImplWithParam(sdID, param, implWithProxy, false)
}

// inject do tag autowire and monkey inject
func (w *WrapperAutowireImpl) inject(impledPtr interface{}, sdId string) error {
	sd := w.Autowire.GetAllStructDescriptors()[sdId]

	// 1. reflect
	valueOf := reflect.ValueOf(impledPtr)
	if valueOf.Kind() != reflect.Interface && valueOf.Kind() != reflect.Ptr {
		// not struct pointer, no needs to inject fields, just return
		return nil
	}
	valueOfElem := valueOf.Elem()
	typeOf := valueOfElem.Type()
	if typeOf.Kind() != reflect.Struct {
		// element not struct, no needs to inject fields, just return
		return nil
	}

	// deal with struct
	// 2. tag inject
	numField := valueOfElem.NumField()
	for i := 0; i < numField; i++ {
		field := typeOf.Field(i)
		var subImpledPtr interface{}
		var subService reflect.Value
		tagKey := ""
		tagValue := ""
		for _, aw := range w.allAutowires {
			if val, ok := field.Tag.Lookup(aw.TagKey()); ok {
				// check field
				subService = valueOfElem.Field(i)
				tagKey = aw.TagKey()
				tagValue = val
				if !(subService.IsValid() && subService.CanSet()) {
					errMsg := fmt.Sprintf("Failed to autowire struct %s's impl %s service. It's field type %s with tag '%s:\"%s\"', please check if the field name is exported",
						sd.ID(), util.GetStructName(impledPtr), field.Type.Name(), tagKey, tagValue)
					logger.Red("[Autowire Wrapper] Inject field failed with error: %s", errMsg)
					return perrors.New(errMsg)
				}

				fieldType := buildFiledTypeFullName(field.Type)
				fieldInfo := &FieldInfo{
					FieldName:         field.Name,
					FieldType:         fieldType,
					TagKey:            aw.TagKey(),
					TagValue:          val,
					FieldReflectType:  field.Type,
					FieldReflectValue: subService,
				}
				// create param from field info
				var err error

				subImpledPtr, err = aw.implWithField(fieldInfo)
				if err != nil {
					return err
				}
				break // only one tag is supported
			}
		}
		if tagKey == "" && tagValue == "" {
			continue
		}
		// set field
		subService.Set(reflect.ValueOf(subImpledPtr))
	}
	return nil
}

// todo can we downward the parsing of field to autowire impl but not autowire core?
func buildFiledTypeFullName(fieldType reflect.Type) string {
	// todo find unsupported type and log warning, like 'struct' field
	if util.IsPointerField(fieldType) || util.IsSliceField(fieldType) {
		return util.GetSDIDByReflectType(fieldType.Elem())
	}
	// interface field
	return util.GetSDIDByReflectType(fieldType)
}
