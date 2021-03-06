/*
 * provendb-verify
 * Copyright (C) 2019  Southbank Software Ltd.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 *
 * @Author: guiguan
 * @Date:   2018-08-24T09:56:10+10:00
 * @Last modified by:   guiguan
 * @Last modified time: 2019-04-02T13:25:43+11:00
 */

package anchor

import (
	"context"
	"strings"
	"testing"

	"github.com/SouthbankSoftware/provendb-verify/pkg/proof/testutil"
	log "github.com/sirupsen/logrus"
)

func Test_verifyAnchorURIs(t *testing.T) {
	var canceledCtx context.Context

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type args struct {
		ctx           context.Context
		uris          []interface{}
		expectedValue string
	}
	tests := []struct {
		name         string
		args         args
		wantedErrStr string
	}{
		{
			"Verify Calendar anchor URIs",
			args{
				context.Background(),
				[]interface{}{
					"https://a.chainpoint.org/calendar/985635/hash",
					"https://a.chainpoint.org/calendar/985635/hash",
				},
				"4690932f928fb7f7ce6e6c49ee95851742231709360be28b7ce2af7b92cfa95b",
			},
			"",
		},
		{
			"Verify Bitcoin anchor URIs",
			args{
				context.Background(),
				[]interface{}{
					"https://a.chainpoint.org/calendar/985814/data",
				},
				"c617f5faca34474bea7020d75c39cb8427a32145f9646586ecb9184002131ad9",
			},
			"",
		},
		{
			"Verify unknown URIs",
			args{
				context.Background(),
				[]interface{}{
					"http://skldfjklasdfk.com",
				},
				"",
			},
			"Get http://skldfjklasdfk.com: dial tcp: lookup skldfjklasdfk.com: no such host",
		},
		{
			"Verify with canceled context",
			args{
				canceledCtx,
				[]interface{}{
					"http://skldfjklasdfk.com",
				},
				"",
			},
			"context canceled",
		},
		{
			"Verify with 404 status code",
			args{
				context.Background(),
				[]interface{}{
					"https://a.chainpoint.org/notexists",
				},
				"",
			},
			`got 404 Not Found from https://a.chainpoint.org/notexists`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyAnchorURIs(tt.args.ctx, tt.args.uris, tt.args.expectedValue)

			if err != nil {
				log.Error(err)

				errStr := err.Error()

				if errStr != tt.wantedErrStr {
					const commSuffix = "no such host"

					if strings.HasSuffix(errStr, commSuffix) &&
						strings.HasSuffix(tt.wantedErrStr, commSuffix) {
						return
					}

					t.Errorf("verifyAnchorURIs() errStr = %v, wantErrStr %v", errStr, tt.wantedErrStr)
				}
			} else if "" != tt.wantedErrStr {
				t.Errorf("verifyAnchorURIs() errStr = , wantErrStr %v", tt.wantedErrStr)
			}
		})
	}
}

func Test_verifyBitcoinBlockMerkleRoot(t *testing.T) {
	var canceledCtx context.Context

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	type args struct {
		ctx           context.Context
		blockHeight   string
		expectedValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Verify BTC block merkle root",
			args{
				context.Background(),
				"503275",
				"c617f5faca34474bea7020d75c39cb8427a32145f9646586ecb9184002131ad9",
			},
			false,
		},
		{
			"Verify BTC block merkle root with canceled context",
			args{
				canceledCtx,
				"503275",
				"c617f5faca34474bea7020d75c39cb8427a32145f9646586ecb9184002131ad9",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := verifyBitcoinBlockMerkleRoot(tt.args.ctx, tt.args.blockHeight, tt.args.expectedValue)

			if err != nil {
				log.Error(err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("verifyBitcoinBlockMerkleRoot() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyBitcoinTxOpReturn(t *testing.T) {
	type args struct {
		ctx           context.Context
		txID          string
		expectedValue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Verify BTC transaction OP_RETURN",
			args{
				context.Background(),
				"ba3c8c3e547ed73471c28a69659373f3f0a3b726aab31cdecd14513d9c581f1e",
				"267335262e21e7adb4220068b4b90b7ff066324935d7f61ceab2a64080b06b1b",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyBitcoinTxOpReturn(tt.args.ctx, tt.args.txID, tt.args.expectedValue); (err != nil) != tt.wantErr {
				t.Errorf("verifyBitcoinTxOpReturn() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyCalendarBranch(t *testing.T) {
	type args struct {
		ctx    context.Context
		branch map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Verify Calendar anchor branch",
			args{
				context.Background(),
				map[string]interface{}{
					"label": "cal_anchor_branch",
					"anchors": []interface{}{
						map[string]interface{}{
							"type":      "cal",
							"anchor_id": "985637",
							"uris": []string{
								"https://a.chainpoint.org/calendar/985637/hash",
							},
							"expected_value": "4690932f928fb7f7ce6e6c49ee95851742231709360be28b7ce2af7b92cfa95b",
						},
						map[string]interface{}{
							"type":      "cal",
							"anchor_id": "985635",
							"uris": []string{
								"https://a.chainpoint.org/calendar/985635/hash",
								"https://a.chainpoint.org/calendar/985635/hash",
								"https://a.chainpoint.org/calendar/985635/hash",
							},
							"expected_value": "4690932f928fb7f7ce6e6c49ee95851742231709360be28b7ce2af7b92cfa95b",
						},
						map[string]interface{}{
							"type":      "cal",
							"anchor_id": "985635",
							"uris": []string{
								"https://a.chainpoint.org/calendar/985635/hash",
							},
							"expected_value": "4690932f928fb7f7ce6e6c49ee95851742231709360be28b7ce2af7b92cfa95b",
						},
					},
				},
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyCalendarBranch(tt.args.ctx, tt.args.branch); (err != nil) != tt.wantErr {
				t.Errorf("verifyCalendarBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_verifyBitcoinBranch(t *testing.T) {
	type args struct {
		ctx    context.Context
		branch map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Verify Bitcoin anchor branch",
			args{
				context.Background(),
				map[string]interface{}{
					"label": "btc_anchor_branch",
					"anchors": []interface{}{
						map[string]interface{}{
							"type":      "btc",
							"anchor_id": "503275",
							"uris": []interface{}{
								"https://a.chainpoint.org/calendar/985814/data",
							},
							"expected_value": "c617f5faca34474bea7020d75c39cb8427a32145f9646586ecb9184002131ad9",
						},
					},
					"opReturnValue": "267335262e21e7adb4220068b4b90b7ff066324935d7f61ceab2a64080b06b1b",
					"btcTxId":       "ba3c8c3e547ed73471c28a69659373f3f0a3b726aab31cdecd14513d9c581f1e",
				},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := verifyBitcoinBranch(tt.args.ctx, tt.args.branch); (err != nil) != tt.wantErr {
				t.Errorf("verifyBitcoinBranch() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	type args struct {
		ctx            context.Context
		evaluatedProof interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"Verify evaluated Chainpoint v3 Proof - evaluated_proof1.json",
			args{
				context.Background(),
				testutil.LoadJSON(t, "evaluated_proof1.json"),
			},
			false,
		},
		{
			"Verify evaluated Chainpoint v3 Proof - evaluated_proof2.json",
			args{
				context.Background(),
				testutil.LoadJSON(t, "evaluated_proof2.json"),
			},
			false,
		},
		{
			"Verify evaluated Chainpoint v3 Proof - evaluated_proof3.json",
			args{
				context.Background(),
				testutil.LoadJSON(t, "evaluated_proof3.json"),
			},
			false,
		},
	}
	for _, tt := range tests {
		tt := tt // capture range variable tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // run subtest in parallel

			if err := Verify(tt.args.ctx, tt.args.evaluatedProof); (err != nil) != tt.wantErr {
				t.Errorf("Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
