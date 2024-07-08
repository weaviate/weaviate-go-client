package graphql

import "github.com/weaviate/weaviate/entities/dto"

type NearVectorMultiTargetArgBuilder struct {
	targetCombination dto.TargetCombination
	targetVectors     []string
	vectorPerTarget   map[string][]float32
}
type NearVectorMultiTargetArgBuilder2 struct {
	base *NearVectorMultiTargetArgBuilder
}

func (m *NearVectorMultiTargetArgBuilder) Sum(targetVectors ...string) *NearVectorMultiTargetArgBuilder2 {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Sum}
	}
	return &NearVectorMultiTargetArgBuilder2{m}
}

func (m *NearVectorMultiTargetArgBuilder) Average(targetVectors ...string) *NearVectorMultiTargetArgBuilder2 {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Average}
	}
	return &NearVectorMultiTargetArgBuilder2{m}
}

func (m *NearVectorMultiTargetArgBuilder) Min(targetVectors ...string) *NearVectorMultiTargetArgBuilder2 {
	if len(targetVectors) > 0 {
		m.targetVectors = targetVectors
		m.targetCombination = dto.TargetCombination{Type: dto.Minimum}
	}
	return &NearVectorMultiTargetArgBuilder2{m}
}

func (m *NearVectorMultiTargetArgBuilder) ManualWeights(targetVectors map[string]float32) *NearVectorMultiTargetArgBuilder2 {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		for k := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
		}
		m.targetVectors = targetVectorsTmp
		m.targetCombination = dto.TargetCombination{Type: dto.ManualWeights, Weights: targetVectors}
	}
	return &NearVectorMultiTargetArgBuilder2{m}
}

func (m *NearVectorMultiTargetArgBuilder) RelativeScore(targetVectors map[string]float32) *NearVectorMultiTargetArgBuilder2 {
	if len(targetVectors) > 0 {
		targetVectorsTmp := make([]string, 0, len(targetVectors))
		for k := range targetVectors {
			targetVectorsTmp = append(targetVectorsTmp, k)
		}
		m.targetVectors = targetVectorsTmp
		m.targetCombination = dto.TargetCombination{Type: dto.RelativeScore, Weights: targetVectors}
	}
	return &NearVectorMultiTargetArgBuilder2{m}
}

func (m *NearVectorMultiTargetArgBuilder2) WithVectorPerTarget(vectorPerTarget map[string][]float32) *NearVectorMultiTargetArgBuilder {
	if len(vectorPerTarget) > 0 {
		m.base.vectorPerTarget = vectorPerTarget
	}
	return m.base
}

func (m *NearVectorMultiTargetArgBuilder2) WithVector(vector []float32) *NearVectorMultiTargetArgBuilder {
	vectorPerTarget := make(map[string][]float32, len(m.base.targetVectors))
	for _, target := range m.base.targetVectors {
		vectorPerTarget[target] = vector
	}
	m.base.vectorPerTarget = vectorPerTarget
	return m.base
}
