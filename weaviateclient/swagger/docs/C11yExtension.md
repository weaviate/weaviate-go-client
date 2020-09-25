# C11yExtension

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Concept** | **string** | The new concept you want to extend. Must be an all-lowercase single word, or a space delimited compound word. Examples: &#39;foobarium&#39;, &#39;my custom concept&#39; | [optional] [default to null]
**Definition** | **string** | A list of space-delimited words or a sentence describing what the custom concept is about. Avoid using the custom concept itself. An Example definition for the custom concept &#39;foobarium&#39;: would be &#39;a naturally occourring element which can only be seen by programmers&#39; | [optional] [default to null]
**Weight** | **float32** | Weight of the definition of the new concept where 1&#x3D;&#39;override existing definition entirely&#39; and 0&#x3D;&#39;ignore custom definition&#39;. Note that if the custom concept is not present in the contextionary yet, the weight cannot be less than 1. | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


