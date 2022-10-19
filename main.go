package main

import (
	"encoding/json"
	"log"
	"net/url"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
)

type SubscribeMsg struct {
	Id     int32    `json:"id"`
	Method string   `json:"method"`
	Params []string `json:"params"`
}

type SubscribeResponseMsg struct {
	JsonRpc string `json:"jsonrpc"`
	Id      int32  `json:"id"`
	Result  string `json:"result"`
}

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u, err := url.Parse("ws://127.0.0.1:8546")
	if err != nil {
		log.Fatal(err)
	}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	newPedingTranSubMsg := SubscribeMsg{
		Id:     1,
		Method: "eth_subscribe",
		Params: []string{"newPendingTransactions"},
	}
	sendMsg(c, newPedingTranSubMsg)

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")
			return
		}
	}
}

func sendMsg(ws *websocket.Conn, msg interface{}) {
	ws.WriteJSON(msg)
	b, _ := json.Marshal(msg)
	log.Printf("Send: %#v\n", string(b))
}

/**
Another implementation

func StreamTx(rpcClient *rpc.Client, fullMode bool) {

	// Go channel to pipe data from client subscription
	txChannel := make(chan common.Hash)

	// Subscribe to receive one time events for new txs
	rpcClient.EthSubscribe(
		context.Background(), txChannel, "newPendingTransactions", // no additional args
	)
	client := GetCurrentClient()
	fmt.Println("Subscribed to mempool txs")

	// Configure chain ID and signer to ensure you're configured to mainnet
	chainID, _ := client.NetworkID(context.Background())
	signer := types.NewEIP155Signer(chainID)

	for {
		select {
		case transactionHash := <-txChannel:
			tx, is_pending, _ := client.TransactionByHash(context.Background(), transactionHash)
			if is_pending {
				_, _ = signer.Sender(tx)
				//fmt.Println(fullMode)
				//handleTransaction(tx, client, false, fullMode)
			}
		}
	}
}
**/
