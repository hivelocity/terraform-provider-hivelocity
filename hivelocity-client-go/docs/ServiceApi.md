# \ServiceApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetServiceIdResource**](ServiceApi.md#GetServiceIdResource) | **Get** /service/{serviceId} | Return a dictionary with data from a specific account service
[**GetServiceResource**](ServiceApi.md#GetServiceResource) | **Get** /service/ | Return a list of all account services


# **GetServiceIdResource**
> ServiceDump GetServiceIdResource(ctx, serviceId, optional)
Return a dictionary with data from a specific account service

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **serviceId** | **int32**|  | 
 **optional** | ***ServiceApiGetServiceIdResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ServiceApiGetServiceIdResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**ServiceDump**](Service-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetServiceResource**
> []ServiceDump GetServiceResource(ctx, status, optional)
Return a list of all account services

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **status** | **string**| The current status of the service | [default to all]
 **optional** | ***ServiceApiGetServiceResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ServiceApiGetServiceResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]ServiceDump**](Service-dump.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

