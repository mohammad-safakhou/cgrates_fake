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
	"github.com/cgrates/cgrates/scheduler"
	"github.com/cgrates/cgrates/utils"
)

func init() {
	c := &CmdGetScheduledActions{
		name:      "scheduler_queue",
		rpcMethod: utils.APIerSv1GetScheduledActions,
		rpcParams: &scheduler.ArgsGetScheduledActions{},
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

// Commander implementation
type CmdGetScheduledActions struct {
	name      string
	rpcMethod string
	rpcParams *scheduler.ArgsGetScheduledActions
	*CommandExecuter
}

func (self *CmdGetScheduledActions) Name() string {
	return self.name
}

func (self *CmdGetScheduledActions) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdGetScheduledActions) RpcParams(reset bool) any {
	if reset || self.rpcParams == nil {
		self.rpcParams = &scheduler.ArgsGetScheduledActions{}
	}
	return self.rpcParams
}

func (self *CmdGetScheduledActions) PostprocessRpcParams() error {
	return nil
}

func (self *CmdGetScheduledActions) RpcResult() any {
	s := make([]*scheduler.ScheduledAction, 0)
	return &s
}
