/*
Real-time Online/Offline Charging System (OCS) for Telecom & ISP environments
Copyright (C) ITsysCOM GmbH

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>
*/

package console

import (
	v1 "github.com/cgrates/cgrates/apier/v1"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func init() {
	c := &CmdGetTriggers{
		name:      "triggers",
		rpcMethod: utils.APIerSv1GetActionTriggers,
		rpcParams: &v1.AttrGetActionTriggers{},
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

// Commander implementation
type CmdGetTriggers struct {
	name      string
	rpcMethod string
	rpcParams *v1.AttrGetActionTriggers
	*CommandExecuter
}

func (self *CmdGetTriggers) Name() string {
	return self.name
}

func (self *CmdGetTriggers) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdGetTriggers) RpcParams(reset bool) any {
	if reset || self.rpcParams == nil {
		self.rpcParams = &v1.AttrGetActionTriggers{}
	}
	return self.rpcParams
}

func (self *CmdGetTriggers) PostprocessRpcParams() error {
	return nil
}

func (self *CmdGetTriggers) RpcResult() any {
	var atr engine.ActionTriggers
	return &atr
}

func (self *CmdGetTriggers) GetFormatedResult(result any) string {
	return GetFormatedSliceResult(result, utils.StringSet{
		utils.MinSleep: {},
	})
}
