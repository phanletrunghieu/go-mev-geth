package bundle

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"

	myhttp "github.com/phanletrunghieu/go-mev-geth/common/http"
)

var (
	Relay           = "https://relay.flashbots.net"
	TestnetRelay    = "https://relay-goerli.flashbots.net"
	method_simulate = "eth_callBundle"
	method_send     = "eth_sendBundle"
	id              = int64(0)
)

type (
	JsonRpc struct {
		Method  string        `json:"method"`
		Params  []interface{} `json:"params"`
		ID      int64         `json:"id"`
		Jsonrpc string        `json:"jsonrpc"`
	}

	Bundle struct {
		Relay              string   `json:"-"`
		Signer             string   `json:"-"`
		SignedTransactions []string `json:"txs"`
		BlockNumber        string   `json:"blockNumber"`
		MinTimestamp       *int     `json:"minTimestamp,omitempty"`
		MaxTimestamp       *int     `json:"maxTimestamp,omitempty"`
		StateBlockNumber   *string  `json:"stateBlockNumber,omitempty"`
	}

	Response interface{}

	ErrResponse struct {
		Error struct {
			Message string `json:"message"`
			Code    int    `json:"code"`
		} `json:"error"`
	}
)

func NewBundle(
	relay string,
	signer string,
	signedTransactions []string,
	blockNumber uint64,
) *Bundle {
	return &Bundle{
		Relay:              relay,
		Signer:             signer,
		SignedTransactions: signedTransactions,
		BlockNumber:        "0x" + fmt.Sprintf("%x", blockNumber),
	}
}

func (b *Bundle) Send(client *http.Client) (res Response, err error) {
	payload := b.prepareRequest(method_send)
	signature, err := b.sign(payload)
	if err != nil {
		return nil, err
	}
	signerAddress, err := b.signerAddress()
	if err != nil {
		return nil, err
	}
	err = myhttp.Post(b.Relay, payload, &res, client, map[string]string{"X-Flashbots-Signature": signerAddress + ":" + signature})
	return res, err
}

func (b *Bundle) sign(jsonRpc JsonRpc) (signature string, err error) {
	marshal, err := json.Marshal(jsonRpc)
	if err != nil {
		return "", err
	}

	ecdsaPrivateKey, err := crypto.HexToECDSA(b.Signer)
	if err != nil {
		return "", err
	}
	signatureBytes, err := crypto.Sign(accounts.TextHash([]byte(hexutil.Encode(crypto.Keccak256(marshal)))), ecdsaPrivateKey)
	if err != nil {
		return "", err
	}
	return hexutil.Encode(signatureBytes), nil
}

func (b *Bundle) Simulate(client *http.Client) (res Response, err error) {
	payload := b.prepareRequest(method_simulate)
	signature, err := b.sign(payload)
	if err != nil {
		return nil, err
	}
	singerAddress, err := b.signerAddress()
	if err != nil {
		return nil, err
	}
	err = myhttp.Post(b.Relay, payload, &res, client, map[string]string{"X-Flashbots-Signature": singerAddress + ":" + signature})
	return res, err
}

func (b *Bundle) prepareRequest(method string) JsonRpc {
	id++

	stateBlockNumber := "latest"
	bundleParam := *b

	if method == method_simulate {
		bundleParam.StateBlockNumber = &stateBlockNumber
	}

	return JsonRpc{
		Jsonrpc: "2.0",
		Method:  method,
		Params: []interface{}{
			bundleParam,
		},
		ID: id,
	}
}

func (b *Bundle) signerAddress() (address string, err error) {
	ecdsaPrivateKey, err := crypto.HexToECDSA(b.Signer)
	if err != nil {
		return "", err
	}
	publicKey := ecdsaPrivateKey.PublicKey
	addr := crypto.PubkeyToAddress(publicKey)
	return addr.String(), nil
}
