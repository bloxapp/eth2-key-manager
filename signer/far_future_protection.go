package signer

import (
	"time"

	"github.com/attestantio/go-eth2-client/spec/phase0"
)

// FarFutureMaxValidEpoch is the max epoch of far future signing
var FarFutureMaxValidEpoch = int64(time.Minute.Seconds() * 20)

// IsValidFarFutureEpoch prevents far into the future signing request, verify a slot is within the current epoch
// https://github.com/ethereum/eth2.0-specs/blob/dev/specs/phase0/validator.md#protection-best-practices
func IsValidFarFutureEpoch(network network, epoch phase0.Epoch) bool {
	maxValidEpoch := network.EstimatedEpochAtSlot(network.EstimatedSlotAtTime(time.Now().Unix() + FarFutureMaxValidEpoch))
	return epoch <= maxValidEpoch
}

// IsValidFarFutureSlot returns true if the given slot is valid
func IsValidFarFutureSlot(network network, slot phase0.Slot) bool {
	maxValidSlot := network.EstimatedSlotAtTime(time.Now().Unix() + FarFutureMaxValidEpoch)
	return slot <= maxValidSlot
}
