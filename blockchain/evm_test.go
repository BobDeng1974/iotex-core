// Copyright (c) 2018 IoTeX
// This is an alpha (internal) release and is not suitable for production. This source code is provided 'as is' and no
// warranties are given as to title or non-infringement, merchantability or fitness for purpose and, to the extent
// permitted by law, all liability for your use of the code is disclaimed. This source code is governed by Apache
// License 2.0 that can be found in the LICENSE file.

package blockchain

import (
	"context"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iotexproject/iotex-core/blockchain/action"
	"github.com/iotexproject/iotex-core/config"
	"github.com/iotexproject/iotex-core/pkg/hash"
	"github.com/iotexproject/iotex-core/pkg/util/byteutil"
	ta "github.com/iotexproject/iotex-core/test/testaddress"
	"github.com/iotexproject/iotex-core/testutil"
)

func TestEVM(t *testing.T) {
	fmt.Printf("Test EVM\n")
	require := require.New(t)
	testutil.CleanupPath(t, testTriePath)
	defer testutil.CleanupPath(t, testTriePath)
	testutil.CleanupPath(t, testDBPath)
	defer testutil.CleanupPath(t, testDBPath)

	ctx := context.Background()
	cfg := config.Default
	cfg.Chain.TrieDBPath = testTriePath
	cfg.Chain.ChainDBPath = testDBPath
	bc := NewBlockchain(&cfg, DefaultStateFactoryOption(), BoltDBDaoOption())
	require.NotNil(bc)
	defer func() {
		err := bc.Stop(ctx)
		require.NoError(err)
	}()
	_, err := bc.CreateState(ta.Addrinfo["producer"].RawAddress, Gen.TotalSupply)
	require.NoError(err)
	// data, _ := hex.DecodeString("6080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146078575b600080fd5b348015605957600080fd5b5060766004803603810190808035906020019092919050505060a0565b005b348015608357600080fd5b50608a60aa565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a723058202b8e3ee299d6212c404a3f109eb874d5af929b6d2d701819421e3686c4c82fbd0029")
	data, _ := hex.DecodeString("608060405234801561001057600080fd5b5060df8061001f6000396000f3006080604052600436106049576000357c0100000000000000000000000000000000000000000000000000000000900463ffffffff16806360fe47b114604e5780636d4ce63c146078575b600080fd5b348015605957600080fd5b5060766004803603810190808035906020019092919050505060a0565b005b348015608357600080fd5b50608a60aa565b6040518082815260200191505060405180910390f35b8060008190555050565b600080549050905600a165627a7a7230582002faabbefbbda99b20217cf33cb8ab8100caf1542bf1f48117d72e2c59139aea0029")
	// data, _ := hex.DecodeString("6060604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166341c0e1b581146100585780637bf786f81461006b578063fbf788d61461009c575b005b341561006357600080fd5b6100566100ca565b341561007657600080fd5b61008a600160a060020a03600435166100f1565b60405190815260200160405180910390f35b34156100a757600080fd5b610056600160a060020a036004351660243560ff60443516606435608435610103565b60005433600160a060020a03908116911614156100ef57600054600160a060020a0316ff5b565b60016020526000908152604090205481565b600160a060020a0385166000908152600160205260408120548190861161012957600080fd5b3087876040516c01000000000000000000000000600160a060020a03948516810282529290931690910260148301526028820152604801604051809103902091506001828686866040516000815260200160405260006040516020015260405193845260ff90921660208085019190915260408085019290925260608401929092526080909201915160208103908084039060008661646e5a03f115156101cf57600080fd5b505060206040510351600054600160a060020a039081169116146101f257600080fd5b50600160a060020a03808716600090815260016020526040902054860390301631811161026257600160a060020a0387166000818152600160205260409081902088905582156108fc0290839051600060405180830381858888f19350505050151561025d57600080fd5b6102b7565b6000547f2250e2993c15843b32621c89447cc589ee7a9f049c026986e545d3c2c0c6f97890600160a060020a0316604051600160a060020a03909116815260200160405180910390a186600160a060020a0316ff5b505050505050505600a165627a7a72305820533e856fc37e3d64d1706bcc7dfb6b1d490c8d566ea498d9d01ec08965a896ca0029")
	execution, err := action.NewExecution(ta.Addrinfo["producer"].RawAddress, action.EmptyAddress, 1, big.NewInt(0), uint32(100000), uint32(10), data)
	require.NoError(err)
	execution, err = execution.Sign(ta.Addrinfo["producer"])
	require.NoError(err)
	blk, err := bc.MintNewBlock(nil, nil, []*action.Execution{execution}, ta.Addrinfo["producer"], "")
	require.NoError(err)
	require.Nil(bc.CommitBlock(blk))

	h, _ := hex.DecodeString("8db0d504897721a2cc48659d741f763de31d3567")
	contractAddrHash := byteutil.BytesTo20B(h)
	code, err := bc.GetFactory().GetCode(contractAddrHash)
	require.Nil(err)
	require.Equal(data[31:], code)

	// store to key 0
	contractAddr := "io1qyqsyqcy3kcd2pyfwus69nzgvkwhg8mk8h336dt86pg6cj"
	data, _ = hex.DecodeString("60fe47b1000000000000000000000000000000000000000000000000000000000000000f")
	execution, err = action.NewExecution(ta.Addrinfo["producer"].RawAddress, contractAddr, 2, big.NewInt(0), uint32(120000), uint32(10), data)
	require.NoError(err)
	execution, err = execution.Sign(ta.Addrinfo["producer"])
	fmt.Printf("execution %+v\n", execution)
	require.NoError(err)
	blk, err = bc.MintNewBlock(nil, nil, []*action.Execution{execution}, ta.Addrinfo["producer"], "")
	require.NoError(err)
	require.Nil(bc.CommitBlock(blk))

	v, err := bc.GetFactory().GetContractState(contractAddrHash, hash.ZeroHash32B)
	require.Nil(err)
	require.Equal(byte(15), v[31])

	// read from key 0
	contractAddr = "io1qyqsyqcy3kcd2pyfwus69nzgvkwhg8mk8h336dt86pg6cj"
	data, _ = hex.DecodeString("6d4ce63c")
	execution, err = action.NewExecution(ta.Addrinfo["producer"].RawAddress, contractAddr, 3, big.NewInt(0), uint32(120000), uint32(10), data)
	require.NoError(err)
	execution, err = execution.Sign(ta.Addrinfo["producer"])
	fmt.Printf("execution %+v\n", execution)
	require.NoError(err)
	blk, err = bc.MintNewBlock(nil, nil, []*action.Execution{execution}, ta.Addrinfo["producer"], "")
	require.NoError(err)
	require.Nil(bc.CommitBlock(blk))
}