package core

type Config struct {
	CacheExpiryTime int
	Port            int
	EtcdHostname    string
	EtcdPort        int
	EtcdUsername    string
	EtcdPassword    string
}

type IdentifierType string

const (
	IdentifierDatabase IdentifierType = "mdb"
)

type Identifier struct {
	IdentifierType IdentifierType
	//UUID           UUID
}

type Database struct {
	Identifier  Identifier
	Name        string
	Description string
	Comment     string
	Plugin      string
	Config      Config
}

// TODO
type TransactionOperation interface{}

type Consistency int

const (
	ConsistencySerializable Consistency = iota
	ConsistencyRelaxed
	ConsistencyNone
)

type Transaction struct {
	Operations  []TransactionOperation
	Consistency Consistency
	UseIDs      bool
}
