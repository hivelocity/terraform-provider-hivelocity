# \CancellationApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetCancellationDeviceResource**](CancellationApi.md#GetCancellationDeviceResource) | **Get** /cancellation/device/{deviceId} | Return the Cancellation found for a Device
[**GetCancellationIdResource**](CancellationApi.md#GetCancellationIdResource) | **Get** /cancellation/cancellation/{cancellationId} | Return any cancellation by ID
[**GetCancellationResource**](CancellationApi.md#GetCancellationResource) | **Get** /cancellation/cancellation | Return the services cancellations of a client
[**GetCancellationServiceResource**](CancellationApi.md#GetCancellationServiceResource) | **Get** /cancellation/service/{serviceId} | Return the Cancellation found for a Service
[**PostCancellationResource**](CancellationApi.md#PostCancellationResource) | **Post** /cancellation/cancellation | Creates Cancellation for a device/service


# **GetCancellationDeviceResource**
> Cancellation GetCancellationDeviceResource(ctx, deviceId, optional)
Return the Cancellation found for a Device

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**|  | 
 **optional** | ***CancellationApiGetCancellationDeviceResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CancellationApiGetCancellationDeviceResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Cancellation**](Cancellation.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCancellationIdResource**
> Cancellation GetCancellationIdResource(ctx, cancellationId, optional)
Return any cancellation by ID

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **cancellationId** | **int32**|  | 
 **optional** | ***CancellationApiGetCancellationIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CancellationApiGetCancellationIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Cancellation**](Cancellation.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCancellationResource**
> []Cancellation GetCancellationResource(ctx, optional)
Return the services cancellations of a client

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***CancellationApiGetCancellationResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CancellationApiGetCancellationResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]Cancellation**](Cancellation.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetCancellationServiceResource**
> Cancellation GetCancellationServiceResource(ctx, serviceId, optional)
Return the Cancellation found for a Service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **serviceId** | **int32**|  | 
 **optional** | ***CancellationApiGetCancellationServiceResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CancellationApiGetCancellationServiceResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Cancellation**](Cancellation.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostCancellationResource**
> Cancellation PostCancellationResource(ctx, payload, optional)
Creates Cancellation for a device/service

reason options are:     'My current service is being upgraded',     'Hardware/software/network or power issues',     'I think I have found a better deal',     'I have experienced tech support problems',     'I am consolidating my Hivelocity accounts',     'I am a reseller and my customer cancelled',     'I am moving to a different technology solution: Public Cloud',     'I am moving to a different technology solution: Managed Hosting Company',     'A product difference (Example: AWS)',     ''

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**CancellationCreate**](CancellationCreate.md)|  | 
 **optional** | ***CancellationApiPostCancellationResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a CancellationApiPostCancellationResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Cancellation**](Cancellation.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

