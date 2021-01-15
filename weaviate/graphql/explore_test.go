package graphql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExploreBuilder(t *testing.T) {

	t.Run("Simple Explore", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := Explore{
			connection: conMock,
		}
		nearTextBuilder := NearTextArgumentBuilder{}

		withNearText := nearTextBuilder.WithConcepts([]string{"Cheese", "pineapple"})
		query := builder.WithFields([]ExploreFields{Certainty, Beacon}).
			WithNearText(withNearText).
			build()

		expected := `{Explore(nearText:{concepts: ["Cheese","pineapple"] } ){certainty beacon }}`
		assert.Equal(t, expected, query)
	})

	t.Run("Explore limit and certainty", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := Explore{
			connection: conMock,
		}
		nearTextBuilder := NearTextArgumentBuilder{}

		withNearText := nearTextBuilder.
			WithConcepts([]string{"Cheese"}).
			WithLimit(5).
			WithCertainty(0.71)
		query := builder.WithFields([]ExploreFields{Beacon}).
			WithNearText(withNearText).
			build()

		expected := `{Explore(nearText:{concepts: ["Cheese"] limit: 5 certainty: 0.71 } ){beacon }}`
		assert.Equal(t, expected, query)
	})

	t.Run("Explore with move", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := Explore{
			connection: conMock,
		}
		nearTextBuilder := NearTextArgumentBuilder{}

		fields := []ExploreFields{Beacon}
		concepts := []string{"Cheese"}
		moveTo := &MoveParameters{
			Concepts: []string{"pizza", "pineapple"},
			Force:    0.2,
		}
		moveAwayFrom := &MoveParameters{
			Concepts: []string{"fish"},
			Force:    0.1,
		}

		withNearText := nearTextBuilder.WithConcepts(concepts).
			WithMoveTo(moveTo).
			WithMoveAwayFrom(moveAwayFrom)

		query := builder.WithFields(fields).
			WithNearText(withNearText).
			build()

		expected := `{Explore(nearText:{concepts: ["Cheese"] moveTo: {concepts: ["pizza","pineapple"] force: 0.2} moveAwayFrom: {concepts: ["fish"] force: 0.1} } ){beacon }}`
		assert.Equal(t, expected, query)
	})

	t.Run("Explore with all params", func(t *testing.T) {
		conMock := &MockRunREST{}

		builder := Explore{
			connection: conMock,
		}
		nearTextBuilder := NearTextArgumentBuilder{}

		concepts := []string{"New Yorker"}
		certainty := float32(0.95)
		moveTo := &MoveParameters{
			Concepts: []string{"publisher", "articles"},
			Force:    0.5,
		}
		moveAwayFrom := &MoveParameters{
			Concepts: []string{"fashion", "shop"},
			Force:    0.2,
		}
		fields := []ExploreFields{Beacon, Certainty, ClassName}

		withNearText := nearTextBuilder.
			WithConcepts(concepts).
			WithCertainty(certainty).
			WithMoveAwayFrom(moveAwayFrom).
			WithMoveTo(moveTo)

		query := builder.WithFields(fields).
			WithNearText(withNearText).
			build()

		expected := `{Explore(nearText:{concepts: ["New Yorker"] certainty: 0.95 moveTo: {concepts: ["publisher","articles"] force: 0.5} moveAwayFrom: {concepts: ["fashion","shop"] force: 0.2} } ){beacon certainty className }}`
		assert.Equal(t, expected, query)
	})

	t.Run("Missuse", func(t *testing.T) {
		conMock := &MockRunREST{}
		builder := Explore{
			connection: conMock,
		}
		query := builder.build()
		assert.NotEmpty(t, query, "Check that there is no panic if query is not validly build")
	})

}
