package reprovide_test

import (
	"context"
	"testing"

	testutil "gx/ipfs/QmVvkK7s5imCiq3JVbL3pGfnhcCnf3LrFJPF4GE2sAoGZf/go-testutil"
	pstore "gx/ipfs/QmXauCuJzmzapetmC6W4TuDJLL1yFFrVzSHoWv8YdbmnxH/go-libp2p-peerstore"
	blockstore "gx/ipfs/QmayRSLCiM2gWR7Kay8vqu3Yy5mf7yPqocF9ZRgDUPYMcc/go-ipfs-blockstore"
	mock "gx/ipfs/QmbhrE4WDBUYTujDaKMw6siq2V2dG3XHPizhzVdkf8agPK/go-ipfs-routing/mock"
	ds "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore"
	dssync "gx/ipfs/QmeiCcJfDW1GJnWUArudsv5rQsihpi4oyddPhdqo3CfX6i/go-datastore/sync"
	blocks "gx/ipfs/Qmej7nf81hi2x2tvjRBF3mcp74sQyuDH4VMYDGd1YtXjb2/go-block-format"

	. "github.com/ipfs/go-ipfs/exchange/reprovide"
)

func TestReprovide(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mrserv := mock.NewServer()

	idA := testutil.RandIdentityOrFatal(t)
	idB := testutil.RandIdentityOrFatal(t)

	clA := mrserv.Client(idA)
	clB := mrserv.Client(idB)

	bstore := blockstore.NewBlockstore(dssync.MutexWrap(ds.NewMapDatastore()))

	blk := blocks.NewBlock([]byte("this is a test"))
	bstore.Put(blk)

	keyProvider := NewBlockstoreProvider(bstore)
	reprov := NewReprovider(ctx, clA, keyProvider)
	err := reprov.Reprovide()
	if err != nil {
		t.Fatal(err)
	}

	var providers []pstore.PeerInfo
	maxProvs := 100

	provChan := clB.FindProvidersAsync(ctx, blk.Cid(), maxProvs)
	for p := range provChan {
		providers = append(providers, p)
	}

	if len(providers) == 0 {
		t.Fatal("Should have gotten a provider")
	}

	if providers[0].ID != idA.ID() {
		t.Fatal("Somehow got the wrong peer back as a provider.")
	}
}
