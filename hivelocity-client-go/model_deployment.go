/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Deployment struct {
	OrderNumber string `json:"orderNumber,omitempty"`
	DeploymentId int32 `json:"deploymentId,omitempty"`
	DeploymentConfiguration []interface{} `json:"deploymentConfiguration,omitempty"`
	Empty bool `json:"empty,omitempty"`
	DeploymentName string `json:"deploymentName,omitempty"`
	Price float32 `json:"price,omitempty"`
	StartedProvisioning bool `json:"startedProvisioning,omitempty"`
}
