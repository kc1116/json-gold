package lds

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/kc1116/json-gold/ld"
)

//Signature . . .
//Struct that defines a signature object based on the https://web-payments.org/specs/source/ld-signatures/#linked-data-signature-overview spec
type Signature struct {
	Type           string `json:"type"`
	Creator        string `json:"creator"`
	Created        string `json:"created"`
	Domain         string `json:"json"`
	Nonce          string `json:"nonce"`
	SignatureValue string `json:"signatureValue"`
}

//SignatureOptions . . .
//Struct that defines a signature options object based on the https://web-payments.org/specs/source/ld-signatures/#signature-algorithm spec
//Used to pass signature options into signature function
type SignatureOptions struct {
	Context string `json:"context"`
	Creator string `json:"creator"`
	Created string `json:"created"`
	Domain  string `json:"json"`
	Nonce   string `json:"nonce"`
}

//S2015 . . .
//Function S2015 (Signature 2015) that signs a linked data document based on the https://web-payments.org/specs/source/ld-signatures/#signature-algorithm spec
func S2015(document map[string]interface{}, options SignatureOptions, privateKey string) {
	var proc = ld.NewJsonLdProcessor()
	var nOpts = ld.NewJsonLdOptions("")
	nOpts.Format = "application/nquads"

	log.Println("Here is your element:", document)
	log.Println("\nHere are your options:", options)
	log.Println("Here is your privateKey:", privateKey)

	//make copy of original document
	output := copyMap(document)
	//remove any signature nodes from output
	removeSignature(output)

	//normalize original document
	normDocument, err := proc.Normalize(document, nOpts)
	if err != nil {
		panic(err)
	}

	tbs := createVerifyHash(normDocument, "normalize", "SHA256", options, proc, nOpts)
	fmt.Printf("HERE IS YOUR HASH: %x \n", tbs)

	signature := getSignature(tbs, privateKey)
	log.Println(signature)
	//log.Println("Here is your tbs:", tbs)
	log.Println("Here is your normalized document:", normDocument.(string))

	log.Println("Here is the output", output)
}

func removeSignature(document map[string]interface{}) {
	if _, ok := document["signature"]; ok {
		log.Println("Removing signature.")
		delete(document, "signature")
	}
}

func copyMap(orig map[string]interface{}) map[string]interface{} {
	// Copy from the original map to the target map
	var newMap = make(map[string]interface{})
	for key, value := range orig {
		newMap[key] = value
	}

	return newMap
}

func getISOTime() string {
	ISOString := "2006-01-02T15:04:05.999Z07:00"
	t := time.Now().UTC().Format(ISOString)

	return t
}

func doHashing(options interface{}) string {
	b, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}

	//hash canonicalized options object
	h1 := sha256.Sum256(b)
	//hash it again
	h2 := sha256.Sum256(b)
	//append second hash to first hash
	h12 := string(h1[:]) + string(h2[:])
	//hash output
	h3 := sha256.Sum256([]byte(h12))

	//return output
	return string(h3[:])
}

/*func sanitizeSigOpts(options *SignatureOptions) {
(interface{}, error)
}*/

func createVerifyHash(document interface{}, canAlgo string, digAlgo string, options SignatureOptions, proc *ld.JsonLdProcessor, nOpts *ld.JsonLdOptions) string {
	opts := SignatureOptions{
		Context: options.Context,
		Creator: options.Creator,
		Created: options.Created,
		Domain:  options.Domain,
		Nonce:   options.Nonce,
	}

	//sanitizeSigOpts(&opts)
	if opts.Created == "" {
		opts.Created = getISOTime()
	}

	normOptions, err := proc.Normalize(opts, nOpts)
	if err != nil {
		panic(err)
	}

	output := doHashing(normOptions)

	return output
}

func getSignature(tbs string, privateKey string) string {
	r := strings.NewReader(privateKey)
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatal(err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		log.Println(block)
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	log.Println(privateKey)
	log.Println(key)

	return "lol"
}
