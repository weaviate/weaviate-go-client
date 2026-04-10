package query

type (
	// KeywordSimilarity conrols the similarity threshold for BM25 (keyword) search.
	KeywordSimilarity struct {
		// allTokensMatch requires each token in the search query to be present
		// in a candidate object for it to be considered a match.
		allTokensMatch bool

		// minimumTokensMatch is the lower threshold for the number of times
		// _each_ token needs be present in a candidate object for it to be considered a match.
		mininumTokensMatch *int32
	}
)

func (kws *KeywordSimilarity) AllTokensMatch() bool       { return kws.allTokensMatch }
func (kws *KeywordSimilarity) MinimumTokensMatch() *int32 { return kws.mininumTokensMatch }

// AllTokensMatch is a [KeywordSimilarity] parameter with AllTokensMatch=true.
var AllTokensMatch = KeywordSimilarity{allTokensMatch: true}

// MinimumTokensMatch returns [KeywordSimilarity] with MinimumTokensMatch=n.
func MinimumTokensMatch(n int32) KeywordSimilarity {
	return KeywordSimilarity{mininumTokensMatch: &n}
}
