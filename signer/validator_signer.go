package signer

import (
	"sync"

	"github.com/attestantio/go-eth2-client/api"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/capella"
	"github.com/attestantio/go-eth2-client/spec/phase0"
	ssz "github.com/ferranbt/fastssz"
	"github.com/google/uuid"

	"github.com/bloxapp/eth2-key-manager/core"
)

// ValidatorSigner represents the behavior of the validator signer
type ValidatorSigner interface {
	SignBeaconBlock(block *api.VersionedProposal, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignBlindedBeaconBlock(block *api.VersionedBlindedProposal, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignBeaconAttestation(attestation *phase0.AttestationData, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignAggregateAndProof(agg *phase0.AggregateAndProof, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignSlot(slot phase0.Slot, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignEpoch(epoch phase0.Epoch, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignSyncCommittee(msgBlockRoot []byte, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignSyncCommitteeSelectionData(data *altair.SyncAggregatorSelectionData, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignSyncCommitteeContributionAndProof(contribAndProof *altair.ContributionAndProof, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignRegistration(registration *api.VersionedValidatorRegistration, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignVoluntaryExit(voluntaryExit *phase0.VoluntaryExit, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
	SignBLSToExecutionChange(blsToExecutionChange *capella.BLSToExecutionChange, domain phase0.Domain, pubKey []byte) (sig []byte, root []byte, err error)
}

// SimpleSigner implements ValidatorSigner interface
type SimpleSigner struct {
	wallet            core.Wallet
	slashingProtector core.SlashingProtector
	network           core.Network
	signLocks         map[string]*sync.RWMutex
	mapLock           *sync.RWMutex
}

// NewSimpleSigner is the constructor of SimpleSigner
func NewSimpleSigner(wallet core.Wallet, slashingProtector core.SlashingProtector, network core.Network) *SimpleSigner {
	return &SimpleSigner{
		wallet:            wallet,
		slashingProtector: slashingProtector,
		network:           network,
		signLocks:         map[string]*sync.RWMutex{},
		mapLock:           &sync.RWMutex{},
	}
}

// lock locks signer
func (signer *SimpleSigner) lock(accountID uuid.UUID, operation string) {
	signer.mapLock.Lock()
	defer signer.mapLock.Unlock()

	k := accountID.String() + "_" + operation
	if val, ok := signer.signLocks[k]; ok {
		val.Lock()
	} else {
		signer.signLocks[k] = &sync.RWMutex{}
		signer.signLocks[k].Lock()
	}
}

func (signer *SimpleSigner) unlock(accountID uuid.UUID, operation string) {
	signer.mapLock.RLock()
	defer signer.mapLock.RUnlock()

	k := accountID.String() + "_" + operation
	if val, ok := signer.signLocks[k]; ok {
		val.Unlock()
	}
}

// ComputeETHSigningRoot returns computed root for eth signing
func ComputeETHSigningRoot(obj ssz.HashRoot, domain phase0.Domain) (phase0.Root, error) {
	root, err := obj.HashTreeRoot()
	if err != nil {
		return phase0.Root{}, err
	}
	signingContainer := phase0.SigningData{
		ObjectRoot: root,
		Domain:     domain,
	}
	return signingContainer.HashTreeRoot()
}
