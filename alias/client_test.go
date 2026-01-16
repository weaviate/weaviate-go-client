package alias_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v6/alias"
	"github.com/weaviate/weaviate-go-client/v6/internal/api"
	"github.com/weaviate/weaviate-go-client/v6/internal/testkit"
)

func TestAlias_Create(t *testing.T) {
	for _, tt := range []struct {
		name              string
		alias, collection string
		stubs             []testkit.Stub[api.CreateAliasRequest, api.Alias]
		expectErr         bool // Using require.NoError if nil.
	}{
		{
			name:       "create alias",
			alias:      "Alias",
			collection: "Collection",
			stubs: []testkit.Stub[api.CreateAliasRequest, api.Alias]{
				{
					Request:  &api.CreateAliasRequest{Alias: "Alias", Collection: "Collection"},
					Response: api.Alias{Alias: "Alias", Collection: "Collection"},
				},
			},
		},
		{
			name:       "create alias with non-existent collection",
			alias:      "Alias",
			collection: "DoesntExist",
			stubs: []testkit.Stub[api.CreateAliasRequest, api.Alias]{
				{Err: testkit.ErrWhaam},
			},
			expectErr: true,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			// The first response is always consumed by the test itself to GetCreateStatus.
			transport := testkit.NewTransport(t, tt.stubs)
			c := alias.NewClient(transport)
			require.NotNil(t, c, "nil client")

			// GetCreateStatus is part of test setup, always called with t.Context()
			alias, err := c.Create(t.Context(), tt.alias, tt.collection)
			if tt.expectErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.NotNil(t, alias, "nil backup from get-status")

				require.Equal(t, tt.alias, alias.Alias)
				require.Equal(t, tt.collection, alias.Collection)
			}
		})
	}
}
