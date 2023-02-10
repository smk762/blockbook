//go:build unittest

package komodo

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"os"
	"reflect"
	"testing"

	"github.com/martinboehm/btcutil/chaincfg"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/btc"
)

var (
	testTx1, testTx2 bchain.Tx

	testTxPacked1 = "0a20d5d22c2f9b2c073ef0ce27feefcf4fc007cfa967f8221285e17e0b910103eac712bf030400008085202f8903f91fdc3290ffc8301c6ee9ea288f761b8afb4ab7092ab667b5e753f13d69e5a2000000004847304402202e7097d9ee757ad63944f3b2ef2104fdc7d2545ef5b13ff0577bc7d914938d7502207f904354670bc67e13357b00ceea9f09ba7cbd8cfefd4846ebe805e9446f7a7701feffffff36fa4a4780feec666e7129eab73b2c7a15518539a3e5e0cb5cc986fc5d234364000000004948304502210092a33956ab51482395c7f8dfedb15a6308bce5a81a86b387ce2a58b895630d9f02202524b73f8e0f22ff577620976cc1e26bcd8c5b55700083981ab2477e55dbeaa401feffffff39510819cccd49d372b50afef9f0e924503667ed665a3a5b9cf0f3df2fcdda2a00000000484730440220364688bb679d94b2f92d0a4f943feaf79f8cba2e743508a5993d332c463630690220588095ca66a6bf80322a70b3322953d50e9ad5259e2e5840c6aca8438e5c555a01feffffff028afc5407000000002321038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abac4014502e000000001976a9145d904d4531f4c74f760a14ef057c866a06705dda88ac16ffe563422e3200000000000000000000000018c380989f062096fe979f0628fbdac801327312206443235dfc86c95ccbe0e5a3398551157a2c3bb7ea29716e66ecfe80474afa36224948304502210092a33956ab51482395c7f8dfedb15a6308bce5a81a86b387ce2a58b895630d9f02202524b73f8e0f22ff577620976cc1e26bcd8c5b55700083981ab2477e55dbeaa40128feffffff0f32731220a2e5693df153e7b567b62a09b74afb8a1b768f28eae96e1c30c8ff9032dc1ff9224948304502210092a33956ab51482395c7f8dfedb15a6308bce5a81a86b387ce2a58b895630d9f02202524b73f8e0f22ff577620976cc1e26bcd8c5b55700083981ab2477e55dbeaa40128feffffff0f327212202adacd2fdff3f09c5b3a5a66ed67365024e9f0f9fe0ab572d349cdcc1908513922484730440220364688bb679d94b2f92d0a4f943feaf79f8cba2e743508a5993d332c463630690220588095ca66a6bf80322a70b3322953d50e9ad5259e2e5840c6aca8438e5c555a0128feffffff0f3a4f0a040754fc8a1a2321038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abac222252447261676f4e4864776f767673444c534c4d6941457a4541724144336b7136464e3a470a042e50144010011a1976a9145d904d4531f4c74f760a14ef057c866a06705dda88ac222252486f756e643850707968564c666935366443374d4b335a76766b416d4233627651"
	testTxPacked2 = "0a207f6a417dd1117aefb96632d8db6b54c2c8768bb7345177f43080d9005efdead112ef020400008085202f89029c6cdb08308bbf6e6466a760483f5550e1e71fb98bada5cfe9dabaa07279268200000000484730440220118fb52a5f5d3f652de24fee228122021d18fb8906f5113b9b817aa59174b6e702203c3693915ad53262cbdf79f76e72f29e526c807cd1599cad9ca9d42bc06a8a7101feffffff6229e5b750ef8c82ca8ff50bdcdb5652d815bf7cf0179e31aae76a8973aa5c82140000006a47304402205131f0b42ae373f5d787df3bd69b5a7e37c2ed0e1ed6c123e0e78f286de7db270220770fd72120f07e541029e8df815e323d41e349f53b91dd72fd217d2c40f023990121038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abfeffffff02381d3000000000002321038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abac007841cb020000001976a91432311a35188a9439c6c866e842564d6fefd3a02888ace98de6639b303200000000000000000000000018d59f9a9f0620e99b9a9f0628d4dfc8013272122082267972a0badae9cfa5ad8bb91fe7e150553f4860a766646ebf8b3008db6c9c22484730440220118fb52a5f5d3f652de24fee228122021d18fb8906f5113b9b817aa59174b6e702203c3693915ad53262cbdf79f76e72f29e526c807cd1599cad9ca9d42bc06a8a710128feffffff0f3296011220825caa73896ae7aa319e17f07cbf15d85256dbdc0bf58fca828cef50b7e529621814226a47304402205131f0b42ae373f5d787df3bd69b5a7e37c2ed0e1ed6c123e0e78f286de7db270220770fd72120f07e541029e8df815e323d41e349f53b91dd72fd217d2c40f023990121038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34ab28feffffff0f3a500a0502cb4178001a2321038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abac222252447261676f4e4864776f767673444c534c4d6941457a4541724144336b7136464e3a460a0398968010011a1976a91432311a35188a9439c6c866e842564d6fefd3a02888ac222252447261676f4e4864776f767673444c534c4d6941457a4541724144336b7136464e"
)

func init() {
	testTx1 = bchain.Tx{
		Hex:       "0400008085202f8903f91fdc3290ffc8301c6ee9ea288f761b8afb4ab7092ab667b5e753f13d69e5a2000000004847304402202e7097d9ee757ad63944f3b2ef2104fdc7d2545ef5b13ff0577bc7d914938d7502207f904354670bc67e13357b00ceea9f09ba7cbd8cfefd4846ebe805e9446f7a7701feffffff36fa4a4780feec666e7129eab73b2c7a15518539a3e5e0cb5cc986fc5d234364000000004948304502210092a33956ab51482395c7f8dfedb15a6308bce5a81a86b387ce2a58b895630d9f02202524b73f8e0f22ff577620976cc1e26bcd8c5b55700083981ab2477e55dbeaa401feffffff39510819cccd49d372b50afef9f0e924503667ed665a3a5b9cf0f3df2fcdda2a00000000484730440220364688bb679d94b2f92d0a4f943feaf79f8cba2e743508a5993d332c463630690220588095ca66a6bf80322a70b3322953d50e9ad5259e2e5840c6aca8438e5c555a01feffffff028afc5407000000002321038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abac4014502e000000001976a9145d904d4531f4c74f760a14ef057c866a06705dda88ac16ffe563422e32000000000000000000000000",
		Blocktime: 1676017731,
		Time:      1676017731,
		Txid:      "d5d22c2f9b2c073ef0ce27feefcf4fc007cfa967f8221285e17e0b910103eac7",
		LockTime:  1676017430,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "48304502210092a33956ab51482395c7f8dfedb15a6308bce5a81a86b387ce2a58b895630d9f02202524b73f8e0f22ff577620976cc1e26bcd8c5b55700083981ab2477e55dbeaa401",
				},
				Txid:     "6443235dfc86c95ccbe0e5a3398551157a2c3bb7ea29716e66ecfe80474afa36",
				Vout:     0,
				Sequence: 4294967294,
			},
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "48304502210092a33956ab51482395c7f8dfedb15a6308bce5a81a86b387ce2a58b895630d9f02202524b73f8e0f22ff577620976cc1e26bcd8c5b55700083981ab2477e55dbeaa401",
				},
				Txid:     "a2e5693df153e7b567b62a09b74afb8a1b768f28eae96e1c30c8ff9032dc1ff9",
				Vout:     0,
				Sequence: 4294967294,
			},
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "4730440220364688bb679d94b2f92d0a4f943feaf79f8cba2e743508a5993d332c463630690220588095ca66a6bf80322a70b3322953d50e9ad5259e2e5840c6aca8438e5c555a01",
				},
				Txid:     "2adacd2fdff3f09c5b3a5a66ed67365024e9f0f9fe0ab572d349cdcc19085139",
				Vout:     0,
				Sequence: 4294967294,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(123010186),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "21038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abac",
					Addresses: []string{
						"RDragoNHdwovvsDLSLMiAEzEArAD3kq6FN",
					},
				},
			},
			{
				ValueSat: *big.NewInt(777000000),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a9145d904d4531f4c74f760a14ef057c866a06705dda88ac",
					Addresses: []string{
						"RHound8PpyhVLfi56dC7MK3ZvvkAmB3bvQ",
					},
				},
			},
		},
	}

	testTx2 = bchain.Tx{
		Hex:       "0400008085202f89029c6cdb08308bbf6e6466a760483f5550e1e71fb98bada5cfe9dabaa07279268200000000484730440220118fb52a5f5d3f652de24fee228122021d18fb8906f5113b9b817aa59174b6e702203c3693915ad53262cbdf79f76e72f29e526c807cd1599cad9ca9d42bc06a8a7101feffffff6229e5b750ef8c82ca8ff50bdcdb5652d815bf7cf0179e31aae76a8973aa5c82140000006a47304402205131f0b42ae373f5d787df3bd69b5a7e37c2ed0e1ed6c123e0e78f286de7db270220770fd72120f07e541029e8df815e323d41e349f53b91dd72fd217d2c40f023990121038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abfeffffff02381d3000000000002321038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abac007841cb020000001976a91432311a35188a9439c6c866e842564d6fefd3a02888ace98de6639b3032000000000000000000000000",
		Blocktime: 1676054485,
		Time:      1676054485,
		Txid:      "7f6a417dd1117aefb96632d8db6b54c2c8768bb7345177f43080d9005efdead1",
		LockTime:  1676053993,
		Vin: []bchain.Vin{
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "4730440220118fb52a5f5d3f652de24fee228122021d18fb8906f5113b9b817aa59174b6e702203c3693915ad53262cbdf79f76e72f29e526c807cd1599cad9ca9d42bc06a8a7101",
				},
				Txid:     "82267972a0badae9cfa5ad8bb91fe7e150553f4860a766646ebf8b3008db6c9c",
				Vout:     0,
				Sequence: 4294967294,
			},
			{
				ScriptSig: bchain.ScriptSig{
					Hex: "47304402205131f0b42ae373f5d787df3bd69b5a7e37c2ed0e1ed6c123e0e78f286de7db270220770fd72120f07e541029e8df815e323d41e349f53b91dd72fd217d2c40f023990121038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34ab",
				},
				Txid:     "825caa73896ae7aa319e17f07cbf15d85256dbdc0bf58fca828cef50b7e52962",
				Vout:     20,
				Sequence: 4294967294,
			},
		},
		Vout: []bchain.Vout{
			{
				ValueSat: *big.NewInt(12000000000),
				N:        0,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "21038e010c33c56b61389409eea5597fe17967398731e23185c84c472a16fc5d34abac",
					Addresses: []string{
						"RDragoNHdwovvsDLSLMiAEzEArAD3kq6FN",
					},
				},
			},
			{
				ValueSat: *big.NewInt(10000000),
				N:        1,
				ScriptPubKey: bchain.ScriptPubKey{
					Hex: "76a91432311a35188a9439c6c866e842564d6fefd3a02888ac",
					Addresses: []string{
						"RDragoNHdwovvsDLSLMiAEzEArAD3kq6FN",
					},
				},
			},
		},
	}
}

func TestMain(m *testing.M) {
	c := m.Run()
	chaincfg.ResetParams()
	os.Exit(c)
}

func TestGetAddrDesc(t *testing.T) {
	type args struct {
		tx     bchain.Tx
		parser *KomodoParser
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "zec-1",
			args: args{
				tx:     testTx1,
				parser: NewKomodoParser(GetChainParams("main"), &btc.Configuration{}),
			},
		},
		{
			name: "zec-2",
			args: args{
				tx:     testTx2,
				parser: NewKomodoParser(GetChainParams("main"), &btc.Configuration{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for n, vout := range tt.args.tx.Vout {
				got1, err := tt.args.parser.GetAddrDescFromVout(&vout)
				if err != nil {
					t.Errorf("getAddrDescFromVout() error = %v, vout = %d", err, n)
					return
				}
				got2, err := tt.args.parser.GetAddrDescFromAddress(vout.ScriptPubKey.Addresses[0])
				if err != nil {
					t.Errorf("getAddrDescFromAddress() error = %v, vout = %d", err, n)
					return
				}
				if !bytes.Equal(got1, got2) {
					t.Errorf("Address descriptors mismatch: got1 = %v, got2 = %v", got1, got2)
				}
			}
		})
	}
}

func TestPackTx(t *testing.T) {
	type args struct {
		tx        bchain.Tx
		height    uint32
		blockTime int64
		parser    *KomodoParser
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "zec-1",
			args: args{
				tx:        testTx1,
				height:    3288443,
				blockTime: 1676017731,
				parser:    NewKomodoParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked1,
			wantErr: false,
		},
		{
			name: "zec-2",
			args: args{
				tx:        testTx2,
				height:    3289044,
				blockTime: 1676054485,
				parser:    NewKomodoParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    testTxPacked2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.args.parser.PackTx(&tt.args.tx, tt.args.height, tt.args.blockTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("packTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			h := hex.EncodeToString(got)
			if !reflect.DeepEqual(h, tt.want) {
				t.Errorf("packTx() = %v, want %v", h, tt.want)
			}
		})
	}
}

func TestUnpackTx(t *testing.T) {
	type args struct {
		packedTx string
		parser   *KomodoParser
	}
	tests := []struct {
		name    string
		args    args
		want    *bchain.Tx
		want1   uint32
		wantErr bool
	}{
		{
			name: "kmd-1",
			args: args{
				packedTx: testTxPacked1,
				parser:   NewKomodoParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx1,
			want1:   3288443,
			wantErr: false,
		},
		{
			name: "kmd-1",
			args: args{
				packedTx: testTxPacked2,
				parser:   NewKomodoParser(GetChainParams("main"), &btc.Configuration{}),
			},
			want:    &testTx2,
			want1:   3289044,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, _ := hex.DecodeString(tt.args.packedTx)
			got, got1, err := tt.args.parser.UnpackTx(b)
			if (err != nil) != tt.wantErr {
				t.Errorf("unpackTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("unpackTx() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("unpackTx() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
