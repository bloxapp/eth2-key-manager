package core

import (
	"github.com/attestantio/go-eth2-client/spec/bellatrix"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/google/uuid"
)

// ValidatorAccount holds the information and actions needed by validator account keys.
// It holds 2 keys, a validation and a withdrawal key.
// As a minimum, the ValidatorAccount should have at least the validation key.
// Withdrawal key is not mandatory to be present.
type ValidatorAccount interface {
	// ID provides the ID for the account.
	ID() uuid.UUID

	// Name provides the name for the account.
	Name() string

	// BasePath provides the basePath of the account.
	BasePath() string

	// ValidatorPublicKey provides the public key for the validation key.
	ValidatorPublicKey() []byte

	// WithdrawalPublicKey provides the public key for the withdrawal key.
	WithdrawalPublicKey() []byte

	// ValidationKeySign signs data with the validation key.
	ValidationKeySign(data []byte) ([]byte, error)

	// GetDepositData returns deposit data
	GetDepositData() (map[string]interface{}, error)

	// SetContext sets the given context
	SetContext(ctx *WalletContext)
}

// ValidatorInfo represents the information of a validator
type ValidatorInfo struct {
	Index                 phase0.ValidatorIndex
	Pubkey                phase0.BLSPubKey
	WithdrawalCredentials []byte
	ToExecutionAddress    bellatrix.ExecutionAddress
}
