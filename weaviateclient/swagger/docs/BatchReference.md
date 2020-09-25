# BatchReference

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**From** | **string** | Long-form beacon-style URI to identify the source of the cross-ref including the property name. Should be in the form of weaviate://localhost/&lt;kinds&gt;/&lt;uuid&gt;/&lt;className&gt;/&lt;propertyName&gt;, where &lt;kinds&gt; must be one of &#39;actions&#39;, &#39;things&#39; and &lt;className&gt; and &lt;propertyName&gt; must represent the cross-ref property of source class to be used. | [optional] [default to null]
**To** | **string** | Short-form URI to point to the cross-ref. Should be in the form of weaviate://localhost/things/&lt;uuid&gt; for the example of a local cross-ref to a thing | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


