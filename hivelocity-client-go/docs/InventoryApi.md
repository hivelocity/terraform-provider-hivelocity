# \InventoryApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetStockByProductResource**](InventoryApi.md#GetStockByProductResource) | **Get** /inventory/product/{productId} | Return a structured sps stock data, grouped by city or facility code for a single product
[**GetStockResource**](InventoryApi.md#GetStockResource) | **Get** /inventory/product | Return structured sps stock data, grouped by city or facility code for all products


# **GetStockByProductResource**
> Stock GetStockByProductResource(ctx, productId, optional)
Return a structured sps stock data, grouped by city or facility code for a single product

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **productId** | **int32**| Product database ID | 
 **optional** | ***InventoryApiGetStockByProductResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a InventoryApiGetStockByProductResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**Stock**](Stock.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **GetStockResource**
> Inventory GetStockResource(ctx, optional)
Return structured sps stock data, grouped by city or facility code for all products

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***InventoryApiGetStockResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a InventoryApiGetStockResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **location** | **optional.String**| Filter products by location | [default to MAIN]
 **groupBy** | **optional.String**| Get results grouped by &#39;city&#39;, &#39;facility&#39;, or &#39;flat&#39; | [default to facility]

### Return type

[**Inventory**](Inventory.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

