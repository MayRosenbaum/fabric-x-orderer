/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package requestfilter

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSignedTransactionService(t *testing.T) {
	numOfTxs := 100
	txSize := 300
	service, err := NewSignedTransactionService(numOfTxs, txSize)
	require.NoError(t, err)

	isValid := service.VerifyTransaction()
	require.True(t, isValid)

	// verify all txs
	start := time.Now()
	for i := 0; i < numOfTxs; i++ {
		isValid = service.VerifyTransaction()
		require.True(t, isValid)
	}
	stop := time.Since(start)
	fmt.Printf("verifying %d txs took: %s\n", numOfTxs, stop)
}
