/*
 * Hivelocity API for Partners
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type Deployment struct {
	StartedProvisioning bool `json:"startedProvisioning,omitempty"`
	DeploymentId int32 `json:"deploymentId,omitempty"`
	Empty bool `json:"empty,omitempty"`
	OrderNumber string `json:"orderNumber,omitempty"`
	Price float32 `json:"price,omitempty"`
	DeploymentName string `json:"deploymentName,omitempty"`
	DeploymentConfiguration []interface{} `json:"deploymentConfiguration,omitempty"`
}
