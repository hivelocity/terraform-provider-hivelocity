/*
 * Hivelocity API for Partners
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type BareMetalDevice struct {
	Tags []string `json:"tags,omitempty"`
	LocationName string `json:"locationName,omitempty"`
	VlanId int32 `json:"vlanId,omitempty"`
	PublicSshKeyId int32 `json:"publicSshKeyId,omitempty"`
	ProductId int32 `json:"productId,omitempty"`
	Period string `json:"period,omitempty"`
	Script string `json:"script,omitempty"`
	OsName string `json:"osName,omitempty"`
	ServiceId int32 `json:"serviceId,omitempty"`
	ProductName string `json:"productName,omitempty"`
	PowerStatus string `json:"powerStatus,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	DeviceId int32 `json:"deviceId,omitempty"`
	OrderId int32 `json:"orderId,omitempty"`
	PrimaryIp string `json:"primaryIp,omitempty"`
}
