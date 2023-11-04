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

package services

import (
	"sync"

	"github.com/cgrates/birpc"
	v2 "github.com/cgrates/cgrates/apier/v2"
	"github.com/cgrates/cgrates/config"
	"github.com/cgrates/cgrates/cores"
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

// NewAPIerSv2Service returns the APIerSv2 Service
func NewAPIerSv2Service(apiv1 *APIerSv1Service, cfg *config.CGRConfig,
	server *cores.Server, internalAPIerSv2Chan chan birpc.ClientConnector,
	anz *AnalyzerService, srvDep map[string]*sync.WaitGroup) *APIerSv2Service {
	return &APIerSv2Service{
		apiv1:    apiv1,
		connChan: internalAPIerSv2Chan,
		cfg:      cfg,
		server:   server,
		anz:      anz,
		srvDep:   srvDep,
	}
}

// APIerSv2Service implements Service interface
type APIerSv2Service struct {
	sync.RWMutex
	cfg    *config.CGRConfig
	server *cores.Server

	apiv1    *APIerSv1Service
	api      *v2.APIerSv2
	connChan chan birpc.ClientConnector
	anz      *AnalyzerService
	srvDep   map[string]*sync.WaitGroup
}

// Start should handle the sercive start
// For this service the start should be called from RAL Service
func (api *APIerSv2Service) Start() error {
	if api.IsRunning() {
		return utils.ErrServiceAlreadyRunning
	}

	apiV1Chan := api.apiv1.GetAPIerSv1Chan()
	apiV1 := <-apiV1Chan
	apiV1Chan <- apiV1

	api.Lock()
	defer api.Unlock()

	api.api = &v2.APIerSv2{
		APIerSv1: *apiV1,
	}
	srv, err := engine.NewService(api.api)
	if err != nil {
		return err
	}
	var apierV1Srv *birpc.Service
	apierV1Srv, err = birpc.NewService(apiV1, "", false)
	if err != nil {
		return err
	}
	engine.RegisterPingMethod(apierV1Srv.Methods)
	srv[utils.APIerSv1] = apierV1Srv
	if !api.cfg.DispatcherSCfg().Enabled {
		for _, s := range srv {
			api.server.RpcRegister(s)
		}
		var legacySrv engine.IntService
		legacySrv, err = engine.NewServiceWithName(api.api, utils.ApierV2, true)
		if err != nil {
			return err
		}
		//backwards compatible
		for _, s := range legacySrv {
			api.server.RpcRegister(s)
		}
	}

	api.connChan <- api.anz.GetInternalCodec(srv, utils.APIerSv2)
	return nil
}

// Reload handles the change of config
func (api *APIerSv2Service) Reload() (err error) {
	return
}

// Shutdown stops the service
func (api *APIerSv2Service) Shutdown() (err error) {
	api.Lock()
	api.api = nil
	<-api.connChan
	api.Unlock()
	return
}

// IsRunning returns if the service is running
func (api *APIerSv2Service) IsRunning() bool {
	api.RLock()
	defer api.RUnlock()
	return api != nil && api.api != nil
}

// ServiceName returns the service name
func (api *APIerSv2Service) ServiceName() string {
	return utils.APIerSv2
}

// ShouldRun returns if the service should be running
func (api *APIerSv2Service) ShouldRun() bool {
	return api.cfg.ApierCfg().Enabled
}
