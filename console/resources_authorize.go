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
	"time"

	"github.com/cgrates/cgrates/utils"
)

func init() {
	c := &CmdResourceAuthorize{
		name:      "resources_authorize",
		rpcMethod: utils.ResourceSv1AuthorizeResources,
		rpcParams: &utils.CGREvent{},
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

// Commander implementation
type CmdResourceAuthorize struct {
	name      string
	rpcMethod string
	rpcParams *utils.CGREvent
	*CommandExecuter
}

func (self *CmdResourceAuthorize) Name() string {
	return self.name
}

func (self *CmdResourceAuthorize) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdResourceAuthorize) RpcParams(reset bool) any {
	if reset || self.rpcParams == nil {
		self.rpcParams = new(utils.CGREvent)
	}
	return self.rpcParams
}

func (self *CmdResourceAuthorize) PostprocessRpcParams() error {
	if self.rpcParams != nil && self.rpcParams.Time == nil {
		self.rpcParams.Time = utils.TimePointer(time.Now())
	}
	return nil
}

func (self *CmdResourceAuthorize) RpcResult() any {
	var atr string
	return &atr
}
