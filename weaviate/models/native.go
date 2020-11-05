package models

import (
	native "github.com/semi-technologies/weaviate/entities/models"
)

type Action native.Action
type Thing native.Thing
type ActionsGetResponse native.ActionsGetResponse
type Schema native.Schema
type BatchReference native.BatchReference
type BatchReferenceResponse native.BatchReferenceResponse
type ThingsGetResponse native.ThingsGetResponse
type Classification native.Classification
type WhereFilter native.WhereFilter
type C11yWordsResponse native.C11yWordsResponse
type C11yExtension native.C11yExtension
type ActionsListResponse native.ActionsListResponse
type PropertySchema native.PropertySchema
type SingleRef native.SingleRef
type MultipleRef native.MultipleRef
type ThingsListResponse native.ThingsListResponse
type GraphQLResponse native.GraphQLResponse
type GraphQLQuery native.GraphQLQuery
type Meta native.Meta
type Class native.Class
type Property native.Property

// "github.com/semi-technologies/weaviate-go-client/weaviate/models"

func CastActionsFromActionsListResponse(actionsListResponse *ActionsListResponse) []*Action {
	casted := make([]*Action, len(actionsListResponse.Actions))
	for i, action := range actionsListResponse.Actions {
		newAction := Action(*action)
		casted[i] = &newAction
	}
	return casted
}

func CastThingsFromThingsListResponse(thingsListRsponse *ThingsListResponse) []*Thing {
	casted := make([]*Thing, len(thingsListRsponse.Things))
	for i, thing := range thingsListRsponse.Things {
		newThing := Thing(*thing)
		casted[i] = &newThing
	}
	return casted
}

func CastFromNativeWhereFilter(filter *native.WhereFilter) *WhereFilter {
	newFilter := WhereFilter(*filter)
	return &newFilter
}

func CastToNativeWhereFilter(filter *WhereFilter) *native.WhereFilter {
	newFilter := native.WhereFilter(*filter)
	return &newFilter
}