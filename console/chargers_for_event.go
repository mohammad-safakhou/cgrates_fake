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
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func init() {
	c := &CmdGetChargersForEvent{
		name:      "chargers_for_event",
		rpcMethod: utils.ChargerSv1GetChargersForEvent,
		rpcParams: &utils.CGREvent{},
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

type CmdGetChargersForEvent struct {
	name      string
	rpcMethod string
	rpcParams *utils.CGREvent
	*CommandExecuter
}

func (self *CmdGetChargersForEvent) Name() string {
	return self.name
}

func (self *CmdGetChargersForEvent) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdGetChargersForEvent) RpcParams(reset bool) any {
	if reset || self.rpcParams == nil {
		self.rpcParams = new(utils.CGREvent)
	}
	return self.rpcParams
}

func (self *CmdGetChargersForEvent) PostprocessRpcParams() error {
	return nil
}

func (self *CmdGetChargersForEvent) RpcResult() any {
	var atr engine.ChargerProfiles
	return &atr
}
