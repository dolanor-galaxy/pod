package blockchain

import (
	"math"

	chaincfg "github.com/p9c/pod/pkg/chain/config"
)

const (
	// // vbLegacyBlockVersion is the highest legacy block version before the version bits scheme became active.
	// vbLegacyBlockVersion = 4
	// vbTopBits defines the bits to set in the version to signal that the version bits scheme is being used.
	vbTopBits = 0x20000000
	// vbTopMask is the bitmask to use to determine whether or not the version bits scheme is in use.
	vbTopMask = 0xe0000000
	// vbNumBits is the total number of bits available for use with the version bits scheme.
	vbNumBits = 29
	// // unknownVerNumToCheck is the number of previous blocks to consider when checking for a threshold of unknown block versions for the purposes of warning the user.
	// unknownVerNumToCheck = 100
	// // unknownVerWarnNum is the threshold of previous blocks that have an unknown version to use for the purposes of warning the user.
	// unknownVerWarnNum = unknownVerNumToCheck / 2
)

type // bitConditionChecker provides a thresholdConditionChecker which can be
// used to test whether or not a specific bit is set when it's not supposed
// to be according to the expected version based on the known deployments and
// the current state of the chain.
// This is useful for detecting and warning about unknown rule activations.
bitConditionChecker struct {
	bit   uint32
	chain *BlockChain
}

// Ensure the bitConditionChecker type implements the
// thresholdConditionChecker interface.
var _ thresholdConditionChecker = bitConditionChecker{}

func // BeginTime returns the unix timestamp for the median block time after
// which voting on a rule change starts (at the next window).
// Since this implementation checks for unknown rules,
// it returns 0 so the rule is always treated as active.
// This is part of the thresholdConditionChecker interface implementation.
(c bitConditionChecker) BeginTime() uint64 {
	return 0
}

func // EndTime returns the unix timestamp for the median block time after which
// an attempted rule change fails if it has not already been locked in or
// activated. Since this implementation checks for unknown rules,
// it returns the maximum possible timestamp so the rule is always treated as
// active. This is part of the thresholdConditionChecker interface
// implementation.
(c bitConditionChecker) EndTime() uint64 {
	return math.MaxUint64
}

func // RuleChangeActivationThreshold is the number of blocks for which the
// condition must be true in order to lock in a rule change.
// This implementation returns the value defined by the chain netparams the
// checker is associated with.
// This is part of the thresholdConditionChecker interface implementation.
(c bitConditionChecker) RuleChangeActivationThreshold() uint32 {
	return c.chain.params.RuleChangeActivationThreshold
}

func // MinerConfirmationWindow is the number of blocks in each threshold state
// retarget window. This implementation returns the value defined by the
// chain netparams the checker is associated with.
// This is part of the thresholdConditionChecker interface implementation.
(c bitConditionChecker) MinerConfirmationWindow() uint32 {
	return c.chain.params.MinerConfirmationWindow
}

func // Condition returns true when the specific bit associated with the checker
// is set and it's not supposed to be according to the expected version based
// on the known deployments and the current state of the chain.
// This function MUST be called with the chain state lock held (for writes).
// This is part of the thresholdConditionChecker interface implementation.
(c bitConditionChecker) Condition(node *BlockNode) (bool, error) {
	conditionMask := uint32(1) << c.bit
	version := uint32(node.version)
	if version&vbTopMask != vbTopBits {
		return false, nil
	}
	if version&conditionMask == 0 {
		return false, nil
	}
	expectedVersion, err := c.chain.calcNextBlockVersion(node.parent)
	if err != nil {
		Error(err)
		return false, err
	}
	return expectedVersion&conditionMask == 0, nil
}

type // deploymentChecker provides a thresholdConditionChecker which can be used
// to test a specific deployment rule.
// This is required for properly detecting and activating consensus rule
// changes.
deploymentChecker struct {
	deployment *chaincfg.ConsensusDeployment
	chain      *BlockChain
}

var // Ensure the deploymentChecker type implements the
// thresholdConditionChecker interface.
_ thresholdConditionChecker = deploymentChecker{}

func // BeginTime returns the unix timestamp for the median block time after
// which voting on a rule change starts (at the next window).
// This implementation returns the value defined by the specific deployment
// the checker is associated with.
// This is part of the thresholdConditionChecker interface implementation.
(c deploymentChecker) BeginTime() uint64 {
	return c.deployment.StartTime
}

func // EndTime returns the unix timestamp for the median block time after which
// an attempted rule change fails if it has not already been locked in or
// activated. This implementation returns the value defined by the specific
// deployment the checker is associated with.
// This is part of the thresholdConditionChecker interface implementation.
(c deploymentChecker) EndTime() uint64 {
	return c.deployment.ExpireTime
}

func // RuleChangeActivationThreshold is the number of blocks for which the
// condition must be true in order to lock in a rule change.
// This implementation returns the value defined by the chain netparams the
// checker is associated with.
// This is part of the thresholdConditionChecker interface implementation.
(c deploymentChecker) RuleChangeActivationThreshold() uint32 {
	return c.chain.params.RuleChangeActivationThreshold
}

func // MinerConfirmationWindow is the number of blocks in each threshold state
// retarget window. This implementation returns the value defined by the
// chain netparams the checker is associated with.
// This is part of the thresholdConditionChecker interface implementation.
(c deploymentChecker) MinerConfirmationWindow() uint32 {
	return c.chain.params.MinerConfirmationWindow
}

func // Condition returns true when the specific bit defined by the deployment
// associated with the checker is set.
// This is part of the thresholdConditionChecker interface implementation.
(c deploymentChecker) Condition(node *BlockNode) (bool, error) {
	conditionMask := uint32(1) << c.deployment.BitNumber
	version := uint32(node.version)
	return (version&vbTopMask == vbTopBits) && (version&conditionMask != 0),
		nil
}

func // calcNextBlockVersion calculates the expected version of the block after
// the passed previous block node based on the state of started and locked in
// rule change deployments.
// This function differs from the exported CalcNextBlockVersion in that the
// exported version uses the current best chain as the previous block node
// while this function accepts any block node.
// This function MUST be called with the chain state lock held (for writes).
(b *BlockChain) calcNextBlockVersion(prevNode *BlockNode) (uint32, error) {
	// Set the appropriate bits for each actively defined rule deployment
	// that is either in the process of being voted on,
	// or locked in for the/ activation at the next threshold window change.
	expectedVersion := uint32(vbTopBits)
	for id := 0; id < len(b.params.Deployments); id++ {
		deployment := &b.params.Deployments[id]
		cache := &b.deploymentCaches[id]
		checker := deploymentChecker{deployment: deployment, chain: b}
		state, err := b.thresholdState(prevNode, checker, cache)
		if err != nil {
			Error(err)
			return 0, err
		}
		if state == ThresholdStarted || state == ThresholdLockedIn {
			expectedVersion |= uint32(1) << deployment.BitNumber
		}
	}
	return expectedVersion, nil
}

func // CalcNextBlockVersion calculates the expected version of the block after
// the end of the current best chain based on the state of started and locked
// in rule change deployments. This function is safe for concurrent access.
(b *BlockChain) CalcNextBlockVersion() (uint32, error) {
	b.chainLock.Lock()
	version, err := b.calcNextBlockVersion(b.BestChain.Tip())
	b.chainLock.Unlock()
	return version, err
}

func // warnUnknownRuleActivations displays a warning when any unknown new rules
// are either about to activate or have been activated.
// This will only happen once when new rules have been activated and every
// block for those about to be activated.
// This function MUST be called with the chain state lock held (for writes)
(b *BlockChain) warnUnknownRuleActivations(node *BlockNode) error {
	// Warn if any unknown new rules are either about to activate or have already been activated.
	for bit := uint32(0); bit < vbNumBits; bit++ {
		checker := bitConditionChecker{bit: bit, chain: b}
		cache := &b.warningCaches[bit]
		state, err := b.thresholdState(node.parent, checker, cache)
		if err != nil {
			Error(err)
			return err
		}
		switch state {
		case ThresholdActive:
			if !b.unknownRulesWarned {
				Warnf("unknown new rules activated (bit %d)", bit)
				b.unknownRulesWarned = true
			}
		case ThresholdLockedIn:
			window := int32(checker.MinerConfirmationWindow())
			activationHeight := window - (node.height % window)
			Warnf("Unknown new rules are about to activate in %d blocks ("+
				"bit %d)", activationHeight, bit)
		}
	}
	return nil
}

// warnUnknownVersions logs a warning if a high enough percentage of the last
// blocks have unexpected versions.
// This function MUST be called with the chain state lock held (for writes)
// func (b *BlockChain) warnUnknownVersions(node *BlockNode) error {
// 	// Nothing to do if already warned.
// 	if b.unknownVersionsWarned {
// 		return nil
// 	}
// 	// Warn if enough previous blocks have unexpected versions.
// 	numUpgraded := uint32(0)
// 	for i := uint32(0); i < unknownVerNumToCheck && node != nil; i++ {
// 		expectedVersion, err := b.calcNextBlockVersion(node.parent)
// 		if err != nil {
// Error(err)
// 			return err
// 		}
// 		if expectedVersion > vbLegacyBlockVersion &&
// 			(node.version & ^expectedVersion) != 0 {
// 			numUpgraded++
// 		}
// 		node = node.parent
// 	}
// 	if numUpgraded > unknownVerWarnNum {
// 		WARN{"Unknown block versions are being mined, so new " +
// 			"rules might be in effect.  Are you running the " +
// 			"latest version of the software?")
// 		b.unknownVersionsWarned = true
// 	}
// 	return nil
// }
