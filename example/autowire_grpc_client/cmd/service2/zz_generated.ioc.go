//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli

package service2

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &Impl1_{}
		},
	})
	singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &Impl1{}
		},
	})
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &Impl2_{}
		},
	})
	singleton.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &Impl2{}
		},
	})
}

type Impl1_ struct {
	Hello_ func(name string) string
}

func (i *Impl1_) Hello(name string) string {
	return i.Hello_(name)
}

type Impl2_ struct {
	Hello_ func(name string) string
}

func (i *Impl2_) Hello(name string) string {
	return i.Hello_(name)
}
func GetImpl1() (*Impl1, error) {
	i, err := singleton.GetImpl(util.GetSDIDByStructPtr(new(Impl1)))
	if err != nil {
		return nil, err
	}
	impl := i.(*Impl1)
	return impl, nil
}
func GetImpl2() (*Impl2, error) {
	i, err := singleton.GetImpl(util.GetSDIDByStructPtr(new(Impl2)))
	if err != nil {
		return nil, err
	}
	impl := i.(*Impl2)
	return impl, nil
}
