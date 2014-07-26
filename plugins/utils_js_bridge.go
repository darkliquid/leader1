package plugins

import (
	"github.com/darkliquid/leader1/utils"
	"github.com/robertkrimen/otto"
)

type pmUtilsJSBridge struct {
	GetPage, ExtractURL, ExtractTitle, Sleep, GetShoutcastStats,
	LikeTrack, HateTrack, Request func(call otto.FunctionCall) otto.Value
}

func (pm *PluginManager) InitUtilsJSBridge() {
	bridge := &pmUtilsJSBridge{
		GetPage: func(call otto.FunctionCall) otto.Value {
			var err error
			switch {
			case len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString():
				if page, err := utils.GetPage(call.Argument(0).String()); err == nil {
					if val, err := pm.js.ToValue(page); err == nil {
						return val
					}
				}
				pm.log.Printf("[UTILS] GetPage errored: %s\n", err)
			case len(call.ArgumentList) == 3 && call.ArgumentList[0].IsString() && call.ArgumentList[1].IsString() && call.ArgumentList[2].IsString():
				if page, err := utils.GetPageWithAuth(call.Argument(0).String(), call.Argument(1).String(), call.Argument(2).String()); err == nil {
					if val, err := pm.js.ToValue(page); err == nil {
						return val
					}
				}
				pm.log.Printf("[UTILS] GetPageWithAuth errored: %s\n", err)
			}
			return otto.FalseValue()
		},
		ExtractURL: func(call otto.FunctionCall) otto.Value {
			var err error
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				if url, err := utils.ExtractURL(call.Argument(0).String()); err == nil {
					if val, err := pm.js.ToValue(url); err == nil {
						return val
					}
				}
				pm.log.Printf("[UTILS] ExtractURL errored: %s\n", err)
			}
			return otto.FalseValue()
		},
		ExtractTitle: func(call otto.FunctionCall) otto.Value {
			var err error
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				if title, err := utils.ExtractTitle(call.Argument(0).String()); err == nil {
					if val, err := pm.js.ToValue(title); err == nil {
						return val
					}
				}
				pm.log.Printf("[UTILS] ExtractTitle errored: %s\n", err)
			}
			return otto.FalseValue()
		},
		Sleep: func(call otto.FunctionCall) otto.Value {
			var err error
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsNumber() {
				if i, err := call.Argument(0).ToInteger(); err == nil {
					utils.Sleep(i)
					return otto.TrueValue()
				}
				pm.log.Printf("[UTILS] Sleep errored: %s\n", err)
			}
			return otto.FalseValue()
		},
		GetShoutcastStats: func(call otto.FunctionCall) otto.Value {
			var err error
			if len(call.ArgumentList) == 0 {
				if stats, err := utils.GetShoutcastStats(pm.cfg); err == nil {
					if val, err := pm.js.ToValue(stats); err == nil {
						return val
					}
				}
				pm.log.Printf("[UTILS] GetShoutcastStats errored: %s\n", err)
			}
			return otto.FalseValue()
		},
		LikeTrack: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 2 && call.ArgumentList[0].IsString() && call.ArgumentList[1].IsString() {
				ok, _ := utils.LikeTrack(call.Argument(0).String(), call.Argument(1).String())
				if ok {
					return otto.TrueValue()
				}
				pm.log.Println("[UTILS] LikeTrack failed")
			}
			return otto.FalseValue()
		},
		HateTrack: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 2 && call.ArgumentList[0].IsString() && call.ArgumentList[1].IsString() {
				ok, _ := utils.HateTrack(call.Argument(0).String(), call.Argument(1).String())
				if ok {
					return otto.TrueValue()
				}
				pm.log.Println("[UTILS] HateTrack failed")
			}
			return otto.FalseValue()
		},
		Request: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 2 && call.ArgumentList[0].IsString() && call.ArgumentList[1].IsString() {
				ok, _ := utils.Request(call.Argument(0).String(), call.Argument(1).String())
				if ok {
					return otto.TrueValue()
				}
				pm.log.Println("[UTILS] Request failed")
			}
			return otto.FalseValue()
		},
	}
	pm.js.Set("UTILS", bridge)
}
