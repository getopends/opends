package internal

type Transaction struct {
	ID         uint64 `json:"id"`
	ExternalID string `json:"external_id"`
}

type CreateTransactionInput struct {
	ExternalID string `json:"external_id"`
}

type ValidateReceivingMethodInput struct {
	SendingCustomer   *SendingCustomer   `json:"sending_customer,omitempty"`
	SendingBusiness   *SendingBusiness   `json:"sending_business,omitempty"`
	ReceivingCustomer *ReceivingCustomer `json:"receiving_customer,omitempty"`
	ReceivingBusiness *ReceivingBusiness `json:"receiving_business,omitempty"`
	ReceivingMethod   *ReceivingMethod   `json:"receiving_method,omitempty"`
}

type RetrieveReceivingMethodInput struct {
	ReceivingMethod *ReceivingMethod `json:"receiving_method,omitempty"`
}

type SendingBusiness struct {
	Name string `json:"name"`
}

type SendingCustomer struct {
	FirstName string `json:"first_name"`
}

type ReceivingBusiness struct {
	Name string `json:"name"`
}

type ReceivingCustomer struct {
	Firstname string `json:"firstname"`
}

type ReceivingMethod struct {
	CardNumber  string `json:"card_number"`
	PhoneNumber string `json:"phone_number"`
}

type Problem struct {
	Type     string `json:"type,omitempty"`
	Title    string `json:"title,omitempty"`
	Status   int64  `json:"status,omitempty"`
	Detail   string `json:"detail,omitempty"`
	Instance string `json:"instance,omitempty"`
}

type ListTransactionOptions struct {
	ExternalID string `json:"external_id"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
}

type ListProductsOptions struct {
	OperatorID string `json:"operator_id"`
	ServiceID  string `json:"service_id"`
}
