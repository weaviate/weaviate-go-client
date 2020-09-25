# UnderscorePropertiesClassification

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | unique identifier of the classification run | [optional] [default to null]
**Completed** | [**time.Time**](time.Time.md) | Timestamp when this particular object was classified. This is usually sooner than the overall completion time of the classification, as the overall completion time will only be set once every object has been classified. | [optional] [default to null]
**Scope** | **[]string** | The properties in scope of the classification. Note that this doesn&#39;t mean that these fields were necessarily classified, this only means that those fields were in scope of the classificiation. See \&quot;classifiedFields\&quot; for details. | [optional] [default to null]
**ClassifiedFields** | **[]string** | The (reference) fields which were classified as part of this classification. Note that this might contain fewere entries than \&quot;scope\&quot;, if one of the fields was already set prior to the classification, for example | [optional] [default to null]
**BasedOn** | **[]string** | The (primitive) field(s) which were used as a basis for classification. For example, if the type of classification is \&quot;knn\&quot; with k&#x3D;3, the 3 nearest neighbors - based on these fields - were considered for the classification. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


