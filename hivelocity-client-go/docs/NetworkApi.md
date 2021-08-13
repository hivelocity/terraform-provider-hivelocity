# \NetworkApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetNetworkTaskClientResource**](NetworkApi.md#GetNetworkTaskClientResource) | **Get** /network/status/ | Returns the Last Status for a long running Network Task, such as modifying a VLAN
[**GetNetworkTaskDeviceResource**](NetworkApi.md#GetNetworkTaskDeviceResource) | **Get** /network/status/{deviceId} | Returns the Last Status for a long running Network Task, such as modifying a VLAN
[**GetNullRouteResource**](NetworkApi.md#GetNullRouteResource) | **Get** /network/null/{ip} | Null route an IP
[**GetRemoveNullRouteResource**](NetworkApi.md#GetRemoveNullRouteResource) | **Get** /network/unnull/{ip} | Remove null route from an IP
[**PostDetailedNullRouteResource**](NetworkApi.md#PostDetailedNullRouteResource) | **Post** /network/null | Null route an IP (with custom comments)


# **GetNetworkTaskClientResource**
> []NetworkTask GetNetworkTaskClientResource(ctx, optional)
Returns the Last Status for a long running Network Task, such as modifying a VLAN

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***NetworkApiGetNetworkTaskClientResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a NetworkApiGetNetworkTaskClientResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]NetworkTask**](NetworkTask.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetNetworkTaskDeviceResource**
> NetworkTask GetNetworkTaskDeviceResource(ctx, deviceId, optional)
Returns the Last Status for a long running Network Task, such as modifying a VLAN

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**|  | 
 **optional** | ***NetworkApiGetNetworkTaskDeviceResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a NetworkApiGetNetworkTaskDeviceResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**NetworkTask**](NetworkTask.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetNullRouteResource**
> GetNullRouteResource(ctx, ip)
Null route an IP

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **ip** | **string**|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetRemoveNullRouteResource**
> GetRemoveNullRouteResource(ctx, ip)
Remove null route from an IP

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **ip** | **string**|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostDetailedNullRouteResource**
> PostDetailedNullRouteResource(ctx, payload)
Null route an IP (with custom comments)

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**DetailedNullIp**](DetailedNullIp.md)|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

