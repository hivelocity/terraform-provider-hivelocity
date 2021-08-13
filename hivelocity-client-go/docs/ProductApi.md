# \ProductApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetProductListResource**](ProductApi.md#GetProductListResource) | **Get** /product/list | Return structured sps stock data in a list
[**GetProductOperatingSystemsResource**](ProductApi.md#GetProductOperatingSystemsResource) | **Get** /product/{productId}/operating-systems | Return List of operating systems found for a Product
[**GetProductOptionResource**](ProductApi.md#GetProductOptionResource) | **Get** /product/{productId}/options | Return List of Options found for a Product
[**GetProductsAndOptionsResource**](ProductApi.md#GetProductsAndOptionsResource) | **Get** /product/options | Return a mapping of Products and Options with pricing per-period
[**PostProductMatchResource**](ProductApi.md#PostProductMatchResource) | **Post** /product/match | Return a list of Products matching the provided lshw output of a server


# **GetProductListResource**
> []ProductInfo GetProductListResource(ctx, optional)
Return structured sps stock data in a list

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ProductApiGetProductListResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProductApiGetProductListResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]ProductInfo**](ProductInfo.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetProductOperatingSystemsResource**
> []Option GetProductOperatingSystemsResource(ctx, productId, optional)
Return List of operating systems found for a Product

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **productId** | **int32**| ID of the Product | 
 **optional** | ***ProductApiGetProductOperatingSystemsResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProductApiGetProductOperatingSystemsResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]Option**](Option.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetProductOptionResource**
> Options GetProductOptionResource(ctx, productId, optional)
Return List of Options found for a Product

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **productId** | **int32**| ID of the Product | 
 **optional** | ***ProductApiGetProductOptionResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProductApiGetProductOptionResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **groupBy** | **optional.String**| Get results grouped by &#39;upgrade&#39; or without grouping. | [default to upgrade]

### Return type

[**Options**](Options.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetProductsAndOptionsResource**
> []ProductOption GetProductsAndOptionsResource(ctx, optional)
Return a mapping of Products and Options with pricing per-period

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ProductApiGetProductsAndOptionsResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a ProductApiGetProductsAndOptionsResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]ProductOption**](ProductOption.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostProductMatchResource**
> PostProductMatchResource(ctx, payload)
Return a list of Products matching the provided lshw output of a server

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **payload** | [**ProductMatch**](ProductMatch.md)|  | 

### Return type

 (empty response body)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

