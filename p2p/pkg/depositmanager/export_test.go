package depositmanager

func (d *DepositManager) GetPendingRefund(commitmentDigest CommitmentDigest) (PendingRefund, bool) {
	return d.pendingRefunds.Get(commitmentDigest)
}
