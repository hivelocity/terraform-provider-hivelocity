/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type ContactDump struct {
	FullName string `json:"fullName"`
	Description string `json:"description,omitempty"`
	Phone string `json:"phone,omitempty"`
	IsClient bool `json:"isClient,omitempty"`
	ContactId int32 `json:"contactId,omitempty"`
	Email string `json:"email"`
	ClientId int32 `json:"clientId,omitempty"`
	Active int32 `json:"active"`
}
