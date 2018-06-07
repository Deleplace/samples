package condors

import (
	"context"

	xnetcontext "golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

func getAndIncrementHitCount(c context.Context, path string) (int, error) {
	key := keyFor(c, path)
	var hits int

	// See https://godoc.org/google.golang.org/appengine/datastore#hdr-Transactions
	err := datastore.RunInTransaction(c, func(tc xnetcontext.Context) error {
		var x HitCount
		if err := datastore.Get(tc, key, &x); err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		x.Count++
		if _, err := datastore.Put(tc, key, &x); err != nil {
			return err
		}
		hits = x.Count
		return nil
	}, nil)
	return hits, err
}

func keyFor(c context.Context, path string) *datastore.Key {
	return datastore.NewKey(c, "HitCount", path, 0, nil)
}

type HitCount struct {
	Count int
}
