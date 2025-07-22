package graphql

import (
	"encoding/json"
	"fmt"
	"strings"

	pb "github.com/weaviate/weaviate/grpc/generated/protocol/v1"
)

// fldMover is a type representing field names of a move sub query
type fldMover string

const (
	fldMoverConcepts fldMover = "concepts"
	fldMoverForce    fldMover = "force"
	fldMoverID       fldMover = "id"
	fldMoverBeacon   fldMover = "beacon"
	fldMoverObjects  fldMover = "objects"
)

// MoveParameters to fine tune Explore queries
type MoveParameters struct {
	// Concepts that should be used as base for the movement operation
	Concepts []string
	// Force to be applied in the movement operation
	Force float32
	// Objects used to adjust the serach direction
	Objects []MoverObject
}

func (m *MoveParameters) String() string {
	concepts, _ := json.Marshal(m.Concepts)
	ms := make([]string, 0, len(m.Objects))
	for _, m := range m.Objects {
		if s := m.String(); s != EmptyObjectStr {
			ms = append(ms, s)
		}
	}
	if len(ms) < 1 {
		return fmt.Sprintf("{%s: %s %s: %v}", fldMoverConcepts, concepts, fldMoverForce, m.Force)
	}

	s := "{"
	if len(m.Concepts) > 0 {
		s = fmt.Sprintf("{%s: %s", fldMoverConcepts, concepts)
	}
	return fmt.Sprintf("%s %s: %v %s: %v}", s, fldMoverForce, m.Force, fldMoverObjects, ms)
}

// MoverObject is the object the search is supposed to move close to (or further away from) it.
type MoverObject struct {
	ID     string
	Beacon string
}

// String returns string representation of m as {"id": "value" beacon:"value"}.
// Empty fields are considered optional and are excluded.
// It returns EmptyObjectStr if both fields are empty
func (m *MoverObject) String() string {
	if m.ID != "" && m.Beacon != "" {
		return fmt.Sprintf(`{%s: "%s" %s: "%s"}`, fldMoverID, m.ID, fldMoverBeacon, m.Beacon)
	}
	if m.ID != "" {
		return fmt.Sprintf(`{%s: "%s"}`, fldMoverID, m.ID)
	}
	if m.Beacon != "" {
		return fmt.Sprintf(`{%s: "%s"}`, fldMoverBeacon, m.Beacon)
	}
	return EmptyObjectStr
}

type NearTextArgumentBuilder struct {
	concepts        []string
	withCertainty   bool
	certainty       float32
	withDistance    bool
	distance        float32
	moveTo          *MoveParameters
	moveAwayFrom    *MoveParameters
	withAutocorrect bool
	autocorrect     bool
	targetVectors   []string
	targets         *MultiTargetArgumentBuilder
}

// WithConcepts the result is based on
func (e *NearTextArgumentBuilder) WithConcepts(concepts []string) *NearTextArgumentBuilder {
	e.concepts = concepts
	return e
}

// WithCertainty that is minimally required for an object to be included in the result set
func (e *NearTextArgumentBuilder) WithCertainty(certainty float32) *NearTextArgumentBuilder {
	e.withCertainty = true
	e.certainty = certainty
	return e
}

// WithDistance that is minimally required for an object to be included in the result set
func (e *NearTextArgumentBuilder) WithDistance(distance float32) *NearTextArgumentBuilder {
	e.withDistance = true
	e.distance = distance
	return e
}

// WithMoveTo specific concept
func (e *NearTextArgumentBuilder) WithMoveTo(parameters *MoveParameters) *NearTextArgumentBuilder {
	e.moveTo = parameters
	return e
}

// WithMoveAwayFrom specific concept
func (e *NearTextArgumentBuilder) WithMoveAwayFrom(parameters *MoveParameters) *NearTextArgumentBuilder {
	e.moveAwayFrom = parameters
	return e
}

// WithAutocorrect this is a setting enabling autocorrect of the concepts texts
func (e *NearTextArgumentBuilder) WithAutocorrect(autocorrect bool) *NearTextArgumentBuilder {
	e.withAutocorrect = true
	e.autocorrect = autocorrect
	return e
}

// WithTargetVectors target vector name
func (e *NearTextArgumentBuilder) WithTargetVectors(targetVectors ...string) *NearTextArgumentBuilder {
	if len(targetVectors) > 0 {
		e.targetVectors = targetVectors
	}
	return e
}

// WithTargets sets the multi target vectors to be used with hybrid query. This builder takes precedence over WithTargetVectors.
// So if WithTargets is used, WithTargetVectors will be ignored.
func (h *NearTextArgumentBuilder) WithTargets(targets *MultiTargetArgumentBuilder) *NearTextArgumentBuilder {
	h.targets = targets
	return h
}

// Build build the given clause
func (e *NearTextArgumentBuilder) build() string {
	clause := []string{}
	concepts, _ := json.Marshal(e.concepts)

	clause = append(clause, fmt.Sprintf("concepts: %s", concepts))
	if e.withCertainty {
		clause = append(clause, fmt.Sprintf("certainty: %v", e.certainty))
	}
	if e.withDistance {
		clause = append(clause, fmt.Sprintf("distance: %v", e.distance))
	}
	if e.moveTo != nil {
		clause = append(clause, fmt.Sprintf("moveTo: %s", e.moveTo))
	}
	if e.moveAwayFrom != nil {
		clause = append(clause, fmt.Sprintf("moveAwayFrom: %s", e.moveAwayFrom))
	}
	if e.withAutocorrect {
		clause = append(clause, fmt.Sprintf("autocorrect: %v", e.autocorrect))
	}
	if e.targets != nil {
		clause = append(clause, fmt.Sprintf("targets:{%s}", e.targets.build()))
	}
	if len(e.targetVectors) > 0 && e.targets == nil {
		targetVectors, _ := json.Marshal(e.targetVectors)
		clause = append(clause, fmt.Sprintf("targetVectors: %s", targetVectors))
	}
	return fmt.Sprintf("nearText:{%v}", strings.Join(clause, " "))
}

func (e *NearTextArgumentBuilder) togrpc() *pb.NearTextSearch {
	nearText := &pb.NearTextSearch{
		Query: e.concepts,
	}
	if e.withCertainty {
		certainty := float64(e.certainty)
		nearText.Certainty = &certainty
	}
	if e.withDistance {
		distance := float64(e.distance)
		nearText.Distance = &distance
	}
	if e.moveTo != nil {
		nearText.MoveTo = e.buildMoveParam(e.moveTo)
	}
	if e.moveAwayFrom != nil {
		nearText.MoveAway = e.buildMoveParam(e.moveAwayFrom)
	}
	if e.targets != nil {
		nearText.Targets = e.targets.togrpc()
	}
	if len(e.targetVectors) > 0 && e.targets == nil {
		nearText.Targets = &pb.Targets{TargetVectors: e.targetVectors}
	}
	return nearText
}

func (e *NearTextArgumentBuilder) buildMoveParam(moveParam *MoveParameters) *pb.NearTextSearch_Move {
	move := &pb.NearTextSearch_Move{
		Concepts: moveParam.Concepts,
	}
	if moveParam.Force != 0.0 {
		move.Force = moveParam.Force
	}
	if len(moveParam.Objects) > 0 {
		uuids := make([]string, len(moveParam.Objects))
		for i := range moveParam.Objects {
			// TODO: handle beacon
			uuids[i] = moveParam.Objects[i].ID
		}
		move.Uuids = uuids
	}
	return move
}
