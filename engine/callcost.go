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

package engine

import (
	"errors"
	"time"

	"github.com/cgrates/cgrates/utils"
)

// The output structure that will be returned with the call cost information.
type CallCost struct {
	Category           string
	Tenant             string
	Subject            string
	Account            string
	Destination        string
	ToR                string
	Cost               float64
	Timespans          TimeSpans
	RatedUsage         float64
	AccountSummary     *AccountSummary
	deductConnectFee   bool
	negativeConnectFee bool // the connect fee went negative on default balance
	maxCostDisconect   bool
}

// Merges the received timespan if they are similar (same activation period, same interval, same minute info.
func (cc *CallCost) Merge(other *CallCost) {
	cc.Timespans = append(cc.Timespans, other.Timespans...)
	cc.Cost += other.Cost
}

func (cc *CallCost) GetStartTime() time.Time {
	if len(cc.Timespans) == 0 {
		return time.Now()
	}
	return cc.Timespans[0].TimeStart
}

func (cc *CallCost) GetEndTime() time.Time {
	if len(cc.Timespans) == 0 {
		return time.Now()
	}
	return cc.Timespans[len(cc.Timespans)-1].TimeEnd
}

func (cc *CallCost) GetDuration() (td time.Duration) {
	for _, ts := range cc.Timespans {
		td += ts.GetDuration()
	}
	return
}

func (cc *CallCost) UpdateRatedUsage() time.Duration {
	if cc == nil {
		return 0
	}
	totalDuration := cc.GetDuration()
	cc.RatedUsage = float64(totalDuration.Nanoseconds())
	return totalDuration
}

func (cc *CallCost) GetConnectFee() float64 {
	if len(cc.Timespans) == 0 ||
		cc.Timespans[0].RateInterval == nil ||
		cc.Timespans[0].RateInterval.Rating == nil {
		return 0
	}
	return cc.Timespans[0].RateInterval.Rating.ConnectFee
}

// Creates a CallDescriptor structure copying related data from CallCost
func (cc *CallCost) CreateCallDescriptor() *CallDescriptor {
	return &CallDescriptor{
		Category:    cc.Category,
		Tenant:      cc.Tenant,
		Subject:     cc.Subject,
		Account:     cc.Account,
		Destination: cc.Destination,
		ToR:         cc.ToR,
	}
}

func (cc *CallCost) ToDataCost() (*DataCost, error) {
	if cc.ToR == utils.MetaVoice {
		return nil, errors.New("Not a data call!")
	}
	dc := &DataCost{
		Category:         cc.Category,
		Tenant:           cc.Tenant,
		Subject:          cc.Subject,
		Account:          cc.Account,
		Destination:      cc.Destination,
		ToR:              cc.ToR,
		Cost:             cc.Cost,
		deductConnectFee: cc.deductConnectFee,
	}
	dc.DataSpans = make([]*DataSpan, len(cc.Timespans))
	for i, ts := range cc.Timespans {
		length := ts.TimeEnd.Sub(ts.TimeStart).Nanoseconds()
		callDuration := ts.DurationIndex.Nanoseconds()
		dc.DataSpans[i] = &DataSpan{
			DataStart:      float64(callDuration - length),
			DataEnd:        float64(callDuration),
			Cost:           ts.Cost,
			ratingInfo:     ts.ratingInfo,
			RateInterval:   ts.RateInterval,
			DataIndex:      float64(callDuration),
			MatchedSubject: ts.MatchedSubject,
			MatchedPrefix:  ts.MatchedPrefix,
			MatchedDestId:  ts.MatchedDestId,
			RatingPlanId:   ts.RatingPlanId,
		}
		dc.DataSpans[i].Increments = make([]*DataIncrement, len(ts.Increments))
		for j, incr := range ts.Increments {
			dc.DataSpans[i].Increments[j] = &DataIncrement{
				Amount:         float64(incr.Duration.Nanoseconds()),
				Cost:           incr.Cost,
				BalanceInfo:    incr.BalanceInfo,
				CompressFactor: incr.CompressFactor,
			}
		}
	}
	return dc, nil
}

func (cc *CallCost) GetLongestRounding() (roundingDecimals int, roundingMethod string) {
	for _, ts := range cc.Timespans {
		if ts.RateInterval != nil && ts.RateInterval.Rating.RoundingDecimals > roundingDecimals { //ToDo: When will ts.RateInterval be empty?
			roundingDecimals = ts.RateInterval.Rating.RoundingDecimals
			roundingMethod = ts.RateInterval.Rating.RoundingMethod
		}
	}
	return
}

// public function to update final (merged) callcost
func (cc *CallCost) UpdateCost() {
	cc.deductConnectFee = true
	cc.updateCost()
}

func (cc *CallCost) updateCost() {
	cost := 0.0
	//if cc.deductConnectFee { // add back the connectFee
	//	cost += cc.GetConnectFee()
	//}
	for _, ts := range cc.Timespans {
		ts.Cost = ts.CalculateCost()
		cost += ts.Cost
		cost = utils.Round(cost, globalRoundingDecimals, utils.MetaRoundingMiddle) // just get rid of the extra decimals
	}
	cc.Cost = cost
}

// Round creates the RoundIncrements in timespans
func (cc *CallCost) Round() {
	if len(cc.Timespans) == 0 || cc.Timespans[0] == nil {
		return
	}
	var totalCorrectionCost float64
	for _, ts := range cc.Timespans {
		if len(ts.Increments) == 0 {
			continue // safe check
		}
		inc := ts.Increments[0]
		if inc.BalanceInfo == nil || inc.BalanceInfo.Monetary == nil || inc.Cost == 0 {
			// this is a unit paid timespan, nothing to round
			continue
		}
		cost := ts.CalculateCost()
		roundedCost := utils.Round(cost,
			ts.RateInterval.Rating.RoundingDecimals,
			ts.RateInterval.Rating.RoundingMethod)
		correctionCost := roundedCost - cost
		//log.Print(cost, roundedCost, correctionCost)
		if correctionCost != 0 {
			ts.RoundIncrement = &Increment{
				Cost:           correctionCost,
				BalanceInfo:    inc.BalanceInfo,
				CompressFactor: 1,
			}
			totalCorrectionCost += correctionCost
			ts.Cost += correctionCost
		}
	}
	cc.Cost += totalCorrectionCost
}

func (cc *CallCost) GetRoundIncrements() (roundIncrements Increments) {
	for _, ts := range cc.Timespans {
		if ts.RoundIncrement != nil && ts.RoundIncrement.Cost != 0 {
			roundIncrements = append(roundIncrements, ts.RoundIncrement)
		}
	}
	return
}

func (cc *CallCost) MatchCCFilter(bf *BalanceFilter) bool {
	if bf == nil {
		return true
	}
	if bf.Categories != nil && cc.Category != "" && !(*bf.Categories)[cc.Category] {
		return false
	}
	// match destination ids
	foundMatchingDestID := false
	if bf.DestinationIDs != nil && cc.Destination != "" {
		for _, p := range utils.SplitPrefix(cc.Destination, MIN_PREFIX_MATCH) {
			if destIDs, err := dm.GetReverseDestination(p,
				true, true, utils.NonTransactional); err == nil {
				for _, dID := range destIDs {
					if _, ok := (*bf.DestinationIDs)[dID]; ok {
						foundMatchingDestID = true
						break // only one found?
					}
				}
			}
			if foundMatchingDestID {
				break
			}
		}
	} else {
		foundMatchingDestID = true
	}
	if !foundMatchingDestID {
		return false
	}
	return true
}
