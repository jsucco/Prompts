package shard_counter

import (
	"fmt"
	"math/rand"
	"golang.org/x/net/context"
	"cloud.google.com/go/datastore"
	"google.golang.org/appengine/memcache"
	"google.golang.org/api/iterator"
)

type counterConfig struct {
	Shards int
}

type shard struct {
	Name  string
	Count int
}

const (
	defaultShards = 20
	configKind    = "GeneralCounterShardConfig"
	shardKind     = "GeneralCounterShard"
	projectID = "project-alpha-170622"
)

func memcacheKey(name string) string {
	return shardKind + ":" + name
}

// Count retrieves the value of the named counter.
func Count(ctx context.Context, name string) (int, error) {
	total := 0
	client, _ := datastore.NewClient(ctx, projectID)
	mkey := memcacheKey(name)
	if _, errm := memcache.JSON.Get(ctx, mkey, &total); errm == nil {
		return total, nil
	}
	q := datastore.NewQuery(shardKind).Filter("Name =", name)
	it := client.Run(ctx, q)
	for  {
		var s shard

		if _, err := it.Next(&s); err == iterator.Done {
			break
		} else if err != nil {
			return -1, err
		}
		total += s.Count
	}
	memcache.JSON.Set(ctx, &memcache.Item{
		Key:        mkey,
		Object:     &total,
		Expiration: 60,
	})
	return total, nil
}

// Increment increments the named counter.
func Increment(ctx context.Context, name string) error {
	client, _ := datastore.NewClient(ctx, projectID)
	var cfg counterConfig
	ckey := datastore.NameKey(configKind, name, nil)
	_,err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		err := client.Get(ctx, ckey, &cfg)
		if err == datastore.ErrNoSuchEntity {
			cfg.Shards = defaultShards
			_, err = client.Put(ctx, ckey, &cfg)
		}
		return err
	})
	if err != nil {
		return err
	}
	var s shard
	_, err = client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		shardName := fmt.Sprintf("%s-shard%d", name, rand.Intn(cfg.Shards))
		key := datastore.NameKey(shardKind, shardName, nil)
		err := client.Get(ctx, key, &s)
		// A missing entity and a present entity will both work.
		if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}
		s.Name = name
		s.Count++
		_, err = client.Put(ctx, key, &s)
		return err
	})
	if err != nil {
		return err
	}
	memcache.IncrementExisting(ctx, memcacheKey(name), 1)
	return nil
}

// IncreaseShards increases the number of shards for the named counter to n.
// It will never decrease the number of shards.
func IncreaseShards(ctx context.Context, name string, n int) error {
	client, _ := datastore.NewClient(ctx, projectID)
	ckey := datastore.NameKey(configKind, name, nil)
	_, err := client.RunInTransaction(ctx, func(tx *datastore.Transaction) error {
		var cfg counterConfig
		mod := false
		err := client.Get(ctx, ckey, &cfg)
		if err == datastore.ErrNoSuchEntity {
			cfg.Shards = defaultShards
			mod = true
		} else if err != nil {
			return err
		}
		if cfg.Shards < n {
			cfg.Shards = n
			mod = true
		}
		if mod {
			_, err = client.Put(ctx, ckey, &cfg)
		}
		return err
	})
	return err
}
