package main

import (
	"fmt"
	"main/common"
	"net/rpc"
	"sync"
	"time"
)

// KeyValueStoreCausale rappresenta il servizio di memorizzazione chiave-valore
type KeyValueStoreCausale struct {
	datastore map[string]string // Mappa -> struttura dati che associa chiavi a valori
	id        int               // Id che identifica il server stesso

	vectorClock [common.Replicas]int // Orologio vettoriale
	mutexClock  sync.Mutex

	queue      []MessageC
	mutexQueue sync.Mutex // Mutex per proteggere l'accesso concorrente alla coda
}

type MessageC struct {
	Id            string // Id del messaggio stesso
	IdSender      int    // IdSender rappresenta l'indice del server che invia il messaggio
	TypeOfMessage string
	Args          common.Args
	VectorClock   [common.Replicas]int // Orologio vettoriale
	NumberAck     int
}

// Get restituisce il valore associato alla chiave specificata -> Lettura -> Evento interno
func (kvc *KeyValueStoreCausale) Get(args common.Args, response *common.Response) error {
	kvc.mutexClock.Lock()
	kvc.vectorClock[kvc.id]++
	message := MessageC{common.GenerateUniqueID(), kvc.id, "Get", args, kvc.vectorClock, 3}
	kvc.mutexClock.Unlock()

	kvc.addToQueue(message)

	// Controllo in while se il messaggio può essere inviato a livello applicativo
	for {
		canSend := kvc.controlSendToApplication(message)
		if canSend {
			// Invio a livello applicativo
			err := kvc.RealFunction(message, response)
			if err != nil {
				return err
			}
			break
		}

		// Altrimenti, attendi un breve periodo prima di riprovare
		time.Sleep(time.Millisecond * 100)
	}
	return nil
}

// Put inserisce una nuova coppia chiave-valore, se la chiave è già presente, sovrascrive il valore associato
func (kvc *KeyValueStoreCausale) Put(args common.Args, response *common.Response) error {
	fmt.Println("KeyValueStoreCausale: Comando PUT eseguito")

	kvc.mutexClock.Lock()
	kvc.vectorClock[kvc.id]++

	// CREO IL MESSAGGIO E DEVO FAR SI CHE TUTTI LO SCRIVONO NEL DATASTORE
	message := MessageC{common.GenerateUniqueID(), kvc.id, "Put", args, kvc.vectorClock, 0}
	kvc.mutexClock.Unlock()

	err := kvc.sendToAllServer("KeyValueStoreCausale.CausallyOrderedMulticast", message)
	if err != nil {
		return err
	}

	response.Reply = "true"
	return nil
}

// Delete elimina una coppia chiave-valore, è un operazione di scrittura
func (kvc *KeyValueStoreCausale) Delete(args common.Args, response *common.Response) error {
	fmt.Println("KeyValueStoreCausale: Comando PUT eseguito")

	kvc.mutexClock.Lock()
	kvc.vectorClock[kvc.id]++
	// CREO IL MESSAGGIO E DEVO FAR SI CHE TUTTI LO SCRIVONO NEL DATASTORE
	message := MessageC{common.GenerateUniqueID(), kvc.id, "Delete", args, kvc.vectorClock, 0}
	kvc.mutexClock.Unlock()

	err := kvc.sendToAllServer("KeyValueStoreCausale.CausallyOrderedMulticast", message)
	if err != nil {
		return err
	}

	response.Reply = "true"
	return nil
}

// sendToAllServer invia a tutti i server la richiesta rpcName
func (kvc *KeyValueStoreCausale) sendToAllServer(rpcName string, message MessageC) error {
	// Canale per ricevere i risultati delle chiamate RPC
	resultChan := make(chan error, common.Replicas)

	// Itera su tutte le repliche e avvia le chiamate RPC
	for i := 0; i < common.Replicas; i++ {
		go kvc.callRPC(rpcName, message, resultChan, i)
	}

	// Raccoglie i risultati dalle chiamate RPC
	for i := 0; i < common.Replicas; i++ {
		if err := <-resultChan; err != nil {
			return err
		}
	}
	return nil
}

// callRPC è una funzione ausiliaria per effettuare la chiamata RPC a una replica specifica
func (kvc *KeyValueStoreCausale) callRPC(rpcName string, message MessageC, resultChan chan<- error, replicaIndex int) {

	serverName := common.GetServerName(common.ReplicaPorts[replicaIndex], replicaIndex)

	//fmt.Println("sendToAllServer: Contatto", serverName)
	conn, err := rpc.Dial("tcp", serverName)
	if err != nil {
		// Gestione dell'errore durante la connessione al server
		resultChan <- fmt.Errorf("errore durante la connessione con %s: %v", serverName, err)
		return
	}

	// Chiama il metodo "rpcName" sul server
	var response bool
	err = conn.Call(rpcName, message, &response)
	if err != nil {
		// Gestione dell'errore durante la chiamata RPC
		resultChan <- fmt.Errorf("errore durante la chiamata RPC %s a %s: %v", rpcName, serverName, err)
		return
	}

	// Aggiungi il risultato al canale dei risultati
	resultChan <- nil
}

// RealFunction esegue l'operazione di put e di delete realmente
func (kvc *KeyValueStoreCausale) RealFunction(message MessageC, _ *common.Response) error {

	if message.TypeOfMessage == "Put" { // Scrittura
		kvc.mutexClock.Lock()
		fmt.Println("Key inserita ", message.Args.Key)
		kvc.datastore[message.Args.Key] = message.Args.Value
		fmt.Println("DATASTORE:")
		fmt.Println(kvc.datastore)
		kvc.mutexClock.Unlock()

	} else if message.TypeOfMessage == "Delete" { // Scrittura
		kvc.mutexClock.Lock()
		fmt.Println("Key eliminata ", message.Args.Key)
		delete(kvc.datastore, message.Args.Key)
		fmt.Println("DATASTORE:")
		fmt.Println(kvc.datastore)
		kvc.mutexClock.Unlock()

	} else {
		return fmt.Errorf("command not found")
	}

	return nil
}
