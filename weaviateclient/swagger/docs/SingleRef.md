# SingleRef

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Class** | **string** | If using a concept reference (rather than a direct reference), specify the desired class name here | [optional] [default to null]
**Schema** | [***PropertySchema**](PropertySchema.md) | If using a concept reference (rather than a direct reference), specify the desired properties here | [optional] [default to null]
**Beacon** | **string** | If using a direct reference, specify the URI to point to the cross-ref here. Should be in the form of weaviate://localhost/things/&lt;uuid&gt; for the example of a local cross-ref to a thing | [optional] [default to null]
**Href** | **string** | If using a direct reference, this read-only fields provides a link to the refernced resource. If &#39;origin&#39; is globally configured, an absolute URI is shown - a relative URI otherwise. | [optional] [default to null]
**Meta** | [***ReferenceMeta**](ReferenceMeta.md) | Additional Meta information about this particular reference. Only shown if meta&#x3D;&#x3D;true. | [optional] [default to null]
**Classification** | [***ReferenceMetaClassification**](ReferenceMetaClassification.md) | Additional Meta information about classifications if the item was part of one | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


