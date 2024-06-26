package sequential

import (
	"fmt"
	"github.com/fatih/color"
	"main/common"
	"main/server/message"
	"time"
)

// ----- Stampa messaggi ----- //

const layoutTime = "15:04:05.0000"

// Usato per messaggi di ricezione
func printDebugBlue(blueString string, message commonMsg.MessageS, kvs *KeyValueStoreSequential) {

	// Ottieni l'orario corrente
	now := time.Now()
	// Formatta l'orario corrente come stringa nel formato desiderato
	formattedTime := now.Format(layoutTime)

	if common.GetDebug() {
		switch message.GetTypeOfMessage() {
		case common.Del:
			fmt.Println(color.BlueString(blueString), message.GetTypeOfMessage(), message.GetKey(),
				"clockClient", message.GetSendingFIFO(), "clockMsg:", message.GetClock(), "clockServer:", kvs.GetClock(), formattedTime)
		default:
			fmt.Println(color.BlueString(blueString), message.GetTypeOfMessage(), message.GetKey()+":"+message.GetValue(),
				"clockClient", message.GetSendingFIFO(), "clockMsg:", message.GetClock(), "clockServer:", kvs.GetClock(), formattedTime)
		}
	} else {
		switch message.GetTypeOfMessage() {
		case common.Del:
			fmt.Println(color.BlueString(blueString), message.GetTypeOfMessage(), message.GetKey())
		default:
			fmt.Println(color.BlueString(blueString), message.GetTypeOfMessage(), message.GetKey()+":"+message.GetValue())
		}
	}
}

func printGreen(greenString string, message commonMsg.MessageS, kvs *KeyValueStoreSequential) {
	// Ottieni l'orario corrente
	now := time.Now()
	// Formatta l'orario corrente come stringa nel formato desiderato
	formattedTime := now.Format(layoutTime)

	if common.GetDebug() {
		switch message.GetTypeOfMessage() {
		case common.Del:
			fmt.Println(color.GreenString(greenString), message.GetTypeOfMessage(), message.GetKey(),
				"clockClient", message.GetSendingFIFO(), "clockMsg:", message.GetClock(), "clockServer:", kvs.GetClock(), formattedTime)

		default:
			fmt.Println(color.GreenString(greenString), message.GetTypeOfMessage(), message.GetKey()+":"+message.GetValue(),
				"clockClient", message.GetSendingFIFO(), "clockMsg:", message.GetClock(), "clockServer:", kvs.GetClock(), formattedTime)
		}

	} else {
		switch message.GetTypeOfMessage() {
		case common.Del:
			fmt.Println(color.GreenString(greenString), message.GetTypeOfMessage(), message.GetKey())
		default:
			fmt.Println(color.GreenString(greenString), message.GetTypeOfMessage(), message.GetKey()+":"+message.GetValue())
		}
	}
}

func printRed(redString string, message commonMsg.MessageS, kvs *KeyValueStoreSequential) {

	// Ottieni l'orario corrente
	now := time.Now()
	// Formatta l'orario corrente come stringa nel formato desiderato
	formattedTime := now.Format(layoutTime)

	if common.GetDebug() {
		fmt.Println(color.RedString(redString), message.GetTypeOfMessage(), message.GetKey(), "datastore:", kvs.GetDatastore(),
			"clockClient", message.GetSendingFIFO(), "clockMsg:", message.GetClock(), "clockServer:", kvs.GetClock(), formattedTime)
	} else {
		fmt.Println(color.RedString(redString), message.GetTypeOfMessage(), message.GetKey())
	}
}
