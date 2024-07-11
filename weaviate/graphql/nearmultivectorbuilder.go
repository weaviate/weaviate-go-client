package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate/entities/dto"
)

type NearMultiVectorArgumentBuilder struct {
	certainty         float32
	distance          float32
	targetCombination dto.TargetCombination
	targetVectors     []string
	vectorPerTarget   map[string][]float32
	withCertainty     bool
	withDistance      bool
}

func (m *NearMultiVectorArgumentBuilder) getCombinationMethod() string {
	combinationMethod := ""
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

func (m *NearMultiVectorArgumentBuilder) Sum(targetVectors ...string) *NearMultiVectorArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Sum}
	}
	return m
}

func (m *NearMultiVectorArgumentBuilder) Average(targetVectors ...string) *NearMultiVectorArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Average}
	}
	return m
}

func (m *NearMultiVectorArgumentBuilder) Min(targetVectors ...string) *NearMultiVectorArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Minimum}
	}
	return m
}

func (m *NearMultiVectorArgumentBuilder) ManualWeights(targetVectors map[string]float32) *NearMultiVectorArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		for k := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
		}
		m.targetVectors = targetVectorsTmp
		m.targetCombination = dto.TargetCombination{Type: dto.ManualWeights, Weights: targetVectors}
	}
	return m
}

func (m *NearMultiVectorArgumentBuilder) RelativeScore(targetVectors map[string]float32) *NearMultiVectorArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		for k := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
		}
		m.targetVectors = targetVectorsTmp
		m.targetCombination = dto.TargetCombination{Type: dto.RelativeScore, Weights: targetVectors}
	}
	return m
}

func (m *NearMultiVectorArgumentBuilder) WithVectorPerTarget(vectorPerTarget map[string][]float32) *NearMultiVectorArgumentBuilder {
	if len(vectorPerTarget) > 0 {
		m.vectorPerTarget = vectorPerTarget
	}
	return m
}

func (m *NearMultiVectorArgumentBuilder) WithVector(vector []float32) *NearMultiVectorArgumentBuilder {
	vectorPerTarget := make(map[string][]float32, len(m.targetVectors))
	for _, target := range m.targetVectors {
		vectorPerTarget[target] = vector
	}
	m.vectorPerTarget = vectorPerTarget
	return m
}

func (m *NearMultiVectorArgumentBuilder) WithCertainty(certainty float32) *NearMultiVectorArgumentBuilder {
	m.certainty = certainty
	m.withCertainty = true
	return m
}

func (m *NearMultiVectorArgumentBuilder) WithDistance(distance float32) *NearMultiVectorArgumentBuilder {
	m.distance = distance
	m.withDistance = true
	return m
}

func (m *NearMultiVectorArgumentBuilder) build() string {
	clause := []string{}
	targetVectors := m.targetVectors
	if m.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty:%v", m.certainty))
	}
	if m.withDistance {
		clause = append(clause, fmt.Sprintf("distance:%v", m.distance))
	}
	if len(m.vectorPerTarget) > 0 {
		vectorPerTarget := make([]string, 0, len(m.vectorPerTarget))
		for k, v := range m.vectorPerTarget {
			vBytes, err := json.Marshal(v)
			if err != nil {
				panic(fmt.Sprintf("could not marshal vector: %v", err))
			}
			vectorPerTarget = append(vectorPerTarget, fmt.Sprintf("%s:%v", k, string(vBytes)))
		}
		clause = append(clause, fmt.Sprintf("vectorPerTarget:{%s}", strings.Join(vectorPerTarget, ",")))
		if len(targetVectors) == 0 {
			targetVectors = make([]string, 0, len(m.vectorPerTarget))
			for k := range m.vectorPerTarget {
				targetVectors = append(targetVectors, k)
			}
		}
	}
	if len(targetVectors) > 0 {
		targetVectorsBytes, err := json.Marshal(targetVectors)
		if err != nil {
			panic(fmt.Sprintf("could not marshal target vectors: %v", err))
		}
		targetVectorsString := fmt.Sprintf("targetVectors:%s", string(targetVectorsBytes))

		weightsString := ""
		combinationWeights := m.targetCombination.Weights
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
			combinationMethodString = fmt.Sprintf(", combinationMethod:%s", combinationMethod)
		}

		clause = append(clause, fmt.Sprintf(
			"targets:{%s%s%s}",
			targetVectorsString,
			combinationMethodString,
			weightsString,
		))
	}
	return fmt.Sprintf("nearVector:{%s}", strings.Join(clause, " "))
}
