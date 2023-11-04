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

package migrator

import (
	"github.com/cgrates/cgrates/engine"
	"github.com/cgrates/cgrates/utils"
)

func (m *Migrator) migrateCurrentTPTiming() (err error) {
	tpids, err := m.storDBIn.StorDB().GetTpIds(utils.TBLTPTimings)
	if err != nil {
		return err
	}

	for _, tpid := range tpids {
		ids, err := m.storDBIn.StorDB().GetTpTableIds(tpid, utils.TBLTPTimings,
			utils.TPDistinctIds{"tag"}, map[string]string{}, nil)
		if err != nil {
			return err
		}
		for _, id := range ids {
			tm, err := m.storDBIn.StorDB().GetTPTimings(tpid, id)
			if err != nil {
				return err
			}
			if tm != nil {
				if !m.dryRun {
					if err := m.storDBOut.StorDB().SetTPTimings(tm); err != nil {
						return err
					}
					for _, timing := range tm {
						if err := m.storDBIn.StorDB().RemTpData(utils.TBLTPTimings,
							timing.TPid, map[string]string{"tag": timing.ID}); err != nil {
							return err
						}
					}
					m.stats[utils.TpTiming] += 1
				}
			}
		}
	}
	return
}

func (m *Migrator) migrateTpTimings() (err error) {
	var vrs engine.Versions
	current := engine.CurrentStorDBVersions()
	if vrs, err = m.getVersions(utils.TpTiming); err != nil {
		return
	}
	switch vrs[utils.TpTiming] {
	case current[utils.TpTiming]:
		if m.sameStorDB {
			break
		}
		if err := m.migrateCurrentTPTiming(); err != nil {
			return err
		}
	}
	return m.ensureIndexesStorDB(utils.TBLTPTimings)
}
