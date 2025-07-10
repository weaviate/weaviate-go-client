package alias

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/alias"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/fault"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestAlias_integration(t *testing.T) {
	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create Alias for non-existing class should fail", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)

		alias := &alias.Alias{
			Alias: "Band-Alias",
			Class: "Band",
		}
		err := client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.Error(t, err) // should cause error.
	})

	t.Run("Create Alias for existing class", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:       "Band",
			Description: "Band that plays and produces music",
		}

		alias := &alias.Alias{
			Alias: "Band-Alias",
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)
		defer func() {
			errRm := client.Schema().AllDeleter().Do(context.Background())
			assert.Nil(t, errRm)
		}()

		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias.Alias).Do(ctx))
		}()
	})

	t.Run("Create same Alias for same existing class should fail", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:       "Band",
			Description: "Band that plays and produces music",
		}

		alias := &alias.Alias{
			Alias: "Band-Alias",
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)
		defer func() {
			errRm := client.Schema().AllDeleter().Do(context.Background())
			assert.Nil(t, errRm)
		}()

		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias.Alias).Do(ctx))
		}()

		// do it again for same alias and same class. Should throw error
		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.Error(t, err)
	})

	t.Run("Alias creation should be globally unique", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:       "Band",
			Description: "Band that plays and produces music",
		}

		schemaClass2 := &models.Class{
			Class:       "Band2",
			Description: "Band that plays and produces different music",
		}

		alias2 := &alias.Alias{
			Alias: "Band-Alias",
			Class: schemaClass2.Class,
		}

		alias := &alias.Alias{
			Alias: "Band-Alias",
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassDeleter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassCreator().WithClass(schemaClass2).Do(ctx)
		require.NoError(t, err)
		defer func() {
			errRm := client.Schema().AllDeleter().Do(context.Background())
			assert.Nil(t, errRm)
		}()

		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias.Alias).Do(ctx))
		}()

		// do it again for same alias for different class. Should throw error. Because same alias cannot point to
		// different class (thus globally unique)
		err = client.Alias().AliasCreator().WithAlias(alias2).Do(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "already exists") // we should lock that error is not any other error, it's error saying alias already exists.
	})

	t.Run("Get all alias", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)
		cn1 := "TestGetAllBand"
		cn2 := "TestGetAllBand2"
		an1 := "TestGetAllAlias"
		an2 := "TestGetAllAlias2"

		schemaClass := &models.Class{
			Class:       cn1,
			Description: "Band that plays and produces music",
		}

		schemaClass2 := &models.Class{
			Class:       cn2,
			Description: "Band that plays and produces different music",
		}

		alias2 := &alias.Alias{
			Alias: an2,
			Class: schemaClass2.Class,
		}

		alias := &alias.Alias{
			Alias: an1,
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassDeleter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassCreator().WithClass(schemaClass2).Do(ctx)
		require.NoError(t, err)
		defer func() {
			errRm := client.Schema().AllDeleter().Do(context.Background())
			assert.Nil(t, errRm)
		}()

		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias.Alias).Do(ctx))
		}()

		err = client.Alias().AliasCreator().WithAlias(alias2).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias2.Alias).Do(ctx))
		}()

		// list all alias
		resp, err := client.Alias().Getter().Do(ctx)
		require.NoError(t, err)
		assert.Len(t, resp, 2)
	})
	t.Run("Get alias for specific collection", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)
		cn1 := "TestGetAllBand"
		cn2 := "TestGetAllBand2"
		an1 := "TestGetAllAlias"
		an2 := "TestGetAllAlias2"

		schemaClass := &models.Class{
			Class:       cn1,
			Description: "Band that plays and produces music",
		}

		schemaClass2 := &models.Class{
			Class:       cn2,
			Description: "Band that plays and produces different music",
		}

		alias2 := &alias.Alias{
			Alias: an2,
			Class: schemaClass2.Class,
		}

		alias := &alias.Alias{
			Alias: an1,
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassDeleter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassCreator().WithClass(schemaClass2).Do(ctx)
		require.NoError(t, err)
		defer func() {
			errRm := client.Schema().AllDeleter().Do(context.Background())
			assert.Nil(t, errRm)
		}()

		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias.Alias).Do(ctx))
		}()

		err = client.Alias().AliasCreator().WithAlias(alias2).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias2.Alias).Do(ctx))
		}()

		// list alias for specific class
		resp, err := client.Alias().Getter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, schemaClass.Class, resp[0].Class)

		// list alias for specific class
		resp, err = client.Alias().Getter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, schemaClass2.Class, resp[0].Class)

		// Also verify via /alises/{alias} endpoint.
		respSingle, err := client.Alias().AliasGetter().WithAliasName(alias.Alias).Do(ctx)
		require.NoError(t, err)
		require.NotNil(t, respSingle)
		require.Equal(t, alias, respSingle)

		// list alias for specific class
		respSingle, err = client.Alias().AliasGetter().WithAliasName(alias2.Alias).Do(ctx)
		require.NoError(t, err)
		require.NotNil(t, respSingle)
		require.Equal(t, alias2, respSingle)
	})

	t.Run("Update alias from one collection to another", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)
		cn1 := "TestUpdateBand"
		cn2 := "TestUpdateBand2"
		an1 := "TestUpdateAlias"

		schemaClass := &models.Class{
			Class:       cn1,
			Description: "Band that plays and produces music",
		}

		schemaClass2 := &models.Class{
			Class:       cn2,
			Description: "Band that plays and produces different music",
		}

		alias := &alias.Alias{
			Alias: an1,
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassDeleter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassCreator().WithClass(schemaClass2).Do(ctx)
		require.NoError(t, err)
		defer func() {
			errRm := client.Schema().AllDeleter().Do(context.Background())
			assert.Nil(t, errRm)
		}()

		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias.Alias).Do(ctx))
		}()

		// list alias for specific class
		resp, err := client.Alias().Getter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, schemaClass.Class, resp[0].Class)
		assert.Equal(t, alias.Alias, resp[0].Alias)

		// update
		alias.Class = schemaClass2.Class
		err = client.Alias().AliasUpdater().WithAlias(alias).Do(ctx)
		require.NoError(t, err)

		// list alias for specific class
		resp, err = client.Alias().Getter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, schemaClass2.Class, resp[0].Class)
	})
	t.Run("Update alias from one collection to another which doesn't exist", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)
		cn1 := "TestUpdateBand"
		an1 := "TestUpdateAlias"

		schemaClass := &models.Class{
			Class:       cn1,
			Description: "Band that plays and produces music",
		}

		bandAlias := &alias.Alias{
			Alias: an1,
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)
		defer func() {
			errRm := client.Schema().AllDeleter().Do(context.Background())
			assert.Nil(t, errRm)
		}()

		err = client.Alias().AliasCreator().WithAlias(bandAlias).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(bandAlias.Alias).Do(ctx))
		}()

		// list alias for specific class
		resp, err := client.Alias().Getter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, schemaClass.Class, resp[0].Class)
		assert.Equal(t, bandAlias.Alias, resp[0].Alias)

		// update should fail
		err = client.Alias().AliasUpdater().WithAlias(&alias.Alias{
			Alias: an1,
			Class: "Unknown",
		}).Do(ctx)
		require.Error(t, err)
		require.IsType(t, (*fault.WeaviateClientError)(nil), err)
		require.Equal(t, 422, err.(*fault.WeaviateClientError).StatusCode)
	})
	t.Run("Delete alias that doesn't exist", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)
		cn1 := "TestUpdateBand"
		an1 := "TestUpdateAlias"

		schemaClass := &models.Class{
			Class:       cn1,
			Description: "Band that plays and produces music",
		}

		alias := &alias.Alias{
			Alias: an1,
			Class: schemaClass.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)

		err = client.Schema().ClassCreator().WithClass(schemaClass).Do(ctx)
		require.NoError(t, err)
		defer func() {
			errRm := client.Schema().AllDeleter().Do(context.Background())
			assert.Nil(t, errRm)
		}()

		err = client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().WithAliasName(alias.Alias).Do(ctx))
		}()

		// list alias for specific class. Make sure alias "foo" doesn't exist
		resp, err := client.Alias().AliasGetter().WithAliasName("foo").Do(ctx)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "404")
		assert.Nil(t, resp)

		err = client.Alias().AliasDeleter().WithAliasName("foo").Do(ctx) // that doesn't exist
		require.Error(t, err)
		assert.Contains(t, err.Error(), "404")
	})

	t.Run("tear down weaviate", func(t *testing.T) {
		err := testenv.TearDownLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to tear down weaviate: %s", err)
		}
	})
}
