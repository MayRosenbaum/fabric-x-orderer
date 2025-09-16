/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package requestfilter_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/hyperledger/fabric-protos-go-apiv2/common"
	"github.com/hyperledger/fabric-x-common/common/policies"
	policyMock "github.com/hyperledger/fabric-x-orderer/common/policy/mocks"
	"github.com/hyperledger/fabric-x-orderer/common/requestfilter"
	"github.com/hyperledger/fabric-x-orderer/common/requestfilter/mocks"
	"github.com/hyperledger/fabric-x-orderer/node/protos/comm"
	"github.com/hyperledger/fabric-x-orderer/testutil/tx"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

func TestSigVerifyFilter(t *testing.T) {
	var v requestfilter.RulesVerifier
	fc := &mocks.FakeFilterConfig{}

	v.AddRule(requestfilter.NewSigFilter(fc, policies.ChannelWriters))
	err := v.Verify(nil)
	require.EqualError(t, err, "failed to convert request to signedData : nil request")

	req := &comm.Request{}
	err = v.Verify(req)
	require.EqualError(t, err, "failed to convert request to signedData : missing header in request's payload")

	payload := &common.Payload{Header: &common.Header{ChannelHeader: make([]byte, 10), SignatureHeader: nil}}
	p, err := proto.Marshal(payload)
	require.NoError(t, err)
	req.Payload = p
	err = v.Verify(req)
	require.EqualError(t, err, "failed to convert request to signedData : missing signature header in payload's header")

	payload = &common.Payload{Header: &common.Header{ChannelHeader: make([]byte, 10), SignatureHeader: make([]byte, 10)}}
	p, err = proto.Marshal(payload)
	require.NoError(t, err)
	req.Payload = p
	err = v.Verify(req)
	require.ErrorContains(t, err, "failed unmarshalling signature header")

	sigheader, err := proto.Marshal(&common.SignatureHeader{
		Creator: []byte("user"),
		Nonce:   []byte("nonce"),
	})
	require.NoError(t, err)

	payload = &common.Payload{Header: &common.Header{ChannelHeader: make([]byte, 10), SignatureHeader: sigheader}}
	p, err = proto.Marshal(payload)
	require.NoError(t, err)
	req.Payload = p
	err = v.Verify(req)
	require.ErrorContains(t, err, "failed unmarshalling channel header")

	chdr := &common.ChannelHeader{ChannelId: "ChannelId"}
	chdrBytes, err := proto.Marshal(chdr)
	require.NoError(t, err)
	payload = &common.Payload{Header: &common.Header{ChannelHeader: chdrBytes, SignatureHeader: sigheader}}
	p, err = proto.Marshal(payload)
	require.NoError(t, err)
	req.Payload = p
	err = v.Verify(req)
	require.NoError(t, err)
}

func TestSigValidationFlag(t *testing.T) {
	var v requestfilter.RulesVerifier
	req := tx.CreateStructuredRequest([]byte("data"))
	fc := &mocks.FakeFilterConfig{}
	pm := &policyMock.FakePolicyManager{}
	p := &policyMock.FakePolicyEvaluator{}

	pm.GetPolicyReturns(p, false)
	fc.GetPolicyManagerReturns(pm)
	fc.GetClientSignatureVerificationRequiredReturns(true)

	v.AddRule(requestfilter.NewSigFilter(fc, policies.ChannelWriters))

	err := v.Verify(req)
	require.ErrorContains(t, err, "no policies in config block")

	pm.GetPolicyReturns(p, true)
	err = v.Verify(req)
	require.NoError(t, err)

	p.EvaluateSignedDataReturns(errors.New("error"))
	err = v.Verify(req)
	require.ErrorContains(t, err, "signature did not satisfy policy")

	p.EvaluateSignedDataReturns(nil)
	err = v.Verify(req)
	require.NoError(t, err)

	fc.GetClientSignatureVerificationRequiredReturns(false)
	err = v.Update(fc)
	require.NoError(t, err)
	err = v.Verify(req)
	require.NoError(t, err)
}
