# \BandwidthApi

All URIs are relative to *https://localhost/api/v2*

Method | HTTP request | Description
------------- | ------------- | -------------
[**PostDeviceIdBandwidthImageResource**](BandwidthApi.md#PostDeviceIdBandwidthImageResource) | **Post** /bandwidth/device/{deviceId}/image | Returns RRDTool Graph based bandwidth in PNG format
[**PostDeviceIdBandwidthResource**](BandwidthApi.md#PostDeviceIdBandwidthResource) | **Post** /bandwidth/device/{deviceId} | Returns RRDTool Xport based bandwidth data in JSON format
[**PostServiceIdBandwidthImageResource**](BandwidthApi.md#PostServiceIdBandwidthImageResource) | **Post** /bandwidth/service/{serviceId}/image | Returns RRDTool Graph based bandwidth in PNG format
[**PostServiceIdBandwidthResource**](BandwidthApi.md#PostServiceIdBandwidthResource) | **Post** /bandwidth/service/{serviceId} | Returns RRDTool Xport based bandwidth data in JSON format


# **PostDeviceIdBandwidthImageResource**
> []BandwidthImage PostDeviceIdBandwidthImageResource(ctx, deviceId, period, interface_, optional)
Returns RRDTool Graph based bandwidth in PNG format

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View | 
  **period** | **string**| Preconfigured Time Periods for Graph Data | [default to day]
  **interface_** | **string**| Network Interface to use for Graph Data | [default to eth0]
 **optional** | ***BandwidthApiPostDeviceIdBandwidthImageResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a BandwidthApiPostDeviceIdBandwidthImageResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **start** | **optional.Int32**| Start Time of Custom Time Period. (Unix Epoch Time) | [default to 0]
 **end** | **optional.Int32**| End Time of Custom Time Period (Unix Epoch Time) | [default to 1620236450]
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]BandwidthImage**](BandwidthImage.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostDeviceIdBandwidthResource**
> []Bandwidth PostDeviceIdBandwidthResource(ctx, deviceId, period, interface_, step, optional)
Returns RRDTool Xport based bandwidth data in JSON format

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **deviceId** | **int32**| ID of Device to View | 
  **period** | **string**| Preconfigured Time Periods for Graph Data | [default to day]
  **interface_** | **string**| Network Interface to use for Graph Data | [default to eth0]
  **step** | **int32**| Interval of Graph in Seconds | [default to 300]
 **optional** | ***BandwidthApiPostDeviceIdBandwidthResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a BandwidthApiPostDeviceIdBandwidthResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **historical** | **optional.Bool**| Include Historical Interface Data for Device for Resellers | [default to false]
 **start** | **optional.Int32**| Start Time of Custom Time Period. (Unix Epoch Time) | [default to 0]
 **end** | **optional.Int32**| End Time of Custom Time Period (Unix Epoch Time) | [default to 1620236450]
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]Bandwidth**](Bandwidth.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostServiceIdBandwidthImageResource**
> []BandwidthImage PostServiceIdBandwidthImageResource(ctx, serviceId, period, interface_, optional)
Returns RRDTool Graph based bandwidth in PNG format

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **serviceId** | **int32**| ID of Service to View | 
  **period** | **string**| Preconfigured Time Periods for Graph Data | [default to day]
  **interface_** | **string**| Network Interface to use for Graph Data | [default to eth0]
 **optional** | ***BandwidthApiPostServiceIdBandwidthImageResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a BandwidthApiPostServiceIdBandwidthImageResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------



 **start** | **optional.Int32**| Start Time of Custom Time Period. (Unix Epoch Time) | [default to 0]
 **end** | **optional.Int32**| End Time of Custom Time Period (Unix Epoch Time) | [default to 1620236450]
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]BandwidthImage**](BandwidthImage.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **PostServiceIdBandwidthResource**
> []Bandwidth PostServiceIdBandwidthResource(ctx, serviceId, period, interface_, step, optional)
Returns RRDTool Xport based bandwidth data in JSON format

### Required Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
  **serviceId** | **int32**| ID of Service to View | 
  **period** | **string**| Preconfigured Time Periods for Graph Data | [default to day]
  **interface_** | **string**| Network Interface to use for Graph Data | [default to eth0]
  **step** | **int32**| Interval of Graph in Seconds | [default to 300]
 **optional** | ***BandwidthApiPostServiceIdBandwidthResourceOpts** | optional parameters | nil if no parameters

### Optional Parameters
Optional parameters are passed through a pointer to a BandwidthApiPostServiceIdBandwidthResourceOpts struct

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------




 **start** | **optional.Int32**| Start Time of Custom Time Period. (Unix Epoch Time) | [default to 0]
 **end** | **optional.Int32**| End Time of Custom Time Period (Unix Epoch Time) | [default to 1620236450]
 **xFields** | **optional.String**| An optional fields mask | 

### Return type

[**[]Bandwidth**](Bandwidth.md)

### Authorization

[apiKey](../README.md#apiKey)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

