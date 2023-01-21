package signer

import (
	"encoding/hex"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/bloxapp/ssv-spec/types"
	"github.com/pkg/errors"
)

// SignAggregateAndProof signs aggregate and proof
func (signer *SimpleSigner) SignAggregateAndProof(agg *phase0.AggregateAndProof, domain phase0.Domain, pubKey []byte) ([]byte, error) {
	// 1. check we can even sign this
	// TODO - should we?

	// 2. get the account
	if pubKey == nil {
		return nil, errors.New("account was not supplied")
	}

	account, err := signer.wallet.AccountByPublicKey(hex.EncodeToString(pubKey))
	if err != nil {
		return nil, err
	}

	root, err := types.ComputeETHSigningRoot(agg, domain)
	if err != nil {
		return nil, err
	}

	sig, err := account.ValidationKeySign(root[:])
	if err != nil {
		return nil, err
	}

	return sig, nil
}
