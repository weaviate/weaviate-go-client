# Weaviate Go Client v6 - Design Decisions

**Audience:** Internal stakeholders
**Purpose:** High-level overview of key architectural decisions through examples

## Overview

Goals: a native Go experience with:
- **Explicit, readable APIs**
- **Minimal boilerplate**
- **Type safety**

Here are some concrete examples illustrating the design choices.

---

## Single insert

```go
songs := client.Collections.Use("Songs")
songs.Data.Insert(ctx,
    data.WithProperties(Song{
        Title:  "Bohemian Rhapsody",
        Artist: "Queen",
        Year:   1975,
    }),
    data.WithVector(types.Vector{Single: []float32{0.1, 0.2, 0.3}}),
)
```

- **Functional, variadic options**
- **Flexible properties**: Accept a struct or `map[string]any`
- **Unified Vector type**: Single `Vector` type for both 1D (`Single`) and 2D (`Multi`) vectors

---

## NearVector query

```go
queryVector := types.Vector{Single: []float32{0.1, 0.2, 0.3}}
result, err := songs.Query.NearVector(ctx, queryVector,
    query.WithLimit(10),
    query.WithDistance(0.5),
    query.WithOffset(20),
)

// Access results
for _, obj := range result.Objects {
    fmt.Println(obj.UUID)                     // Direct access
    fmt.Println(obj.Vectors["text"].Single)   // Vectors always accessible
    title := obj.Properties["title"].(string) // Map-based by default
}
```

- **Variadic options**: Similar to insert
- **Map-based objects**: Default object returns are `map[string]any` for convenient prototyping
- **UUID and Vector fields**

---

## Opt-In generics

```go
result, err := songs.Query.NearVector(ctx, queryVector)
// result.Objects -> []WeaviateObject[map[string]any]

type Song struct {
    Title  string `json:"title"`
    Artist string `json:"artist"`
    Year   int    `json:"year"`
}

// Use `Scan` to convert to typed objects
typedObjects := query.Scan[Song](result)

for _, obj := range typedObjects {
    fmt.Println(obj.Properties.Title)  // No type assertion!
    fmt.Println(obj.Properties.Artist) // IDE autocomplete works
    fmt.Println(obj.UUID)
}
```

- **Additional function for type safety**: `Scan[T]()` converts map-based results to typed structs
    - Simplicity with optional type safety

---

## GroupBy queries

```go
// Standard vs grouped NearVector queries - function-as-receiver pattern
single, err := songs.Query.NearVector(ctx, vector, query.WithLimit(10))
groups, err := songs.Query.NearVector.GroupBy(ctx, vector, "category", query.WithLimit(10))
```

- **Function-as-receiver**: More common query prioritized, less common GroupBy as method
- **Shared options with different return types**

---

## Multiple Vector Formats

```go
// Default vector (unnamed)
songs.Data.Insert(ctx,
    data.WithProperties(data),
    data.WithVector(types.Vector{Single: []float32{0.1, 0.2, 0.3}}),
)

// Named single-dimensional vector
songs.Data.Insert(ctx,
    data.WithProperties(data),
    data.WithVector(types.Vector{
        Name:   "text_embedding",
        Single: []float32{0.1, 0.2, 0.3},
    }),
)

// Named multi-dimensional vector
songs.Data.Insert(ctx,
    data.WithProperties(imageData),
    data.WithVector(types.Vector{
        Name:  "colbert",
        Multi: [][]float32{
          {0.1, 0.2}, {0.3, 0.4}
        },
    }),
)

// Multiple vectors at once
songs.Data.Insert(ctx,
    data.WithProperties(data),
    data.WithVector([]types.Vector{
        {Name: "single_vec", Single: []float32{0.1, 0.2, 0.3}},
        {Name: "matrix_vec", Multi: [][]float32{
            {0.4, 0.5}, {0.6, 0.7}
        }},
    }
)
```

---

## Multi-Vector search

```go
singleVec := types.Vector{Name: "single_vec", Single: []float32{0.1, 0.2, 0.3}}
matrixVec := types.Vector{Name: "matrix_vec", Single: []float32{0.4, 0.5, 0.6}}

result, err := songs.Query.NearVector(ctx,
    query.Average(singleVec, matrixVec),
    query.WithLimit(10),
)

result, err := songs.Query.NearVector(ctx,
    query.ManualWeights(
        query.Target(singleVec, 0.7),
        query.Target(matrixVec, 0.3),
    ),
    query.WithLimit(10),
)
```

---

## Reuse Vectors from results

```go
result, err := songs.Query.NearVector(ctx, queryVector, query.WithLimit(10))

// Re-insert using returned vectors
for _, obj := range result.Objects {
    newID, err := songs.Data.Insert(ctx,
        data.WithProperties(obj.Properties),  // Reuse properties map
        data.WithVector(obj.Vectors),         // Reuse entire vector map
    )
}
```

---

### Pointers for optional values

Implication: Use pointers for optional fields in structs to distinguish between zero values and unset fields.
