package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate/entities/dto"
)

type MultiTargetArgumentBuilder struct {
	targetCombination *dto.TargetCombinationType
	targetVectors     []string
	weights           [][]float32
}

func (m *MultiTargetArgumentBuilder) getCombinationMethod() string {
	combinationMethod := ""
	if m.targetCombination == nil {
		return combinationMethod
	}
	switch *m.targetCombination {
	case dto.Sum:
		combinationMethod = "sum"
	case dto.Average:
		combinationMethod = "average"
	case dto.Minimum:
		combinationMethod = "minimum"
	case dto.ManualWeights:
		combinationMethod = "manualWeights"
	case dto.RelativeScore:
		combinationMethod = "relativeScore"
	}
	return combinationMethod
}

func (m *MultiTargetArgumentBuilder) Sum(targetVectors ...string) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		comb := dto.Sum
		m.targetCombination = &comb
	}
	return m
}

func (m *MultiTargetArgumentBuilder) Average(targetVectors ...string) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		comb := dto.Average
		m.targetCombination = &comb
	}
	return m
}

func (m *MultiTargetArgumentBuilder) Minimum(targetVectors ...string) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		comb := dto.Minimum
		m.targetCombination = &comb
	}
	return m
}

func (m *MultiTargetArgumentBuilder) ManualWeights(targetVectors map[string]float32) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		weightsTmp := make([][]float32, 0, len(targetVectors))
		for k, v := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
			weightsTmp = append(weightsTmp, []float32{v})
		}
		m.targetVectors = targetVectorsTmp
		m.weights = weightsTmp
		comb := dto.ManualWeights
		m.targetCombination = &comb
	}
	return m
}

func (m *MultiTargetArgumentBuilder) ManualWeightsMulti(targetVectors map[string][]float32) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		weightsTmp := make([][]float32, 0, len(targetVectors))
		for k, vv := range targetVectors {
			for range vv {
				targetVectorsTmp = append(targetVectorsTmp, k)
			}
			weightsTmp = append(weightsTmp, vv)
		}
		m.weights = weightsTmp
		m.targetVectors = targetVectorsTmp
		comb := dto.ManualWeights
		m.targetCombination = &comb
	}
	return m
}

func (m *MultiTargetArgumentBuilder) RelativeScore(targetVectors map[string]float32) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		weightsTmp := make([][]float32, 0, len(targetVectors))
		for k, v := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
			weightsTmp = append(weightsTmp, []float32{v})
		}
		m.targetVectors = targetVectorsTmp
		m.weights = weightsTmp
		comb := dto.RelativeScore
		m.targetCombination = &comb
	}
	return m
}

func (m *MultiTargetArgumentBuilder) RelativeScoreMulti(targetVectors map[string][]float32) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		weightsTmp := make([][]float32, 0, len(targetVectors))
		for k, vv := range targetVectors {
			for range vv {
				targetVectorsTmp = append(targetVectorsTmp, k)
			}
			weightsTmp = append(weightsTmp, vv)
		}
		m.weights = weightsTmp
		m.targetVectors = targetVectorsTmp
		comb := dto.RelativeScore
		m.targetCombination = &comb
	}
	return m
}

func (m *MultiTargetArgumentBuilder) build() string {
	clause := []string{}
	targetVectors := m.targetVectors

	if len(targetVectors) > 0 {
		targetVectorsBytes, err := json.Marshal(targetVectors)
		if err != nil {
			panic(fmt.Sprintf("could not marshal target vectors: %v", err))
		}
		targetVectorsString := fmt.Sprintf(", targetVectors: %s", string(targetVectorsBytes))

		weightsString := ""
		if len(m.weights) > 0 {
			weights := make([]string, 0, len(m.weights))
			targetCount := 0
			for _, v := range m.weights {
				if len(v) > 1 {
					vectorStr := fmt.Sprintf("%v", v[0])
					for _, vv := range v[1:] {
						vectorStr += fmt.Sprintf(",%v", vv)
					}
					weights = append(weights, fmt.Sprintf("%s: [%v]", targetVectors[targetCount], vectorStr))
				} else {
					weights = append(weights, fmt.Sprintf("%s: %v", targetVectors[targetCount], v[0]))
				}
				targetCount += len(v)
			}
			weightsString = fmt.Sprintf(", weights: {%s}", strings.Join(weights, ","))
		}

		combinationMethodString := ""
		combinationMethod := m.getCombinationMethod()
		if combinationMethod != "" {
			combinationMethodString = fmt.Sprintf("combinationMethod: %s", combinationMethod)
		}

		clause = append(clause, fmt.Sprintf(
			"%s%s%s",
			combinationMethodString,
			targetVectorsString,
			weightsString,
		))
	}
	return strings.Join(clause, " ")
}
