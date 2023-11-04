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
	c := &CmdSetDispatcherProfile{
		name:      "dispatchers_profile_set",
		rpcMethod: utils.APIerSv1SetDispatcherProfile,
	}
	commands[c.Name()] = c
	c.CommandExecuter = &CommandExecuter{c}
}

// Commander implementation
type CmdSetDispatcherProfile struct {
	name      string
	rpcMethod string
	rpcParams *v1.DispatcherWithAPIOpts
	*CommandExecuter
}

func (self *CmdSetDispatcherProfile) Name() string {
	return self.name
}

func (self *CmdSetDispatcherProfile) RpcMethod() string {
	return self.rpcMethod
}

func (self *CmdSetDispatcherProfile) RpcParams(reset bool) any {
	if reset || self.rpcParams == nil {
		self.rpcParams = &v1.DispatcherWithAPIOpts{
			DispatcherProfile: new(engine.DispatcherProfile),
			APIOpts:           make(map[string]any),
		}
	}
	return self.rpcParams
}

func (self *CmdSetDispatcherProfile) PostprocessRpcParams() error {
	return nil
}

func (self *CmdSetDispatcherProfile) RpcResult() any {
	var s string
	return &s
}
