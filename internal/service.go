package internal

type ServiceOptions struct{}

type Service struct{}

func (s Service) ListTransactions(opts *ListTransactionOptions) ([]Transaction, *Problem) {
	return []Transaction{
		{
			ID:         1,
			ExternalID: "1",
		},
	}, nil
}

func (s Service) CreateTransaction(body *CreateTransactionInput) (*Transaction, *Problem) {
	return &Transaction{
		ID:         1,
		ExternalID: body.ExternalID,
	}, nil
}

func (s Service) GetTransaction(id uint64) (*Transaction, *Problem) {
	return &Transaction{
		ID: id,
	}, nil
}
