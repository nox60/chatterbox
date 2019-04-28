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

package chatterbox

import (
	//	"bytes" //un-comment for helpers like bytes.equal
	"encoding/binary"
	"errors"
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

/*
var sender_cache sync.Map

var receiver_cache sync.Map

var counters sync.Map

var global_rootchain_cache sync.Map
*/
//make(map[PublicKey]*PublicKey)

//make(map[chatter.PublicKey](*chatter.Session))

// Chatter represents a chat participant. Each Chatter has a single long-term
// key Identity, and a map of open sessions with other users (indexed by their
// identity keys). You should not need to modify this.
type Chatter struct {
	Identity *KeyPair
	Sessions map[PublicKey]*Session
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
	ReadMessages     map[int]*Message
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

	global_public_keys.Store(c.Identity.PublicKey, &c.Identity.PublicKey)

	return c
}

// EndSession erases all data for a session with the designated partner.
// All outstanding key material should be zeroized and the session erased.
func (c *Chatter) EndSession(partnerIdentity *PublicKey) error {

	if _, exists := c.Sessions[*partnerIdentity]; !exists {
		return errors.New("Don't have that session open to tear down")
	}

	delete(c.Sessions, *partnerIdentity)

	return nil
}

// InitiateHandshake prepares the first message sent in a handshake, containing
// an ephemeral DH share. The partner which initiates should be
// the first to choose a new DH ratchet value. Part of this code has been
// provided for you, you will need to fill in the key derivation code.
func (c *Chatter) InitiateHandshake(partnerIdentity *PublicKey) (*PublicKey, error) {

	DHCombine(&c.Identity.PublicKey, &c.Identity.PrivateKey)

	//KEY_SERVER  = make(map[PublicKey]*PublicKey)

	if _, exists := c.Sessions[*partnerIdentity]; exists {
		return nil, errors.New("Already have session open")
	}

	newKeyPair := NewKeyPair()

	c.Sessions[*partnerIdentity] = &Session{
		StaleReceiveKeys: make(map[int]*SymmetricKey),
		PartnerDHRatchet: partnerIdentity,
		MyDHRatchet:      newKeyPair,
		SendCounter:      0,
		LastUpdate:       0,
		ReceiveCounter:   0,
		RootChain:        nil,
		SendChain:        nil,
		ReceiveChain:     nil,
	}

	return &newKeyPair.PublicKey, nil
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

	bNewPairs := NewKeyPair()

	dh1 := DHCombine(partnerIdentity, &bNewPairs.PrivateKey)
	dh2 := DHCombine(partnerEphemeral, &c.Identity.PrivateKey)
	dh3 := DHCombine(partnerEphemeral, &bNewPairs.PrivateKey)

	fff := CombineKeys(dh1, dh2, dh3)

	//Ab, aB, ab

	rootKey := fff.DeriveKey(HANDSHAKE_CHECK_LABEL)

	c.Sessions[*partnerIdentity] = &Session{
		StaleReceiveKeys: make(map[int]*SymmetricKey),
		PartnerDHRatchet: partnerEphemeral,
		MyDHRatchet:      bNewPairs,
		SendCounter:      0,
		LastUpdate:       0,
		ReceiveCounter:   0,
		RootChain:        fff.DeriveKey(CHAIN_LABEL),
		SendChain:        fff.DeriveKey(CHAIN_LABEL).DeriveKey(KEY_LABEL),
		ReceiveChain:     fff.DeriveKey(CHAIN_LABEL).DeriveKey(KEY_LABEL),
	}

	return &bNewPairs.PublicKey, rootKey, nil
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

	dh2 := DHCombine(partnerIdentity, &c.Sessions[*partnerIdentity].MyDHRatchet.PrivateKey)
	dh1 := DHCombine(partnerEphemeral, &c.Identity.PrivateKey)
	dh3 := DHCombine(partnerEphemeral, &c.Sessions[*partnerIdentity].MyDHRatchet.PrivateKey)
	kkk := CombineKeys(dh1, dh2, dh3)

	myNewKey := NewKeyPair()

	/*
	 *	1. Pick a new DH ratchet key ​gb​ 2
	 *	2. Update his root key by combining with ​ga​ 2·b2
	 *	3. Derive a new sending key chain
	 *	4. Use this to encrypt his message and send it (along with ​g​b2​) to Alice so she can update her
	 *	   root key in the same way and derive the keys needed to decrypt Bob’s message
	 */

	rootKey := kkk.DeriveKey(HANDSHAKE_CHECK_LABEL)

	c.Sessions[*partnerIdentity].PartnerDHRatchet = partnerEphemeral
	c.Sessions[*partnerIdentity].RootChain = kkk.DeriveKey(CHAIN_LABEL)
	c.Sessions[*partnerIdentity].ReceiveChain = kkk.DeriveKey(CHAIN_LABEL).DeriveKey(KEY_LABEL)
	c.Sessions[*partnerIdentity].SendChain = kkk.DeriveKey(CHAIN_LABEL).DeriveKey(KEY_LABEL)
	c.Sessions[*partnerIdentity].MyDHRatchet = myNewKey

	return rootKey, nil

}

// SendMessage is used to send the given plaintext string as a message.
// You'll need to implement the code to ratchet, derive keys and encrypt this message.
func (c *Chatter) SendMessage(partnerIdentity *PublicKey,
	plaintext string) (*Message, error) {

	if _, exists := c.Sessions[*partnerIdentity]; !exists {
		return nil, errors.New("Can't send message to partner with no open session")
	}

	iv := NewIV()

	/*
		sendingDHRatchet := NewKeyPair()

		if c.Sessions[*partnerIdentity].SendCounter == 0 {
			sendingDHRatchet = c.Sessions[*partnerIdentity].MyDHRatchet
		}*/

	//newSendingDH := DHCombine( c.Sessions[*partnerIdentity].PartnerDHRatchet, &nextDHRatchet.PrivateKey )

	c.Sessions[*partnerIdentity].SendCounter = c.Sessions[*partnerIdentity].SendCounter + 1

	message := &Message{
		Sender:        &c.Identity.PublicKey,
		Receiver:      partnerIdentity,
		Ciphertext:    nil,
		IV:            iv,
		Counter:       c.Sessions[*partnerIdentity].SendCounter,
		NextDHRatchet: &c.Sessions[*partnerIdentity].MyDHRatchet.PublicKey,
	}

	data := message.EncodeAdditionalData()

	encrypt := c.Sessions[*partnerIdentity].SendChain

	//fmt.Println("sendcount: " , c.Sessions[*partnerIdentity].SendCounter,"  sendchain: ", encrypt)

	ciphertext := encrypt.AuthenticatedEncrypt(plaintext, data, iv)

	message.Ciphertext = ciphertext

	//newRootChain := c.Sessions[*partnerIdentity].RootChain.DeriveKey(CHAIN_LABEL)

	//newSendingChain := newRootChain.DeriveKey(KEY_LABEL)

	//newReceivingChain := newRootChain.DeriveKey(KEY_LABEL)

	//c.Sessions[*partnerIdentity].RootChain = newRootChain

	c.Sessions[*partnerIdentity].SendChain = c.Sessions[*partnerIdentity].SendChain.DeriveKey(KEY_LABEL)

	//c.Sessions[*partnerIdentity].ReceiveChain = newReceivingChain

	return message, nil
}

// ReceiveMessage is used to receive the given message and return the correct
// plaintext. This method is where most of the key derivation, ratcheting
// and out-of-order message handling logic happens.
func (c *Chatter) ReceiveMessage(message *Message) (string, error) {

	if _, exists := c.Sessions[*message.Sender]; !exists {
		return "", errors.New("Can't receive message from partner with no open session")
	}

	data := message.EncodeAdditionalData()

	for {
		//fmt.Println(">>> message.Counter ", message.Counter)

		if c.Sessions[*message.Sender].ReceiveCounter+1 > message.Counter {
			break
		}

		oldReceiverChain := c.Sessions[*message.Sender].ReceiveChain

		c.Sessions[*message.Sender].StaleReceiveKeys[c.Sessions[*message.Sender].ReceiveCounter+1] = oldReceiverChain

		//c.Sessions[*message.Sender].RootChain = c.Sessions[*message.Sender].RootChain.DeriveKey(CHAIN_LABEL)

		c.Sessions[*message.Sender].ReceiveChain = c.Sessions[*message.Sender].ReceiveChain.DeriveKey(KEY_LABEL)

		//c.Sessions[*message.Sender].SendChain = c.Sessions[*message.Sender].RootChain.DeriveKey(KEY_LABEL)

		c.Sessions[*message.Sender].ReceiveCounter = c.Sessions[*message.Sender].ReceiveCounter + 1

	}

	//fmt.Print(" ======   ",    c.Sessions[*message.Sender].StaleReceiveKeys)

	receiveChain := c.Sessions[*message.Sender].StaleReceiveKeys[message.Counter]

	//fmt.Println("receivecount: " , message.Counter,"  receiveChain: ", receiveChain)

	if c.Sessions[*message.Sender].ReadMessages[message.Counter] != nil {
		//fmt.Println("no  nulll....  ")
		return "", errors.New("error")

	}

	plaintext, err := receiveChain.AuthenticatedDecrypt(message.Ciphertext, data, message.IV)

	//c.Sessions[*message.Sender].StaleReceiveKeys[message.Counter] = NewSymmetricKey()

	if err == nil {

		if c.Sessions[*message.Sender].ReadMessages == nil {
			c.Sessions[*message.Sender].ReadMessages = make(map[int](*Message))
		}

		c.Sessions[*message.Sender].ReadMessages[message.Counter] = message
	}

	//remove the old chain

	return plaintext, err
}
