package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate/entities/dto"
)

type MultiTargetArgumentBuilder struct {
	targetCombination *dto.TargetCombination
	targetVectors     []string
}

func (m *MultiTargetArgumentBuilder) getCombinationMethod() string {
	combinationMethod := ""
	if m.targetCombination == nil {
		return combinationMethod
	}
	switch m.targetCombination.Type {
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

func (m *MultiTargetArgumentBuilder) getCombinationWeights() map[string]float32 {
	if m.targetCombination == nil {
		return nil
	}
	return m.targetCombination.Weights
}

func (m *MultiTargetArgumentBuilder) Sum(targetVectors ...string) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = &dto.TargetCombination{Type: dto.Sum}
	}
	return m
}

func (m *MultiTargetArgumentBuilder) Average(targetVectors ...string) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = &dto.TargetCombination{Type: dto.Average}
	}
	return m
}

func (m *MultiTargetArgumentBuilder) Minimum(targetVectors ...string) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = &dto.TargetCombination{Type: dto.Minimum}
	}
	return m
}

func (m *MultiTargetArgumentBuilder) ManualWeights(targetVectors map[string]float32) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		for k := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
		}
		m.targetVectors = targetVectorsTmp
		m.targetCombination = &dto.TargetCombination{Type: dto.ManualWeights, Weights: targetVectors}
	}
	return m
}

func (m *MultiTargetArgumentBuilder) RelativeScore(targetVectors map[string]float32) *MultiTargetArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		for k := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
		}
		m.targetVectors = targetVectorsTmp
		m.targetCombination = &dto.TargetCombination{Type: dto.RelativeScore, Weights: targetVectors}
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
		targetVectorsString := fmt.Sprintf(", targetVectors:%s", string(targetVectorsBytes))

		weightsString := ""
		combinationWeights := m.getCombinationWeights()
		if len(combinationWeights) > 0 {
			weights := make([]string, 0, len(combinationWeights))
			for k, v := range combinationWeights {
				weights = append(weights, fmt.Sprintf("%s:%v", k, v))
			}
			weightsString = fmt.Sprintf(", weights:{%s}", strings.Join(weights, ","))
		}

		combinationMethodString := ""
		combinationMethod := m.getCombinationMethod()
		if combinationMethod != "" {
			combinationMethodString = fmt.Sprintf("combinationMethod:%s", combinationMethod)
		}

		clause = append(clause, fmt.Sprintf(
			"targets:{%s%s%s}",
			combinationMethodString,
			targetVectorsString,
			weightsString,
		))
	}
	return strings.Join(clause, " ")
}
