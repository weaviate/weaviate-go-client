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
type NearMultiVectorTargettedArgumentBuilder struct {
	base *NearMultiVectorArgumentBuilder
}
type targets struct {
	combinationMethod string
	targetVectors     []string
	weights           map[string]float32
}

func (m *NearMultiVectorArgumentBuilder) toTargets() *targets {
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
	return &targets{
		combinationMethod: combinationMethod,
		targetVectors:     m.targetVectors,
		weights:           m.targetCombination.Weights,
	}
}

func (m *NearMultiVectorArgumentBuilder) Sum(targetVectors ...string) *NearMultiVectorTargettedArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Sum}
	}
	return &NearMultiVectorTargettedArgumentBuilder{m}
}

func (m *NearMultiVectorArgumentBuilder) Average(targetVectors ...string) *NearMultiVectorTargettedArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Average}
	}
	return &NearMultiVectorTargettedArgumentBuilder{m}
}

func (m *NearMultiVectorArgumentBuilder) Min(targetVectors ...string) *NearMultiVectorTargettedArgumentBuilder {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Minimum}
	}
	return &NearMultiVectorTargettedArgumentBuilder{m}
}

func (m *NearMultiVectorArgumentBuilder) ManualWeights(targetVectors map[string]float32) *NearMultiVectorTargettedArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		for k := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
		}
		m.targetVectors = targetVectorsTmp
		m.targetCombination = dto.TargetCombination{Type: dto.ManualWeights, Weights: targetVectors}
	}
	return &NearMultiVectorTargettedArgumentBuilder{m}
}

func (m *NearMultiVectorArgumentBuilder) RelativeScore(targetVectors map[string]float32) *NearMultiVectorTargettedArgumentBuilder {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		for k := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
		}
		m.targetVectors = targetVectorsTmp
		m.targetCombination = dto.TargetCombination{Type: dto.RelativeScore, Weights: targetVectors}
	}
	return &NearMultiVectorTargettedArgumentBuilder{m}
}

func (m *NearMultiVectorTargettedArgumentBuilder) WithVectorPerTarget(vectorPerTarget map[string][]float32) *NearMultiVectorArgumentBuilder {
	if len(vectorPerTarget) > 0 {
		m.base.vectorPerTarget = vectorPerTarget
	}
	return m.base
}

func (m *NearMultiVectorTargettedArgumentBuilder) WithVector(vector []float32) *NearMultiVectorArgumentBuilder {
	vectorPerTarget := make(map[string][]float32, len(m.base.targetVectors))
	for _, target := range m.base.targetVectors {
		vectorPerTarget[target] = vector
	}
	m.base.vectorPerTarget = vectorPerTarget
	return m.base
}

func (m *NearMultiVectorTargettedArgumentBuilder) WithCertainty(certainty float32) *NearMultiVectorTargettedArgumentBuilder {
	m.base.certainty = certainty
	m.base.withCertainty = true
	return m
}

func (m *NearMultiVectorTargettedArgumentBuilder) WithDistance(distance float32) *NearMultiVectorTargettedArgumentBuilder {
	m.base.distance = distance
	m.base.withDistance = true
	return m
}

func (m *NearMultiVectorArgumentBuilder) build() string {
	clause := []string{}
	if m.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", m.certainty))
	}
	if m.withDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", m.distance))
	}
	if len(m.vectorPerTarget) > 0 {
		vectorPerTarget, err := json.Marshal(m.vectorPerTarget)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal near multi vector search vector: %s", err))
		}
		clause = append(clause, fmt.Sprintf("vectorPerTarget: %s", string(vectorPerTarget)))
	}
	if len(m.targetVectors) > 0 {
		targets := m.toTargets()
		weights, err := json.Marshal(targets.weights)
		if err != nil {
			panic(fmt.Errorf("failed to unmarshal near multi vector search weights: %s", err))
		}
		clause = append(clause, fmt.Sprintf("targets:{combinationMethod: %s, targetVectors: %s, weights: %s}", targets.combinationMethod, targets.targetVectors, string(weights)))
	}
	return fmt.Sprintf("nearVector:{%v}", strings.Join(clause, " "))
}
