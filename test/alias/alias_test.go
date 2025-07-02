package alias

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/weaviate/weaviate-go-client/v5/test/testsuit"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/testenv"
	"github.com/weaviate/weaviate/entities/models"
)

func TestAlias_integration(t *testing.T) {
	// 1. List of all alias
	// 2. Get specif alias that exists
	// 3. Get specif alias that doesn't exists

	// update
	// 1. Update alias from one collection to other (both exists)
	// 1. Update alias from one collection to other (other doesn't exists)
	// 1. Update alias from one collection to other (one doesn't exists)
	// 1. Update alias from one collection to other (both doesn't exists)
	// 1. Update alias from one collection to other without proper payload (invalid payload)

	// delete
	// 1. Delete alias that exists.
	// 2. Delete alias that doesn't exist

	t.Run("up", func(t *testing.T) {
		err := testenv.SetupLocalWeaviate()
		if err != nil {
			t.Fatalf("failed to setup weaviate: %s", err)
		}
	})

	t.Run("Create Alias for non-existing class should fail", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)

		alias := &models.Alias{
			Alias: "Band-Alias",
			Class: "Band",
		}
		err := client.Alias().AliasCreator().WithAlias(alias).Do(ctx)
		require.Error(t, err) // should cause error.

		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx))
		}()

	})

	t.Run("Create Alias for existing class", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:       "Band",
			Description: "Band that plays and produces music",
		}

		alias := &models.Alias{
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
			require.NoError(t, client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx))
		}()

	})

	t.Run("Create same Alias for same existing class should fail", func(t *testing.T) {
		ctx := context.Background()
		client := testsuit.CreateTestClient(false)

		schemaClass := &models.Class{
			Class:       "Band",
			Description: "Band that plays and produces music",
		}

		alias := &models.Alias{
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
			require.NoError(t, client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx))
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

		alias := &models.Alias{
			Alias: "Band-Alias",
			Class: schemaClass.Class,
		}

		alias2 := &models.Alias{
			Alias: "Band-Alias",
			Class: schemaClass2.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassDeleter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx)
		require.NoError(t, err)
		err = client.Alias().AliasDeleter().AliasName(alias2.Alias).Do(ctx)
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
			require.NoError(t, client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx))
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

		alias := &models.Alias{
			Alias: an1,
			Class: schemaClass.Class,
		}

		alias2 := &models.Alias{
			Alias: an2,
			Class: schemaClass2.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassDeleter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx)
		require.NoError(t, err)
		err = client.Alias().AliasDeleter().AliasName(alias2.Alias).Do(ctx)
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
			require.NoError(t, client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx))
		}()

		err = client.Alias().AliasCreator().WithAlias(alias2).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx))
		}()

		// list all alias
		resp, err := client.Alias().AliasLister().Do(ctx)
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

		alias := &models.Alias{
			Alias: an1,
			Class: schemaClass.Class,
		}

		alias2 := &models.Alias{
			Alias: an2,
			Class: schemaClass2.Class,
		}

		err := client.Schema().ClassDeleter().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Schema().ClassDeleter().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)
		err = client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx)
		require.NoError(t, err)
		err = client.Alias().AliasDeleter().AliasName(alias2.Alias).Do(ctx)
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
			require.NoError(t, client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx))
		}()

		err = client.Alias().AliasCreator().WithAlias(alias2).Do(ctx)
		require.NoError(t, err)
		defer func() {
			require.NoError(t, client.Alias().AliasDeleter().AliasName(alias.Alias).Do(ctx))
		}()

		// list alias for specific class
		resp, err := client.Alias().AliasLister().WithClassName(schemaClass.Class).Do(ctx)
		require.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, schemaClass.Class, resp[0].Class)

		// list alias for specific class
		resp, err = client.Alias().AliasLister().WithClassName(schemaClass2.Class).Do(ctx)
		require.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, schemaClass2.Class, resp[0].Class)
	})
}
