package mnemonic_test

import (
	"context"
	"encoding/hex"
	"testing"

	"github.com/prysmaticlabs/prysm/shared/testutil/assert"
	"github.com/prysmaticlabs/remote-signer/keyvault/mnemonic"
)

func TestMnemonicGeneration(t *testing.T) {
	ctx := context.Background()

	mnemonicPhrase := "voice gospel easy verb front diesel sense worth sword equip giggle jeans shoe defy kid degree van frost like blush chef silk spoil obtain"
	mnemonicPassword := ""

	const numOfTests = 3
	startIndexForMnemonic := [numOfTests]int{0, 1, 2}
	numMnemonicKeys := [numOfTests]int{5, 2, 1}

	pubKeysTestVector := [5]string{
		/* Index 0*/ "9731de7d206fcd68bb4fb34c515192adeb63448de22d8d84bd2faad9d1450a6869c46c5ce8a65b4243ad51cff120b9ae",
		/* Index 1*/ "98dbc04dbec1261cc26aebc684c7606288fcb890236b0f92a0436911b09ccb5c11b90867d2b94b1f5d67eb92cb8375b2",
		/* Index 2*/ "a587e0690f2ca201054208c9d2f74286b564977ff6dcdad81cf6f6f604a511a5d8c2df2668d62caa482387e1fb807593",
		/* Index 3*/ "84c545d1a5ae820b39d874117c8d2d8524cfcfd48f2790e9df62b305a30096170126f88f1a12dd7cb1c55c3efa09ef35",
		/* Index 4*/ "814c18e38283dd68021789cd523a8f276230671c3a0a960ba8e9d9a66131da4a091871361506844ba899e8c27d156042"}

	for tc := 0; tc < numOfTests; tc++ {
		vault, _ := mnemonic.NewStore(
			mnemonicPhrase,
			mnemonicPassword,
			startIndexForMnemonic[tc],
			numMnemonicKeys[tc])

		pubKeys, _ := vault.GetPublicKeys(ctx)

		for i := startIndexForMnemonic[tc]; i < (numMnemonicKeys[tc] + startIndexForMnemonic[tc]); i++ {

			expected, _ := hex.DecodeString(pubKeysTestVector[i])

			assert.DeepEqual(t, expected, pubKeys[i-startIndexForMnemonic[tc]].Marshal())
		}
	}
}
