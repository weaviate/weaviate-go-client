# Deprecation

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | **string** | The id that uniquely identifies this particular deprecations (mostly used internally) | [optional] [default to null]
**Status** | **string** | Whether the problematic API functionality is deprecated (planned to be removed) or already removed | [optional] [default to null]
**ApiType** | **string** | Describes which API is effected, usually one of: REST, GraphQL | [optional] [default to null]
**Msg** | **string** | What this deprecation is about | [optional] [default to null]
**Mitigation** | **string** | User-required action to not be affected by the (planned) removal | [optional] [default to null]
**SinceVersion** | **string** | The deprecation was introduced in this version | [optional] [default to null]
**PlannedRemovalVersion** | **string** | A best-effort guess of which upcoming version will remove the feature entirely | [optional] [default to null]
**RemovedIn** | **string** | If the feature has already been removed, it was removed in this version | [optional] [default to null]
**RemovedTime** | [**time.Time**](time.Time.md) | If the feature has already been removed, it was removed at this timestamp | [optional] [default to null]
**SinceTime** | [**time.Time**](time.Time.md) | The deprecation was introduced in this version | [optional] [default to null]
**Locations** | **[]string** | The locations within the specified API affected by this deprecation | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


