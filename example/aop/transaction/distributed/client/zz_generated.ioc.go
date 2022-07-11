//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// Code generated by iocli, run 'iocli gen' to re-generate

package main

import (
	autowire "github.com/alibaba/ioc-golang/autowire"
	normal "github.com/alibaba/ioc-golang/autowire/normal"
	singleton "github.com/alibaba/ioc-golang/autowire/singleton"
	util "github.com/alibaba/ioc-golang/autowire/util"
)

func init() {
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &tradeService_{}
		},
	})
	tradeServiceStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &TradeService{}
		},
		TransactionMethodsMap: map[string]string{
			"DoTradeWithTxFinallyFailed":  "",
			"DoTradeWithTxAddMoneyFailed": "",
			"DoTradeWithTxSuccess":        "",
		},
	}
	singleton.RegisterStructDescriptor(tradeServiceStructDescriptor)
	normal.RegisterStructDescriptor(&autowire.StructDescriptor{
		Factory: func() interface{} {
			return &app_{}
		},
	})
	appStructDescriptor := &autowire.StructDescriptor{
		Factory: func() interface{} {
			return &App{}
		},
	}
	singleton.RegisterStructDescriptor(appStructDescriptor)
}

type tradeService_ struct {
	DoTradeWithTxFinallyFailed_  func(id1, id2, num int) error
	DoTradeWithTxAddMoneyFailed_ func(id1, id2, num int) error
	DoTradeWithTxSuccess_        func(id1, id2, num int) error
}

func (t *tradeService_) DoTradeWithTxFinallyFailed(id1, id2, num int) error {
	return t.DoTradeWithTxFinallyFailed_(id1, id2, num)
}

func (t *tradeService_) DoTradeWithTxAddMoneyFailed(id1, id2, num int) error {
	return t.DoTradeWithTxAddMoneyFailed_(id1, id2, num)
}

func (t *tradeService_) DoTradeWithTxSuccess(id1, id2, num int) error {
	return t.DoTradeWithTxSuccess_(id1, id2, num)
}

type app_ struct {
	Run_ func()
}

func (a *app_) Run() {
	a.Run_()
}

type TradeServiceIOCInterface interface {
	DoTradeWithTxFinallyFailed(id1, id2, num int) error
	DoTradeWithTxAddMoneyFailed(id1, id2, num int) error
	DoTradeWithTxSuccess(id1, id2, num int) error
}

type AppIOCInterface interface {
	Run()
}

var _tradeServiceSDID string

func GetTradeServiceSingleton() (*TradeService, error) {
	if _tradeServiceSDID == "" {
		_tradeServiceSDID = util.GetSDIDByStructPtr(new(TradeService))
	}
	i, err := singleton.GetImpl(_tradeServiceSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*TradeService)
	return impl, nil
}

func GetTradeServiceIOCInterfaceSingleton() (TradeServiceIOCInterface, error) {
	if _tradeServiceSDID == "" {
		_tradeServiceSDID = util.GetSDIDByStructPtr(new(TradeService))
	}
	i, err := singleton.GetImplWithProxy(_tradeServiceSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(TradeServiceIOCInterface)
	return impl, nil
}

var _appSDID string

func GetAppSingleton() (*App, error) {
	if _appSDID == "" {
		_appSDID = util.GetSDIDByStructPtr(new(App))
	}
	i, err := singleton.GetImpl(_appSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(*App)
	return impl, nil
}

func GetAppIOCInterfaceSingleton() (AppIOCInterface, error) {
	if _appSDID == "" {
		_appSDID = util.GetSDIDByStructPtr(new(App))
	}
	i, err := singleton.GetImplWithProxy(_appSDID, nil)
	if err != nil {
		return nil, err
	}
	impl := i.(AppIOCInterface)
	return impl, nil
}
