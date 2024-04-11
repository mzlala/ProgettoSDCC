package common

// Args rappresenta gli argomenti delle chiamate RPC
type Args struct {
	Key   string
	Value string
}

// Response è una struttura creata per memorizzare la risposta delle chiamate RPC
type Response struct {
	Reply string
}

// Datastore mantenuto da ciascun server
type Datastore struct {
	datastore map[string]string // Mappa -> struttura dati che associa chiavi a valori
}

// KeyValueStoreService è un'interfaccia rappresentante che chiamate RPC esposte al client
type KeyValueStoreService interface {
	Get(args Args, reply *Response) error
	Put(args Args, reply *Response) error
	Delete(args Args, reply *Response) error
}
