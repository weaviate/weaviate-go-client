# Property

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DataType** | **[]string** | Can be a reference to another type when it starts with a capital (for example Person), otherwise \&quot;string\&quot; or \&quot;int\&quot;. | [optional] [default to null]
**Cardinality** | **string** | DEPRECATED - do not use anymore. | [optional] [default to null]
**Description** | **string** | Description of the property. | [optional] [default to null]
**VectorizePropertyName** | **bool** | Set this to true if the object vector should include this property&#39;s name in calculating the overall vector position. If set to false (default), only the property value will be used. | [optional] [default to null]
**Name** | **string** | Name of the property as URI relative to the schema URL. | [optional] [default to null]
**Keywords** | [***Keywords**](Keywords.md) |  | [optional] [default to null]
**Index** | **bool** | Optional. By default each property is fully indexed both for full-text, as well as vector-search. You can ignore properties in searches by explicitly setting index to false. Not set is the same as true | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


