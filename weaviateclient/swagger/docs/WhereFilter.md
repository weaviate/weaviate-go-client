# WhereFilter

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Operands** | [**[]WhereFilter**](WhereFilter.md) | combine multiple where filters, requires &#39;And&#39; or &#39;Or&#39; operator | [optional] [default to null]
**Operator** | **string** | operator to use | [optional] [default to null]
**Path** | **[]string** | path to the property currently being filtered | [optional] [default to null]
**ValueInt** | **int64** | value as integer | [optional] [default to null]
**ValueNumber** | **float32** | value as number/float | [optional] [default to null]
**ValueBoolean** | **bool** | value as boolean | [optional] [default to null]
**ValueString** | **string** | value as string | [optional] [default to null]
**ValueText** | **string** | value as text (on text props) | [optional] [default to null]
**ValueDate** | **string** | value as date (as string) | [optional] [default to null]
**ValueGeoRange** | [***WhereFilterGeoRange**](WhereFilterGeoRange.md) | value as geo coordinates and distance | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


