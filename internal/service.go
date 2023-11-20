package internal

type ServiceOptions struct{}

type TransactionService struct{}

func (s TransactionService) ListTransactions(opts *ListTransactionOptions) ([]Transaction, *Problem) {
	return []Transaction{
		{
			ID:         1,
			ExternalID: "1",
		},
	}, nil
}

func (s TransactionService) CreateTransaction(body *CreateTransactionInput) (*Transaction, *Problem) {
	return &Transaction{
		ID:         1,
		ExternalID: body.ExternalID,
	}, nil
}

func (s TransactionService) GetTransaction(id uint64) (*Transaction, *Problem) {
	return &Transaction{
		ID: id,
	}, nil
}
