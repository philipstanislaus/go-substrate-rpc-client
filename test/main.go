package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/go-substrate-rpc-client"
	"github.com/centrifuge/go-substrate-rpc-client/scale"
	"github.com/centrifuge/go-substrate-rpc-client/system"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"golang.org/x/crypto/blake2b"
)

const (
	AnchorCommit    = "anchor.commit"
	AnchorPreCommit = "anchor.preCommit"
	SubKeySign      = "sign-blob"

	// Adjust below params according to your env + chain state + requirement
	//EndPoint = "127.0.0.1"
	EndPoint = "172.42.0.2"
	RPCPort = "9933"
	WSPort  = "9944"

	// SubKeyCmd subkey command to create signatures
	// SubKeyCmd = "/Users/philipstanislaus/.cargo/bin/subkey"
	SubKeyCmd = "subkey"

	NumAnchorsPerThread = 2
	Concurrency         = 1
)

type PreAnchorParams struct {
	AnchorID    [32]byte
	SigningRoot [32]byte
}

func (a PreAnchorParams) Encode(encoder scale.Encoder) error {
	encoder.Write(a.AnchorID[:])
	encoder.Write(a.SigningRoot[:])
	return nil
}

type AnchorParams struct {
	AnchorIDPreimage [32]byte
	DocRoot          [32]byte
	Proof            [32]byte
}

func NewRandomAnchorParam() AnchorParams {
	ap := AnchorParams{}
	copy(ap.AnchorIDPreimage[:], utils.RandomSlice(32))
	copy(ap.DocRoot[:], utils.RandomSlice(32))
	copy(ap.Proof[:], utils.RandomSlice(32))
	return ap
}

func NewRandomAnchorPreAnchorParams() (PreAnchorParams, AnchorParams) {
	pa := PreAnchorParams{}
	ap := AnchorParams{}

	// anchor id hashes
	copy(ap.AnchorIDPreimage[:], utils.RandomSlice(32))
	anchorID := blake2b.Sum256(ap.AnchorIDPreimage[:])
	copy(pa.AnchorID[:], anchorID[:])

	// root hashes
	copy(pa.SigningRoot[:], utils.RandomSlice(32))
	copy(ap.Proof[:], utils.RandomSlice(32))

	var docRootPreimage []byte
	if bytes.Compare(pa.SigningRoot[:], ap.Proof[:]) < 0 {
		docRootPreimage = append(pa.SigningRoot[:], ap.Proof[:]...)
	} else {
		docRootPreimage = append(ap.Proof[:], pa.SigningRoot[:]...)
	}

	ap.DocRoot = blake2b.Sum256(docRootPreimage)
	return pa, ap
}

func NewAnchorParamsFromHex(apre, docr, proof string) AnchorParams {
	a, _ := hexutil.Decode(apre)
	d, _ := hexutil.Decode(docr)
	p, _ := hexutil.Decode(proof)
	ap := AnchorParams{}
	copy(ap.AnchorIDPreimage[:], a)
	copy(ap.DocRoot[:], d)
	copy(ap.Proof[:], p)
	return ap
}

func (a *AnchorParams) AnchorIDHex() string {
	b := blake2b.Sum256(a.AnchorIDPreimage[:])
	return hexutil.Encode(b[:])
}

func (a *AnchorParams) Decode(decoder scale.Decoder) error {
	decoder.Read(a.AnchorIDPreimage[:])
	decoder.Read(a.DocRoot[:])
	decoder.Read(a.Proof[:])
	return nil
}

func (a AnchorParams) Encode(encoder scale.Encoder) error {
	encoder.Write(a.AnchorIDPreimage[:])
	encoder.Write(a.DocRoot[:])
	encoder.Write(a.Proof[:])
	return nil
}

type AnchorData struct {
	ID            [32]byte
	DocRoot       [32]byte
	AnchoredBlock uint64
}

func (a *AnchorData) Decode(decoder scale.Decoder) error {
	decoder.Read(a.ID[:])
	decoder.Read(a.DocRoot[:])
	decoder.Decode(&a.AnchoredBlock)
	return nil
}

func Anchors(client substrate.Client, module string, fn string, anchorIDPreImage []byte) (*AnchorData, error) {
	h := blake2b.Sum256(anchorIDPreImage)
	m, err := client.MetaData(true)
	if err != nil {
		return nil, err
	}

	key, err := substrate.NewStorageKey(*m, module, fn, h[:])
	if err != nil {
		return nil, err
	}

	s := substrate.NewStateRPC(client)
	res, err := s.Storage(key, nil)
	if err != nil {
		return nil, err
	}

	tempDec := res.Decoder()
	a := AnchorData{}
	err = tempDec.Decode(&a)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func main() {
	waitForServer()

	time.Sleep(5 * time.Second)

	// Connect the client.
	client, err := substrate.Connect(fmt.Sprintf("ws://%v:%v", EndPoint, WSPort))
	if err != nil {
		panic(err)
	}
	alice, err := hexutil.Decode(substrate.AlicePubKey)
	if err != nil {
		panic(err)
	}
	nonce, err := system.AccountNonce(client, alice)
	if err != nil && err.Error() == "empty result" {
		nonce = 0
	} else if err != nil {
		panic(err)
	}

	gs, err := system.BlockHash(client, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println("Nonce", nonce, "gs", gs.String(), "Alice", string(alice))

	authRPC := substrate.NewAuthorRPC(client, gs, SubKeyCmd, SubKeySign)
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(Concurrency)
	var counter uint64
	for i := 0; i < Concurrency; i++ {
		go func() {
			for i := 0; i < NumAnchorsPerThread; i++ {
				// a := NewAnchorParamsFromHex("0x0000000000000000000000000000000000000000000000000000000000000901", "0x0000000000000000000000000000000000000000000000000000000000000000", "0x0000000000000000000000000000000000000000000000000000000000000000")
				pa, ap := NewRandomAnchorPreAnchorParams()
				aID := ap.AnchorIDHex()

				res, err := authRPC.SubmitExtrinsic(nonce, AnchorPreCommit, pa)
				if err != nil {
					fmt.Printf("FAIL!!! pre commit for anchor ID %s failed with %s\n", aID, err.Error())
					break
				} else {
					// verify pre anchor
					for i := 0; i < 10; i++ {
						time.Sleep(10 * time.Second)
						anc, err := Anchors(client, "Anchor", "PreAnchors", ap.AnchorIDPreimage[:])
						fmt.Println(err)
						if anc != nil {
							fmt.Printf("SUCCESS!!! pre anchor %v\n", anc)
							break
						}
					}
					fmt.Printf("SUCCESS!!! pre anchor ID %s , tx hash %s\n", aID, res)
					atomic.AddUint64(&nonce, 1)
				}

				// fmt.Println("submitting new anchor with anchor ID", a.AnchorIDHex())
				res, err = authRPC.SubmitExtrinsic(nonce, AnchorCommit, ap)
				if err != nil {
					fmt.Printf("FAIL!!! commit for anchor ID %s failed with %s\n", aID, err.Error())
					break
				} else {
					// verify anchor
					for i := 0; i < 10; i++ {
						time.Sleep(10 * time.Second)
						anc, err := Anchors(client, "Anchor", "Anchors", ap.AnchorIDPreimage[:])
						fmt.Println(err)
						if anc != nil {
							fmt.Printf("SUCCESS!!! anchor %v\n", anc)
							break
						}
					}
					fmt.Printf("SUCCESS!!! anchor ID %s , tx hash %s\n", aID, res)
					atomic.AddUint64(&counter, 1)
					atomic.AddUint64(&nonce, 1)
				}
			}
			wg.Done()
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	tps := float64(counter) / elapsed.Seconds()
	fmt.Printf("Successful execution of %d transactions took %s, amounting to %f TPS\n", counter, elapsed, tps)
}

func waitForServer() {
	reqBody, err := json.Marshal(map[string]string{
		"id": "1", "jsonrpc": "2.0", "method": "system_health",
	})
	if err != nil {
		panic(err)
	}

	for i := 0; true; i++ {
		fmt.Println("iteration", i)
		resp, err := http.Post(fmt.Sprintf("http://%v:%v", EndPoint, RPCPort), "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			fmt.Println("Error receiving HTTP response, will retry. Error:", err.Error())
			time.Sleep(100 * time.Millisecond)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		resp.Body.Close()

		type RpcResp struct {
			JsonRPC string `json:"jsonrpc"`
			Result  struct {
				IsSyncing       bool    `json:"isSyncing"`
				Peers           float64 `json:"peers"`
				ShouldHavePeers bool    `json:"shouldHavePeers"`
			} `json:"result"`
			Id string `json:"id"`
		}

		var rpcResp RpcResp

		err = json.Unmarshal(body, &rpcResp)
		if err != nil {
			fmt.Println("Error unmarhalling JSON, will retry.")
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if rpcResp.Id == "" {
			fmt.Println("Got unexpected response, will retry. Response:", rpcResp)
			time.Sleep(100 * time.Millisecond)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			return
		}
	}
}
