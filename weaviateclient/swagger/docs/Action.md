# Action

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Class** | **string** | Type of the Action, defined in the schema. | [optional] [default to null]
**VectorWeights** | [***VectorWeights**](VectorWeights.md) |  | [optional] [default to null]
**Schema** | [***PropertySchema**](PropertySchema.md) |  | [optional] [default to null]
**Meta** | [***UnderscoreProperties**](UnderscoreProperties.md) |  | [optional] [default to null]
**Id** | **string** | ID of the Action. | [optional] [default to null]
**CreationTimeUnix** | **int64** | Timestamp of creation of this Action in milliseconds since epoch UTC. | [optional] [default to null]
**LastUpdateTimeUnix** | **int64** | Timestamp of the last update made to the Action since epoch UTC. | [optional] [default to null]
**Classification** | [***UnderscorePropertiesClassification**](UnderscorePropertiesClassification.md) | If this object was subject of a classificiation, additional meta info about this classification is available here. (Underscore properties are optional, include them using the ?include&#x3D;_&lt;propName&gt; parameter) | [optional] [default to null]
**Vector** | [***C11yVector**](C11yVector.md) | This object&#39;s position in the Contextionary vector space. (Underscore properties are optional, include them using the ?include&#x3D;_&lt;propName&gt; parameter) | [optional] [default to null]
**Interpretation** | [***Interpretation**](Interpretation.md) | Additional information about how this property was interpreted at vectorization. (Underscore properties are optional, include them using the ?include&#x3D;_&lt;propName&gt; parameter) | [optional] [default to null]
**NearestNeighbors** | [***NearestNeighbors**](NearestNeighbors.md) | Additional information about the neighboring concepts of this element | [optional] [default to null]
**FeatureProjection** | [***FeatureProjection**](FeatureProjection.md) | A feature projection of the object&#39;s vector into lower dimensions for visualization | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


