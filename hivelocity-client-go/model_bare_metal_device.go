/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BareMetalDevice struct {
	PrimaryIp string `json:"primaryIp,omitempty"`
	Period string `json:"period,omitempty"`
	DeviceId int32 `json:"deviceId,omitempty"`
	OsName string `json:"osName,omitempty"`
	LocationName string `json:"locationName,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	ProductName string `json:"productName,omitempty"`
	ProductId int32 `json:"productId,omitempty"`
	ServiceId int32 `json:"serviceId,omitempty"`
	PowerStatus string `json:"powerStatus,omitempty"`
	VlanId int32 `json:"vlanId,omitempty"`
	PublicSshKeyId int32 `json:"publicSshKeyId,omitempty"`
	OrderId int32 `json:"orderId,omitempty"`
	Tags []string `json:"tags,omitempty"`
}
