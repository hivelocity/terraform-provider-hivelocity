/*
 * Hivelocity API
 *
 * Interact with Hivelocity
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

type TicketPost struct {
	Subject string `json:"subject,omitempty"`
	Id float32 `json:"id,omitempty"`
	UbersmithAttachedFiles *TicketAttach `json:"ubersmith_attached_files,omitempty"`
	ClientId float32 `json:"clientId,omitempty"`
	From *interface{} `json:"from,omitempty"`
	Body string `json:"body,omitempty"`
	Date float32 `json:"date,omitempty"`
	Attachments float32 `json:"attachments,omitempty"`
	TicketId float32 `json:"ticketId,omitempty"`
	AdminId float32 `json:"adminId,omitempty"`
	ContactId float32 `json:"contactId,omitempty"`
	Hidden float32 `json:"hidden,omitempty"`
	FromAdmin bool `json:"fromAdmin,omitempty"`
}
