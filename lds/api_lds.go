package lds

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/kc1116/json-gold/ld"
)

//Signature . . .
//Struct that defines a signature object based on the https://web-payments.org/specs/source/ld-signatures/#linked-data-signature-overview spec Crowdbotics
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
func S2015(document map[string]interface{}, options SignatureOptions, privateKey string) map[string]interface{} {
	var proc = ld.NewJsonLdProcessor()
	var nOpts = ld.NewJsonLdOptions("")
	nOpts.Format = "application/nquads"

	//make copy of original document
	output := copyMap(document)
	//remove any signature nodes from output
	removeSignature(output)

	//normalize original document
	normDocument, err := proc.Normalize(document, nOpts)
	if err != nil {
		panic(err)
	}

	//create hash to be signed
	tbs := createVerifyHash(normDocument, "normalize", "SHA512", &options, proc, nOpts)

	//sign hash
	signatureValue := getSignature(tbs, privateKey)

	signature := Signature{
		Type:           "LinkedDataSignature2015",
		Creator:        options.Creator,
		Created:        options.Created,
		Domain:         options.Domain,
		Nonce:          options.Nonce,
		SignatureValue: signatureValue,
	}

	output["signature"] = signature

	return output
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

func doHashing(options interface{}, document interface{}) string {
	optsBytes, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}

	docBytes, err := json.Marshal(document)
	if err != nil {
		panic(err)
	}

	//hash canonicalized options object
	h1 := sha256.Sum256(optsBytes)
	//hash it again
	h2 := sha256.Sum256(docBytes)
	//append hash of normalized document to hash of normalized options document
	h12 := string(h1[:]) + string(h2[:])
	//hash output
	h3 := sha256.Sum256([]byte(h12))

	//return output
	return string(h3[:])
}

func createVerifyHash(document interface{}, canAlgo string, digAlgo string, options *SignatureOptions, proc *ld.JsonLdProcessor, nOpts *ld.JsonLdOptions) string {

	//get ISO8601 combined date and time string
	options.Created = getISOTime()
	opts := SignatureOptions{
		Context: options.Context,
		Creator: options.Creator,
		Created: options.Created,
		Domain:  options.Domain,
		Nonce:   options.Nonce,
	}

	//normalize options
	normOptions, err := proc.Normalize(opts, nOpts)
	if err != nil {
		panic(err)
	}

	//hash options and return data to be signed
	output := doHashing(normOptions, document)

	return output
}

func getSignature(tbs string, privateKey string) string {
	r := strings.NewReader(privateKey)
	pemBytes, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}

	block, _ := pem.Decode(pemBytes)
	if block == nil {
		log.Println(block)
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}

	//sign to be signed hash
	var h crypto.Hash
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, h, []byte(tbs))
	if err != nil {
		panic(err)
	}

	//base64 encode signature
	b64sig := base64.StdEncoding.EncodeToString(signature)

	return b64sig
}
