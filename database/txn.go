package database

// Account is an individual
type Account string

// Txn stores info about each txn
type Txn struct {
	From  Account `json:"from"`
	To    Account `json:"to"`
	Value uint    `json:"value"`
	Data  string  `json:"data"`
}

// NewAccount creates a new account with the given value
func NewAccount(value string) Account {
	return Account(value)
}

// NewTxn creates a new txn based on the given details
func NewTxn(from Account, to Account, value uint, data string) Txn {
	return Txn{from, to, value, data}
}

// IsReward() checks if the txn is a reward
func (t Txn) IsReward() bool {
	return t.Data == "reward"
}
