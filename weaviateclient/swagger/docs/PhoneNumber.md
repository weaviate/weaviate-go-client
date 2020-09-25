# PhoneNumber

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Input** | **string** | The raw input as the phone number is present in your raw data set. It will be parsed into the standardized formats if valid. | [optional] [default to null]
**InternationalFormatted** | **string** | Read-only. Parsed result in the international format (e.g. +49 123 ...) | [optional] [default to null]
**DefaultCountry** | **string** | Optional. The ISO 3166-1 alpha-2 country code. This is used to figure out the correct countryCode and international format if only a national number (e.g. 0123 4567) is provided | [optional] [default to null]
**CountryCode** | **float32** | Read-only. The numerical country code (e.g. 49) | [optional] [default to null]
**National** | **float32** | Read-only. The numerical representation of the national part | [optional] [default to null]
**NationalFormatted** | **string** | Read-only. Parsed result in the national format (e.g. 0123 456789) | [optional] [default to null]
**Valid** | **bool** | Read-only. Indicates whether the parsed number is a valid phone number | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


