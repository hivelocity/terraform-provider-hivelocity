/*
 * Hivelocity API for Partners
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type OrderDump struct {
	Status string `json:"status,omitempty"`
	OrderId int32 `json:"orderId,omitempty"`
	Owner string `json:"owner,omitempty"`
	Info interface{} `json:"info,omitempty"`
	Total float32 `json:"total,omitempty"`
}
