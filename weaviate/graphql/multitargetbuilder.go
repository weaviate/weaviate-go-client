package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/weaviate/weaviate/entities/dto"
	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

type MultiTargetArgumentBuilder struct {
	targetCombination *dto.TargetCombinationType
	targetVectors     []string
	weights           [][]float32
}

func (m *MultiTargetArgumentBuilder) getCombinationMethod() (string, pb.CombinationMethod) {
	if m.targetCombination != nil {
		switch *m.targetCombination {
		case dto.Sum:
			return "sum", pb.CombinationMethod_COMBINATION_METHOD_TYPE_SUM
		case dto.Average:
			return "average", pb.CombinationMethod_COMBINATION_METHOD_TYPE_AVERAGE
		case dto.Minimum:
			return "minimum", pb.CombinationMethod_COMBINATION_METHOD_TYPE_MIN
		case dto.ManualWeights:
			return "manualWeights", pb.CombinationMethod_COMBINATION_METHOD_TYPE_MANUAL
		case dto.RelativeScore:
			return "relativeScore", pb.CombinationMethod_COMBINATION_METHOD_TYPE_RELATIVE_SCORE
		}
	}
	return "", pb.CombinationMethod_COMBINATION_METHOD_UNSPECIFIED
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
		combinationMethod, _ := m.getCombinationMethod()
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

func (m *MultiTargetArgumentBuilder) togrpc() *pb.Targets {
	if len(m.targetVectors) > 0 {
		var weightsForTargets []*pb.WeightsForTarget
		for i, weights := range m.weights {
			for _, w := range weights {
				weightsForTargets = append(weightsForTargets, &pb.WeightsForTarget{
					Target: m.targetVectors[i],
					Weight: w,
				})
			}
		}
		_, combination := m.getCombinationMethod()
		targets := &pb.Targets{
			TargetVectors:     m.targetVectors,
			Combination:       combination,
			WeightsForTargets: weightsForTargets,
		}
		return targets
	}
	return nil
}
