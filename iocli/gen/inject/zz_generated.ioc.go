//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package inject

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	"github.com/alibaba/ioc-golang/autowire/util"
)

func init() {
	objectGenCtxStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &objectGenCtx{}
		},
		ParamFactory: func() interface{} {
			var _ objectGenCtxParamInterface = &objectGenCtxParam{}
			return &objectGenCtxParam{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(objectGenCtxParamInterface)
			impl := i.(*objectGenCtx)
			return param.Init(impl)
		},
		DisableProxy: true,
	}
	normal.RegisterStructDescriptor(objectGenCtxStructDescriptor)
	importsListStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &importsList{}
		},
		ParamFactory: func() interface{} {
			var _ importsListParamInterface = &importsListParam{}
			return &importsListParam{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(importsListParamInterface)
			impl := i.(*importsList)
			return param.Init(impl)
		},
		DisableProxy: true,
	}
	normal.RegisterStructDescriptor(importsListStructDescriptor)
	copyMethodMakerStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &copyMethodMaker{}
		},
		ParamFactory: func() interface{} {
			var _ copyMethodMakerParamInterface = &copyMethodMakerParam{}
			return &copyMethodMakerParam{}
		},
		ConstructFunc: func(i interface{}, p interface{}) (interface{}, error) {
			param := p.(copyMethodMakerParamInterface)
			impl := i.(*copyMethodMaker)
			return param.Init(impl)
		},
		DisableProxy: true,
	}
	normal.RegisterStructDescriptor(copyMethodMakerStructDescriptor)
}

type objectGenCtxParamInterface interface {
	Init(impl *objectGenCtx) (*objectGenCtx, error)
}
type importsListParamInterface interface {
	Init(impl *importsList) (*importsList, error)
}
type copyMethodMakerParamInterface interface {
	Init(impl *copyMethodMaker) (*copyMethodMaker, error)
}

var _objectGenCtxSDID string

func GetobjectGenCtx(p *objectGenCtxParam) (*objectGenCtx, error) {
	if _objectGenCtxSDID == "" {
		_objectGenCtxSDID = util.GetSDIDByStructPtr(new(objectGenCtx))
	}
	i, err := normal.GetImpl(_objectGenCtxSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*objectGenCtx)
	return impl, nil
}

var _importsListSDID string

func GetimportsList(p *importsListParam) (*importsList, error) {
	if _importsListSDID == "" {
		_importsListSDID = util.GetSDIDByStructPtr(new(importsList))
	}
	i, err := normal.GetImpl(_importsListSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*importsList)
	return impl, nil
}

var _copyMethodMakerSDID string

func GetcopyMethodMaker(p *copyMethodMakerParam) (*copyMethodMaker, error) {
	if _copyMethodMakerSDID == "" {
		_copyMethodMakerSDID = util.GetSDIDByStructPtr(new(copyMethodMaker))
	}
	i, err := normal.GetImpl(_copyMethodMakerSDID, p)
	if err != nil {
		return nil, err
	}
	impl := i.(*copyMethodMaker)
	return impl, nil
}
