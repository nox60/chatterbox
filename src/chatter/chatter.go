// Implementation of a forward-secure, end-to-end encrypted messaging client
// supporting key compromise recovery and out-of-order message delivery.
// Directly inspired by Signal/Double-ratchet protocol but missing a few
// features. No asynchronous handshake support (pre-keys) for example.
//
// SECURITY WARNING: This code is meant for educational purposes and may
// contain vulnerabilities or other bugs. Please do not use it for
// security-critical applications.
//
// GRADING NOTES: This is the only file you need to modify for this assignment.
// You may add additional support files if desired. You should modify this file
// to implement the intended protocol, but preserve the function signatures
// for the following methods to ensure your implementation will work with
// standard test code:
//
// *NewChatter
// *EndSession
// *InitiateHandshake
// *ReturnHandshake
// *FinalizeHandshake
// *SendMessage
// *ReceiveMessage
//
// In addition, you'll need to keep all of the following structs' fields:
//
// *Chatter
// *Session
// *Message
//
// You may add fields if needed (not necessary) but don't rename or delete
// any existing fields.
//
// Original version
// Joseph Bonneau February 2019

package chatter

import (
	//	"bytes" //un-comment for helpers like bytes.equal
	"encoding/binary"
	"errors"
	"fmt"

	//"fmt"
	"sync"

	//	"fmt" //un-comment if you want to do any debug printing.
)

// Labels for key derivation

// Label for generating a check key from the initial root.
// Used for verifying the results of a handshake out-of-band.
const HANDSHAKE_CHECK_LABEL byte = 0x01

// Label for ratcheting the main chain of keys
const CHAIN_LABEL = 0x02

// Label for deriving message keys from chain keys.
const KEY_LABEL = 0x03

//var KEY_SERVER map[PublicKey]*PublicKey

var global_public_keys sync.Map

var partner_public_keys_cache sync.Map

var my_private_keys_cache sync.Map

var msg_counter sync.Map

//make(map[PublicKey]*PublicKey)

//make(map[chatter.PublicKey](*chatter.Session))

// Chatter represents a chat participant. Each Chatter has a single long-term
// key Identity, and a map of open sessions with other users (indexed by their
// identity keys). You should not need to modify this.
type Chatter struct {
	Identity *KeyPair
	Sessions map[PublicKey]*Session
}

type MessagePublicKey struct {
	counter int
	publicKey *PublicKey
}

// Session represents an open session between one chatter and another.
// You should not need to modify this, though you can add additional fields
// if you want to.
type Session struct {
	MyDHRatchet      *KeyPair
	PartnerDHRatchet *PublicKey
	RootChain        *SymmetricKey
	SendChain        *SymmetricKey
	ReceiveChain     *SymmetricKey
	StaleReceiveKeys map[int]*SymmetricKey
	SendCounter      int
	LastUpdate       int
	ReceiveCounter   int
}

// Message represents a message as sent over an untrusted network.
// The first 5 fields are send unencrypted (but should be authenticated).
// The ciphertext contains the (encrypted) communication payload.
// You should not need to modify this.
type Message struct {
	Sender        *PublicKey
	Receiver      *PublicKey
	NextDHRatchet *PublicKey
	Counter       int
	LastUpdate    int
	Ciphertext    []byte
	IV            []byte
}

// EncodeAdditionalData encodes all of the non-ciphertext fields of a message
// into a single byte array, suitable for use as additional authenticated data
// in an AEAD scheme. You should not need to modify this code.
func (m *Message) EncodeAdditionalData() []byte {
	buf := make([]byte, 8+3*FINGERPRINT_LENGTH)

	binary.LittleEndian.PutUint32(buf, uint32(m.Counter))
	binary.LittleEndian.PutUint32(buf[4:], uint32(m.LastUpdate))

	if m.Sender != nil {
		copy(buf[8:], m.Sender.Fingerprint())
	}
	if m.Receiver != nil {
		copy(buf[8+FINGERPRINT_LENGTH:], m.Receiver.Fingerprint())
	}
	if m.NextDHRatchet != nil {
		copy(buf[8+2*FINGERPRINT_LENGTH:], m.NextDHRatchet.Fingerprint())
	}

	return buf
}

// NewChatter creates and initializes a new Chatter object. A long-term
// identity key is created and the map of sessions is initialized.
// You should not need to modify this code.
func NewChatter() *Chatter {
	c := new(Chatter)
	c.Identity = NewKeyPair()
	c.Sessions = make(map[PublicKey]*Session)

	global_public_keys.Store(c.Identity.PublicKey,&c.Identity.PublicKey)

	return c
}

// EndSession erases all data for a session with the designated partner.
// All outstanding key material should be zeroized and the session erased.
func (c *Chatter) EndSession(partnerIdentity *PublicKey) error {

	if _, exists := c.Sessions[*partnerIdentity]; !exists {
		return errors.New("Don't have that session open to tear down")
	}

	// TODO: your code here
	delete(c.Sessions, *partnerIdentity)

	return nil
}

// InitiateHandshake prepares the first message sent in a handshake, containing
// an ephemeral DH share. The partner which initiates should be
// the first to choose a new DH ratchet value. Part of this code has been
// provided for you, you will need to fill in the key derivation code.
func (c *Chatter) InitiateHandshake(partnerIdentity *PublicKey) (*PublicKey, error) {

	//KEY_SERVER  = make(map[PublicKey]*PublicKey)

	if _, exists := c.Sessions[*partnerIdentity]; exists {
		return nil, errors.New("Already have session open")
	}

	//User partner's public key and own private key to DH process
	sender := DHCombine(partnerIdentity, &c.Identity.PrivateKey)

	//fmt.Println( sender.DeriveKey(HANDSHAKE_CHECK_LABEL))
	rootChain := CombineKeys(sender)

	c.Sessions[*partnerIdentity] = &Session{
		StaleReceiveKeys: make(map[int]*SymmetricKey),
		// TODO: your code here
		PartnerDHRatchet: partnerIdentity,
		MyDHRatchet: c.Identity,
		SendCounter   :   0,
		LastUpdate    :   0,
		ReceiveCounter :  0,
		RootChain: rootChain,
		SendChain: sender,
		ReceiveChain: sender,
	}

	// TODO: your code

	global_public_keys.Store(c.Identity.PublicKey,&c.Identity.PublicKey)

	return &c.Identity.PublicKey, nil
}

// ReturnHandshake prepares the first message sent in a handshake, containing
// an ephemeral DH share. Part of this code has been provided for you, you will
// need to fill in the key derivation code. The partner which calls this
// method is the responder.
func (c *Chatter) ReturnHandshake(partnerIdentity,
	partnerEphemeral *PublicKey) (*PublicKey, *SymmetricKey, error) {

	if _, exists := c.Sessions[*partnerIdentity]; exists {
		return nil, nil, errors.New("Already have session open")
	}

	receiver := DHCombine(partnerEphemeral, &c.Identity.PrivateKey)

	rootChain := CombineKeys(receiver)

	c.Sessions[*partnerIdentity] = &Session{
		StaleReceiveKeys: make(map[int]*SymmetricKey),
		// TODO: your code here
		PartnerDHRatchet: partnerIdentity,
		MyDHRatchet: c.Identity,
		SendCounter   :   0,
		LastUpdate    :   0,
		ReceiveCounter :  0,
		RootChain: rootChain,
		SendChain: receiver,
		ReceiveChain: receiver,
	}

	global_public_keys.Store(c.Identity.PublicKey,&c.Identity.PublicKey)

	// TODO: your code here
	return &c.Identity.PublicKey, receiver.DeriveKey(HANDSHAKE_CHECK_LABEL), nil
}

// FinalizeHandshake lets the initiator receive the responder's ephemeral key
// and finalize the handshake. Part of this code has been provided, you will
// need to fill in the key derivation code. The partner which calls this
// method is the initiator.
func (c *Chatter) FinalizeHandshake(partnerIdentity,
	partnerEphemeral *PublicKey) (*SymmetricKey, error) {

	if _, exists := c.Sessions[*partnerIdentity]; !exists {
		return nil, errors.New("Can't finalize session, not yet open")
	}

	// TODO: your code here
	b1 := DHCombine(partnerEphemeral, &c.Identity.PrivateKey)

	//return nil, errors.New("Not implemented")
	return b1.DeriveKey(HANDSHAKE_CHECK_LABEL), nil

}

// SendMessage is used to send the given plaintext string as a message.
// You'll need to implement the code to ratchet, derive keys and encrypt this message.
func (c *Chatter) SendMessage(partnerIdentity *PublicKey,
	plaintext string) (*Message, error) {

	if _, exists := c.Sessions[*partnerIdentity]; !exists {
		return nil, errors.New("Can't send message to partner with no open session")
	}

	//encode message
	data := []byte("extra")

	iv := NewIV()

	v, _ :=  global_public_keys.Load(*partnerIdentity)

	tt :=  v.(*PublicKey)

	dhForEnCrypt := DHCombine(tt, &c.Sessions[*partnerIdentity].MyDHRatchet.PrivateKey)


	//temppp, _ := global_public_keys.Load(c.Identity.PublicKey)

	//fmt.Println("     -------------------->>  send    ", plaintext, "partnerIdentity      ",    partnerIdentity, "  receiver encode public key :", tt,       "             c.Identity.PublicKey"  , c.Identity.PublicKey   ," sender global_public_keys.Load(c.Identity.PublicKey) cached",  temppp.(*PublicKey)  , "c.Sessions[*partnerIdentity].MyDHRatchet.PublicKey:    ",c.Sessions[*partnerIdentity].MyDHRatchet.PublicKey   ,"            sender private key ()c.Sessions[*partnerIdentity].MyDHRatchet.PrivateKey: ", &c.Sessions[*partnerIdentity].MyDHRatchet.PrivateKey)

	ciphertext := dhForEnCrypt.AuthenticatedEncrypt(plaintext, data, iv)

	newKeyPair := NewKeyPair()

	//Get counter from cache
	counter, _ := msg_counter.Load(*partnerIdentity)

	if nil == counter {
		msg_counter.Store(*partnerIdentity, 0)
	}

	counter, _ = msg_counter.Load(*partnerIdentity)

	newCounter := counter.(int)

	newCounter += 1

	msg_counter.Store(*partnerIdentity, newCounter)

	message := &Message{
		Sender:  &c.Identity.PublicKey,
		Receiver: partnerIdentity,
		// TODO: your code here
		Ciphertext: ciphertext,
		IV: iv,
		Counter: newCounter,
		NextDHRatchet: &newKeyPair.PublicKey,
	}

	//Put public key into cache

	publickKeyAndCount := new(MessagePublicKey{
		1,
		&newKeyPair.PublicKey,
	})

	fmt.Println(publickKeyAndCount)

	// TODO: your code here

	return message, nil
}



// ReceiveMessage is used to receive the given message and return the correct
// plaintext. This method is where most of the key derivation, ratcheting
// and out-of-order message handling logic happens.
func (c *Chatter) ReceiveMessage(message *Message) (string, error) {

	if _, exists := c.Sessions[*message.Sender]; !exists {
		return "", errors.New("Can't receive message from partner with no open session")
	}

	// TODO: your code here
	data := []byte("extra")

	v, _ :=  global_public_keys.Load(*message.Sender)

	tt :=  v.(*PublicKey)

	theCurrentDh := DHCombine(tt, &c.Sessions[*message.Sender].MyDHRatchet.PrivateKey)

	plaintext,err := theCurrentDh.AuthenticatedDecrypt(message.Ciphertext, data, message.IV)


	if message.Counter > 0{
		return "", errors.New("error of message counter")
	}

	//fmt.Println(">>>>>>    received    Ciphertext", message.Ciphertext, "     plaintext     ", plaintext, "      *message.Sender     ",*message.Sender,  "      sender public key  : ", v ,"     my private key: ", &c.Sessions[*message.Sender].MyDHRatchet.PrivateKey)

	if len(plaintext) == 0 {
		return "", errors.New("error of message body")
	}

	theNewKeyPair := NewKeyPair()

	c.Sessions[*message.Sender].MyDHRatchet = theNewKeyPair

	for key := range c.Sessions {
		c.Sessions[key].MyDHRatchet = theNewKeyPair
	}

	global_public_keys.Store(c.Identity.PublicKey, &theNewKeyPair.PublicKey)

	return plaintext, err
}

