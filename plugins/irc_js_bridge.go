package plugins

import (
	"github.com/robertkrimen/otto"
)

type pmIRCJSBridge struct {
	Nick, GetNick, SendRaw, Privmsg, Notice, Part, Join, Who, Whois, Mode func(call otto.FunctionCall) otto.Value
}

func (pm *PluginManager) InitIRCJSBridge() {
	bridge := &pmIRCJSBridge{
		Nick: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				pm.conn.Nick(call.Argument(0).String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
		GetNick: func(call otto.FunctionCall) otto.Value {
			val, err := otto.ToValue(pm.conn.GetNick())
			if err != nil {
				return otto.UndefinedValue()
			}
			return val
		},
		SendRaw: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				pm.conn.SendRaw(call.Argument(0).String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
		Privmsg: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 2 && call.ArgumentList[0].IsString() && call.ArgumentList[1].IsString() {
				pm.conn.Privmsg(call.Argument(0).String(), call.Argument(1).String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
		Notice: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 2 && call.ArgumentList[0].IsString() && call.ArgumentList[1].IsString() {
				pm.conn.Notice(call.Argument(0).String(), call.Argument(1).String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
		Part: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				pm.conn.Part(call.Argument(0).String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
		Join: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				pm.conn.Join(call.Argument(0).String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
		Who: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				pm.conn.Who(call.Argument(0).String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
		Whois: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				pm.conn.Whois(call.Argument(0).String())
				return otto.TrueValue()
			} else {
				return otto.FalseValue()
			}
		},
		Mode: func(call otto.FunctionCall) otto.Value {
			if len(call.ArgumentList) == 1 && call.ArgumentList[0].IsString() {
				pm.conn.Mode(call.Argument(0).String())
				return otto.TrueValue()
			} else {
				if len(call.ArgumentList) > 1 && call.ArgumentList[0].IsString() {
					var args []string
					for _, arg := range call.ArgumentList[1:] {
						if !arg.IsString() {
							return otto.FalseValue()
						}
						args = append(args, arg.String())
					}
					pm.conn.Mode(call.Argument(0).String(), args...)
				}
				return otto.FalseValue()
			}
		},
	}
	pm.js.Set("IRC", bridge)
}