package updater

import "math/big"

func (u *Updater) ComputeResidualAfterDecay(startTimestamp, endTimestamp, commitTimestamp uint64) *big.Int {
	return u.computeResidualAfterDecay(startTimestamp, endTimestamp, commitTimestamp)
}
