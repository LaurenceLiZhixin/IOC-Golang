//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli

package service1

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
}

type Impl1_ struct {
	Hello_ func(req string) string
}

func (i *Impl1_) Hello(req string) string {
	return i.Hello_(req)
}
func GetImpl1() (*Impl1, error) {
	i, err := singleton.GetImpl(util.GetSDIDByStructPtr(new(Impl1)))
	if err != nil {
		return nil, err
	}
	impl := i.(*Impl1)
	return impl, nil
}
