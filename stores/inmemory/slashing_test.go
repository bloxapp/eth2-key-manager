package inmemory

import (
	"math/big"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/google/uuid"
	"github.com/herumi/bls-eth-go-binary/bls"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/core"
)

func _bigInt(input string) *big.Int {
	res, _ := new(big.Int).SetString(input, 10)
	return res
}

type mockAccount struct {
	id            uuid.UUID
	validationKey *big.Int
}

func (a *mockAccount) ID() uuid.UUID    { return a.id }
func (a *mockAccount) Name() string     { return "" }
func (a *mockAccount) BasePath() string { return "" }
func (a *mockAccount) ValidatorPublicKey() []byte {
	sk := &bls.SecretKey{}
	if err := sk.Deserialize(a.validationKey.Bytes()); err != nil {
		return nil
	}
	return sk.GetPublicKey().Serialize()
}
func (a *mockAccount) WithdrawalPublicKey() []byte                     { return nil }
func (a *mockAccount) ValidationKeySign(data []byte) ([]byte, error)   { return nil, nil }
func (a *mockAccount) GetDepositData() (map[string]interface{}, error) { return nil, nil }
func (a *mockAccount) SetContext(ctx *core.WalletContext)              {}

//func testBlock(t *testing.T) *spec.VersionedBeaconBlock {
//	blockByts := "7b22736c6f74223a312c2270726f706f7365725f696e646578223a38352c22706172656e745f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c2273746174655f726f6f74223a227264584c666d704c2f396a4f662b6c7065753152466d4747486a4571315562633955674257576d505236553d222c22626f6479223a7b2272616e64616f5f72657665616c223a226f734657704c79554f664859583549764b727170626f4d5048464a684153456232333057394b32556b4b41774c38577473496e41573138572f555a5a597652384250777267616c4e45316f48775745397468555277584b4574522b767135684f56744e424868626b5831426f3855625a51532b5230787177386a667177396446222c22657468315f64617461223a7b226465706f7369745f726f6f74223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d222c226465706f7369745f636f756e74223a3132382c22626c6f636b5f68617368223a22704f564553434e6d764a31546876484e444576344e7a4a324257494c39417856464e55642f4b3352536b6f3d227d2c226772616666697469223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d222c2270726f706f7365725f736c617368696e6773223a6e756c6c2c2261747465737465725f736c617368696e6773223a6e756c6c2c226174746573746174696f6e73223a5b7b226167677265676174696f6e5f62697473223a2248773d3d222c2264617461223a7b22736c6f74223a302c22636f6d6d69747465655f696e646578223a302c22626561636f6e5f626c6f636b5f726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d222c22736f75726365223a7b2265706f6368223a302c22726f6f74223a22414141414141414141414141414141414141414141414141414141414141414141414141414141414141413d227d2c22746172676574223a7b2265706f6368223a302c22726f6f74223a224f6b4f6b767962375755666f43634879543333476858794b7741666c4e64534f626b374b49534c396432733d227d7d2c227369676e6174757265223a226c37627963617732537751633147587a4c36662f6f5a39616752386562685278503550675a546676676e30344b367879384a6b4c68506738326276674269675641674347767965357a7446797a4772646936555a655a4850593030595a6d3964513939764352674d34676f31666b3046736e684543654d68522f45454b59626a227d5d2c226465706f73697473223a6e756c6c2c22766f6c756e746172795f6578697473223a6e756c6c7d7d"
//	blk := &spec.VersionedBeaconBlock{}
//	require.NoError(t, json.Unmarshal(_byteArray(blockByts), blk))
//	return blk
//}

func getSlashingStorage() core.SlashingStore {
	return NewInMemStore(core.MainNetwork)
}

func TestSavingProposal(t *testing.T) {
	storage := getSlashingStorage()
	tests := []struct {
		name     string
		proposal phase0.Slot
		account  core.ValidatorAccount
	}{
		{
			name:     "simple save",
			proposal: 100,
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveHighestProposal(test.account.ValidatorPublicKey(), test.proposal)
			require.NoError(t, err)

			// fetch
			proposal, err := storage.RetrieveHighestProposal(test.account.ValidatorPublicKey())
			require.NoError(t, err)
			require.NotNil(t, proposal)
			require.EqualValues(t, test.proposal, proposal)
		})
	}
}

func TestSavingAttestation(t *testing.T) {
	storage := getSlashingStorage()
	tests := []struct {
		name    string
		att     *phase0.AttestationData
		account core.ValidatorAccount
	}{
		{
			name: "simple save",
			att: &phase0.AttestationData{
				Slot:            30,
				Index:           1,
				BeaconBlockRoot: [32]byte{},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{},
				},
				Target: &phase0.Checkpoint{
					Epoch: 4,
					Root:  [32]byte{},
				},
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
		{
			name: "simple save with no change to latest attestation target",
			att: &phase0.AttestationData{
				Slot:            30,
				Index:           1,
				BeaconBlockRoot: [32]byte{},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{},
				},
				Target: &phase0.Checkpoint{
					Epoch: 3,
					Root:  [32]byte{},
				},
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveHighestAttestation(test.account.ValidatorPublicKey(), test.att)
			require.NoError(t, err)

			// fetch
			att, err := storage.RetrieveHighestAttestation(test.account.ValidatorPublicKey())
			require.NoError(t, err)
			require.NotNil(t, att)

			// test equal
			aRoot, err := att.HashTreeRoot()
			require.NoError(t, err)
			bRoot, err := test.att.HashTreeRoot()
			require.NoError(t, err)
			require.EqualValues(t, aRoot, bRoot)
		})
	}
}

func TestSavingHighestAttestation(t *testing.T) {
	storage := getSlashingStorage()
	tests := []struct {
		name    string
		att     *phase0.AttestationData
		account core.ValidatorAccount
	}{
		{
			name: "simple save",
			att: &phase0.AttestationData{
				Slot:            30,
				Index:           1,
				BeaconBlockRoot: [32]byte{},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{},
				},
				Target: &phase0.Checkpoint{
					Epoch: 4,
					Root:  [32]byte{},
				},
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
		{
			name: "simple save with no change to latest attestation target",
			att: &phase0.AttestationData{
				Slot:            30,
				Index:           1,
				BeaconBlockRoot: [32]byte{},
				Source: &phase0.Checkpoint{
					Epoch: 1,
					Root:  [32]byte{},
				},
				Target: &phase0.Checkpoint{
					Epoch: 3,
					Root:  [32]byte{},
				},
			},
			account: &mockAccount{
				id:            uuid.New(),
				validationKey: _bigInt("5467048590701165350380985526996487573957450279098876378395441669247373404218"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// save
			err := storage.SaveHighestAttestation(test.account.ValidatorPublicKey(), test.att)
			require.NoError(t, err)

			// fetch
			att, err := storage.RetrieveHighestAttestation(test.account.ValidatorPublicKey())
			require.NoError(t, err)
			require.NotNil(t, att)

			// test equal
			aRoot, err := att.HashTreeRoot()
			require.NoError(t, err)
			bRoot, err := test.att.HashTreeRoot()
			require.NoError(t, err)
			require.EqualValues(t, aRoot, bRoot)
		})
	}
}
