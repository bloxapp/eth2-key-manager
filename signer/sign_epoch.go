package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/pkg/errors"
)

// SignEpoch signs the given epoch
func (signer *SimpleSigner) SignEpoch(epoch phase0.Epoch, domain phase0.Domain, pubKey []byte) ([]byte, []byte, error) {
	// 1. check we can even sign this
	// TODO - should we?

	// 2. get the account
	if pubKey == nil {
		return nil, nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, nil, err
	}

	root, err := ComputeETHSigningRoot(SSZUint64(epoch), domain)
	if err != nil {
		return nil, nil, err
	}

	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, nil, err
	}

	return sig, root[:], nil
}
