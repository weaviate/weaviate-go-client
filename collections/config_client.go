package collections

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v6/cluster"
	"github.com/weaviate/weaviate-go-client/v6/internal"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/dev"
)

func NewConfigClient(t internal.Transport, rd api.RequestDefaults) *ConfigClient {
	dev.AssertNotNil(t, "transport")
	return &ConfigClient{
		transport: t,
		defaults:  rd,
	}
}

type ConfigClient struct {
	transport internal.Transport
	defaults  api.RequestDefaults
}

// Get returns configuration for the collection.
// Returns nil with nil error if collections does not exist.
func (c *ConfigClient) Get(ctx context.Context) (*Collection, error) {
	var resp api.Collection
	if err := c.transport.Do(ctx, api.GetCollectionRequest(c.defaults.CollectionName), &resp); err != nil {
		return nil, fmt.Errorf("get collection config: %w", err)
	}
	collection, err := collectionFromAPI(&resp)
	if err != nil {
		return nil, err
	}
	return &collection, nil
}

func (c *ConfigClient) AddProperty(ctx context.Context, p Property) error {
	req := &api.AddPropertyRequest{
		RequestDefaults: c.defaults,
		Property: &api.Property{
			Name:             p.Name,
			Description:      p.Description,
			DataType:         api.DataType(p.DataType),
			NestedProperties: nestedPropertiesToAPI(p.NestedProperties),
			Tokenization:     api.Tokenization(p.Tokenization),
			IndexFilterable:  p.IndexFilterable,
			IndexRangeable:   p.IndexRangeable,
			IndexSearchable:  p.IndexSearchable,
		},
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("add property: %w", err)
	}
	return nil
}

func (c *ConfigClient) AddReference(ctx context.Context, ref Reference) error {
	req := &api.AddPropertyRequest{
		RequestDefaults: c.defaults,
		Reference: &api.ReferenceProperty{
			Name:        ref.Name,
			Collections: ref.Collections,
		},
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("add reference property: %w", err)
	}
	return nil
}

type DropPropertyIndexOptions struct {
	PropertyName string
	IndexType    PropertyIndexType
}

type PropertyIndexType api.PropertyIndexType

const (
	IndexFilterable = PropertyIndexType(api.PropertyIndexFilterable)
	IndexSearchable = PropertyIndexType(api.PropertyIndexSearchable)
	IndexRangeable  = PropertyIndexType(api.PropertyIndexRangeable)
)

func (c *ConfigClient) DropPropertyIndex(ctx context.Context, options DropPropertyIndexOptions) error {
	req := &api.DropPropertyIndexRequest{
		RequestDefaults: c.defaults,
		PropertyName:    options.PropertyName,
		IndexType:       api.PropertyIndexType(options.IndexType),
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("drop property index: %w", err)
	}
	return nil
}

type DropVectorIndexOptions struct {
	VectorName string
}

func (c *ConfigClient) DropVectorIndex(ctx context.Context, options DropVectorIndexOptions) error {
	req := &api.DropVectorIndexRequest{
		RequestDefaults: c.defaults,
		VectorName:      options.VectorName,
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("drop vector index: %w", err)
	}
	return nil
}

type Shard struct {
	cluster.Shard
	PerNodeStatus map[string]string
}

func (c *ConfigClient) ListShards(ctx context.Context) ([]Shard, error) {
	req := &api.ListCollectionShardsRequest{
		RequestDefaults: c.defaults,
	}
	var resp api.ListCollectionShardsResponse
	if err := c.transport.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("list collection shards: %w", err)
	}

	shards := make([]Shard, len(resp))
	for i, shard := range resp {
		shards[i] = Shard{
			Shard: cluster.Shard{
				Name: shard.Name,
			},
			PerNodeStatus: shard.PerNodeStatus,
		}
	}

	return shards, nil
}

type ShardStatusOptions struct {
	ShardName string
	Ready     bool
	ReadOnly  bool
}

func (c *ConfigClient) SetShardStatus(ctx context.Context, options ShardStatusOptions) error {
	var status string
	switch {
	case options.Ready:
		status = api.ShardStatusReady
	case options.ReadOnly:
		status = api.ShardStatusReadOnly
	}

	req := &api.UpdateShardStatusRequest{
		RequestDefaults: c.defaults,
		ShardName:       options.ShardName,
		ShardStatus:     status,
	}
	if err := c.transport.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("set status %s for shard %q: %w", status, options.ShardName, err)
	}
	return nil
}

func (c *ConfigClient) SetPropertyDescription(ctx context.Context, propertyName, description string) error {
	if err := c.updateCollection(ctx, func(collection *Collection) (err error) {
		var ok bool
		for i := range collection.Properties {
			p := &collection.Properties[i]
			if p.Name == propertyName {
				p.Description = description
				ok = true
			}
		}
		if !ok {
			err = fmt.Errorf("set property description: no such property %q", propertyName)
		}
		return
	}); err != nil {
		return fmt.Errorf("set property description: %w", err)
	}
	return nil
}

func (c *ConfigClient) UpdateVectorConfig(ctx context.Context, vectorName string, f func(vc *VectorConfig)) error {
	if err := c.updateCollection(ctx, func(collection *Collection) (err error) {
		vc, ok := collection.Vectors[vectorName]
		if !ok {
			err = fmt.Errorf("update vector config: no such vector %q", vectorName)
		}
		f(&vc)
		collection.Vectors[vectorName] = vc
		return
	}); err != nil {
		return fmt.Errorf("set property description: %w", err)
	}
	return nil
}

func (c *ConfigClient) UpdateReplicationConfig(ctx context.Context, f func(rc *ReplicationConfig)) error {
	if err := c.updateCollection(ctx, func(collection *Collection) (err error) {
		rc := collection.Replication
		if rc == nil {
			rc = new(ReplicationConfig)
		}
		f(rc)
		return
	}); err != nil {
		return fmt.Errorf("update replication config: %w", err)
	}
	return nil
}

func (c *ConfigClient) UpdateInvertedIndexConfig(ctx context.Context, f func(iic *InvertedIndexConfig)) error {
	if err := c.updateCollection(ctx, func(collection *Collection) (err error) {
		rc := collection.InvertedIndex
		if rc == nil {
			rc = new(InvertedIndexConfig)
		}
		f(rc)
		return
	}); err != nil {
		return fmt.Errorf("update inverted index config: %w", err)
	}
	return nil
}

func (c *ConfigClient) UpdateObjectTTLConfig(ctx context.Context, f func(ttl *ObjectTTLConfig)) error {
	if err := c.updateCollection(ctx, func(collection *Collection) (err error) {
		ttl := collection.ObjectTTL
		if ttl == nil {
			ttl = new(ObjectTTLConfig)
		}
		f(ttl)
		return
	}); err != nil {
		return fmt.Errorf("update object TTL config: %w", err)
	}
	return nil
}

func (c *ConfigClient) UpdateMultiTenancyConfig(ctx context.Context, f func(mt *MultiTenancyConfig)) error {
	if err := c.updateCollection(ctx, func(collection *Collection) (err error) {
		mt := collection.MultiTenancy
		if mt == nil {
			mt = new(MultiTenancyConfig)
		}
		f(mt)
		return
	}); err != nil {
		return fmt.Errorf("update multi-tenancy config: %w", err)
	}
	return nil
}

func (c *ConfigClient) updateCollection(ctx context.Context, f func(*Collection) error) error {
	collection, err := c.Get(ctx)
	if err != nil {
		return err
	}

	if err := f(collection); err != nil {
		return err
	}

	x, err := collectionToAPI(collection)
	if err != nil {
		return err
	}

	req := &api.UpdateCollectionConfigRequest{Collection: x}

	if err := c.transport.Do(ctx, req, nil); err != nil {
		return err
	}
	return nil
}
