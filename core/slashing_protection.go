package core

import (
	pb "github.com/wealdtech/eth2-signer-api/pb/v1"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

// SlashingProtector represents the behavior of the slashing protector
type SlashingProtector interface {
	IsSlashableAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) (*AttestationSlashStatus, error)
	IsSlashableProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) *ProposalSlashStatus
	// Will potentially update the highest attestation given this latest attestation.
	UpdateLatestAttestation(key e2types.PublicKey, req *pb.SignBeaconAttestationRequest) error
	SaveProposal(key e2types.PublicKey, req *pb.SignBeaconProposalRequest) error
	RetrieveHighestAttestation(key e2types.PublicKey) (*BeaconAttestation, error)
}

// SlashingStore represents the behavior of the slashing store
type SlashingStore interface {
	SaveHighestAttestation(key e2types.PublicKey, req *BeaconAttestation) error
	RetrieveHighestAttestation(key e2types.PublicKey) *BeaconAttestation
	SaveProposal(key e2types.PublicKey, req *BeaconBlockHeader) error
	RetrieveProposal(key e2types.PublicKey, slot uint64) (*BeaconBlockHeader, error)
}
