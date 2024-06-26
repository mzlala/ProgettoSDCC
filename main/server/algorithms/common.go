package algorithms

import (
	"fmt"
	"main/common"
	commonMsg "main/server/message"
	"net/rpc"
)

// SendToAllServer Invia un messaggio a tutte le repliche del server in goroutine
//   - rpcName: nome della procedura remota da chiamare
//   - message: messaggio da inviare (argomento della procedura remota)
//   - response: risposta della procedura remota
func SendToAllServer(rpcName string, message interface{}, response *common.Response) error {
	// Canale per ricevere i risultati delle chiamate RPC
	resultChan := make(chan error, common.Replicas)

	// Itera su tutte le repliche e avvia le chiamate RPC
	for i := 0; i < common.Replicas; i++ {
		go callRPC(rpcName, message, response, resultChan, i)
	}

	// Raccoglie i risultati dalle chiamate RPC
	for i := 0; i < common.Replicas; i++ {
		if err := <-resultChan; err != nil {
			return err
		}
	}
	return nil
}

// callRPC Esegue una chiamata RPC a una replica del server in modo sincrono
func callRPC(rpcName string, message interface{}, response *common.Response, resultChan chan<- error, replicaIndex int) {
	serverName := common.GetServerName(common.ReplicaPorts[replicaIndex], replicaIndex)

	conn, err := rpc.Dial("tcp", serverName)
	if err != nil {
		resultChan <- fmt.Errorf("errore durante la connessione con %s: %v", serverName, err)
		return
	}

	common.RandomDelay()

	switch msg := message.(type) {
	case commonMsg.MessageC:
		err = conn.Call(rpcName, msg, response)
	case commonMsg.MessageS:
		err = conn.Call(rpcName, msg, response)
	default:
		resultChan <- fmt.Errorf("tipo di messaggio non supportato: %T", msg)
		return
	}

	if err != nil {
		resultChan <- fmt.Errorf("errore durante la chiamata RPC %s a %s: %v", rpcName, serverName, err)
		return
	}

	err = conn.Close()
	if err != nil {
		resultChan <- fmt.Errorf("errore durante la connessione in KeyValueStoreCausale.callRPC: %s", err)
		return
	}

	resultChan <- nil
}
