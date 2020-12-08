/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Device struct {
	DeviceId int32 `json:"deviceId,omitempty"`
	Name string `json:"name"`
	Status string `json:"status,omitempty"`
	DeviceType string `json:"deviceType,omitempty"`
	PowerStatus string `json:"powerStatus,omitempty"`
	HasCancellation bool `json:"hasCancellation,omitempty"`
	IsManaged bool `json:"isManaged,omitempty"`
	MonitorsUp int32 `json:"monitorsUp,omitempty"`
	MonitorsTotal int32 `json:"monitorsTotal,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	IpmiEnabled bool `json:"ipmiEnabled,omitempty"`
	DisplayedTags []interface{} `json:"displayedTags,omitempty"`
	Tags []string `json:"tags,omitempty"`
	Location *interface{} `json:"location,omitempty"`
	PrimaryIp string `json:"primaryIp,omitempty"`
	IpmiAddress string `json:"ipmiAddress,omitempty"`
	ServiceMonitors []string `json:"serviceMonitors,omitempty"`
	ServicePlan int32 `json:"servicePlan,omitempty"`
	SelfProvisioning bool `json:"selfProvisioning,omitempty"`
	Metadata *interface{} `json:"metadata,omitempty"`
}
