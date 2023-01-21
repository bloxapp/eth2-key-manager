package slashingprotection

import (
	"fmt"
	"testing"

	"github.com/attestantio/go-eth2-client/spec/phase0"
	"github.com/stretchr/testify/require"

	"github.com/bloxapp/eth2-key-manager/core"
)

func setupAttestation(t *testing.T, withAttestationData bool) (core.SlashingProtector, []core.ValidatorAccount) {
	err := core.InitBLS()
	require.NoError(t, err)

	// seed
	seed := _byteArray("0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1fff")
	// create an account to use
	vault, err := vault()
	require.NoError(t, err)

	w, err := vault.Wallet()
	require.NoError(t, err)

	account1, err := w.CreateValidatorAccount(seed, nil)
	require.NoError(t, err)

	account2, err := w.CreateValidatorAccount(seed, nil)
	require.NoError(t, err)

	protector := NewNormalProtection(vault.Context.Storage.(core.SlashingStore))
	if !withAttestationData {
		return protector, []core.ValidatorAccount{account1, account2}
	}

	err = protector.UpdateHighestAttestation(account1.ValidatorPublicKey(), &phase0.AttestationData{
		Slot:            30,
		Index:           5,
		BeaconBlockRoot: _byteArray32("A"),
		Source: &phase0.Checkpoint{
			Epoch: 1,
			Root:  _byteArray32("B"),
		},
		Target: &phase0.Checkpoint{
			Epoch: 2,
			Root:  _byteArray32("C"),
		},
	})
	require.NoError(t, err)

	err = protector.UpdateHighestAttestation(account1.ValidatorPublicKey(), &phase0.AttestationData{
		Slot:            30,
		Index:           5,
		BeaconBlockRoot: _byteArray32("A"),
		Source: &phase0.Checkpoint{
			Epoch: 2,
			Root:  _byteArray32("B"),
		},
		Target: &phase0.Checkpoint{
			Epoch: 3,
			Root:  _byteArray32("C"),
		},
	})
	require.NoError(t, err)

	err = protector.UpdateHighestAttestation(account1.ValidatorPublicKey(), &phase0.AttestationData{
		Slot:            30,
		Index:           5,
		BeaconBlockRoot: _byteArray32("B"),
		Source: &phase0.Checkpoint{
			Epoch: 3,
			Root:  _byteArray32("C"),
		},
		Target: &phase0.Checkpoint{
			Epoch: 4,
			Root:  _byteArray32("D"),
		},
	})
	require.NoError(t, err)

	err = protector.UpdateHighestAttestation(account1.ValidatorPublicKey(), &phase0.AttestationData{
		Slot:            30,
		Index:           5,
		BeaconBlockRoot: _byteArray32("B"),
		Source: &phase0.Checkpoint{
			Epoch: 4,
			Root:  _byteArray32("C"),
		},
		Target: &phase0.Checkpoint{
			Epoch: 10,
			Root:  _byteArray32("D"),
		},
	})
	require.NoError(t, err)

	err = protector.UpdateHighestAttestation(account1.ValidatorPublicKey(), &phase0.AttestationData{
		Slot:            30,
		Index:           5,
		BeaconBlockRoot: _byteArray32("B"),
		Source: &phase0.Checkpoint{
			Epoch: 5,
			Root:  _byteArray32("C"),
		},
		Target: &phase0.Checkpoint{
			Epoch: 9,
			Root:  _byteArray32("D"),
		},
	})
	require.NoError(t, err)

	return protector, []core.ValidatorAccount{account1, account2}
}

func TestSurroundingVote(t *testing.T) {
	protector, accounts := setupAttestation(t, true)

	t.Run("1 Surrounded vote", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 2,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 5,
				Root:  _byteArray32("C"),
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})

	t.Run("2 Surrounded votes", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 1,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 7,
				Root:  _byteArray32("C"),
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})

	t.Run("1 Surrounding vote", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 5,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 7,
				Root:  _byteArray32("C"),
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})

	t.Run("2 Surrounding vote", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 6,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 7,
				Root:  _byteArray32("C"),
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})
}

func TestDoubleAttestationVote(t *testing.T) {
	protector, accounts := setupAttestation(t, true)

	t.Run("Different committee index, should slash", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 2,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 3,
				Root:  _byteArray32("C"),
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})

	t.Run("Different block root, should slash", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           5,
			BeaconBlockRoot: _byteArray32("AA"),
			Source: &phase0.Checkpoint{
				Epoch: 2,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 3,
				Root:  _byteArray32("C"),
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})

	t.Run("Same attestation, should be slashable (we can't be sure it's not slashable when using highest att.)", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           5,
			BeaconBlockRoot: _byteArray32("B"),
			Source: &phase0.Checkpoint{
				Epoch: 3,
				Root:  _byteArray32("C"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 4,
				Root:  _byteArray32("D"),
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})

	t.Run("new attestation, should not error", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           5,
			BeaconBlockRoot: _byteArray32("E"),
			Source: &phase0.Checkpoint{
				Epoch: 10,
				Root:  _byteArray32("I"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 11,
				Root:  _byteArray32("H"),
			},
		})
		require.False(t, err != nil || res != nil)
	})
}

func TestMinimalSlashingProtection(t *testing.T) {
	protector, accounts := setupAttestation(t, true)
	at, err := protector.RetrieveHighestAttestation(accounts[0].ValidatorPublicKey())
	require.NoError(t, err)
	fmt.Printf("%d", at.Target.Epoch) // 5,10

	t.Run("source lower than highest source", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 4,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 11,
				Root:  _byteArray32("C"),
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})
	t.Run("source equal to highest source, target equal to highest target", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 5,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 10,
				Root:  _byteArray32("C"),
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})
	t.Run("source higher than highest source, target equal to highest target", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 6,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 10,
				Root:  _byteArray32("C"),
			},
		})

		require.NoError(t, err)
		require.NotNil(t, res)
		require.Equal(t, core.HighestAttestationVote, res.Status)
	})
	t.Run("source equal to highest source, target higher than highest target", func(t *testing.T) {
		res, err := protector.IsSlashableAttestation(accounts[0].ValidatorPublicKey(), &phase0.AttestationData{
			Slot:            30,
			Index:           4,
			BeaconBlockRoot: _byteArray32("A"),
			Source: &phase0.Checkpoint{
				Epoch: 6,
				Root:  _byteArray32("B"),
			},
			Target: &phase0.Checkpoint{
				Epoch: 11,
				Root:  _byteArray32("C"),
			},
		})

		require.NoError(t, err)
		require.Nil(t, res)
	})
}

func TestUpdateLatestAttestation(t *testing.T) {
	protector, accounts := setupAttestation(t, false)
	tests := []struct {
		name                  string
		sourceEpoch           phase0.Epoch
		targetEpoch           phase0.Epoch
		expectedHighestSource uint64
		expectedHighestTarget uint64
	}{
		{
			name:                  "source and epoch zero",
			sourceEpoch:           0,
			targetEpoch:           0,
			expectedHighestSource: 0,
			expectedHighestTarget: 0,
		},
		{
			name:                  "source 0 target 1",
			sourceEpoch:           0,
			targetEpoch:           1,
			expectedHighestSource: 0,
			expectedHighestTarget: 1,
		},
		{
			name:                  "source 10 target 11",
			sourceEpoch:           10,
			targetEpoch:           11,
			expectedHighestSource: 10,
			expectedHighestTarget: 11,
		},
		{
			name:                  "source 11 target 9, can't happen in real life",
			sourceEpoch:           11,
			targetEpoch:           9,
			expectedHighestSource: 11,
			expectedHighestTarget: 11,
		},
		{
			name:                  "source 2 target 9",
			sourceEpoch:           2,
			targetEpoch:           9,
			expectedHighestSource: 11,
			expectedHighestTarget: 11,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			k := accounts[0].ValidatorPublicKey()
			err := protector.UpdateHighestAttestation(k, &phase0.AttestationData{
				Source: &phase0.Checkpoint{
					Epoch: test.sourceEpoch,
				},
				Target: &phase0.Checkpoint{
					Epoch: test.targetEpoch,
				},
			})
			require.NoError(tt, err)

			// Validate highest.
			highest, err := protector.RetrieveHighestAttestation(k)
			require.NoError(tt, err)
			require.EqualValues(tt, highest.Source.Epoch, test.expectedHighestSource)
			require.EqualValues(tt, highest.Target.Epoch, test.expectedHighestTarget)
		})
	}
}
