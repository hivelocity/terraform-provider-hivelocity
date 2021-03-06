/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Bandwidth struct {
	Metadata *interface{} `json:"metadata,omitempty"`
	BandwidthData [][]float32 `json:"bandwidthData,omitempty"`
	Interface_ string `json:"interface,omitempty"`
	SwitchId string `json:"switchId,omitempty"`
}
