package main

import (
	"chatter"
	"fmt"
	"math"
	"math/big"
)


func main() {
	//fmt.Println("aaaaaaa")
	//testpackage.Abc()
	//fmt.Println(chatter.NewSymmetricKey())
	//fmt.Println(chatter.RandomnessSource())
	tempObj := new(chatter.KeyPair)

	tempObj = chatter.NewKeyPair()

	fmt.Println("new chatter ", chatter.NewChatter().Identity.PublicKey.X, chatter.NewChatter().Identity.PublicKey.Y)

	fmt.Println(tempObj.PrivateKey)

	//rect2ee := &chatter.Chatter{}

	recc2ee := new(chatter.Chatter)

	recc2ee.Sessions = make(map[chatter.PublicKey](*chatter.Session)) //注意写法

	userDBww := make(map[chatter.PublicKey](*chatter.Session)) //注意写法

	yinzhengjie := map[string]int{
		"尹正杰": 18,
		"饼干":  20,
	}

	yinzhengjie["eee"] = 98

	fmt.Println(yinzhengjie)

	//recc2ee.Sessions[nil] = nil

	tempPublicKey := new(chatter.PublicKey)

	tempPublicKey.Y = new(big.Int).SetUint64(uint64(1000))

	tempPublicKey.X = new(big.Int).SetUint64(uint64(9999))

	tempSession := new(chatter.Session)

	userDBww[*tempPublicKey] = tempSession

	recc2ee.Sessions[*tempPublicKey] = tempSession

	//ß	recc2ee.Sessions[&tempPublicKey] = tempSession

	fmt.Println(recc2ee.Sessions)

	fmt.Println(math.Pow(3, 4)) //次方运算

	//首先模拟素数和本原根

	g := float64(5)

	p := float64(23)

	A := float64(4)

	B := float64(3)

	//fmt.Println(math.Pow(g,A))

	Ya := math.Mod(math.Pow(g, A), p)

	fmt.Println("Ya:", Ya) //Ya,

	Yb := math.Mod(math.Pow(g, B), p)

	fmt.Println("Yb:", Yb) //Yb,

	Ka := math.Mod(math.Pow(Yb, A), p)

	fmt.Println(Ka)

	Kb := math.Mod(math.Pow(Ya, B), p)

	fmt.Println(Kb)

	//fmt.Println(9%7) //取模

	
	/*
		func DoHandshake(t *testing.T, alice, bob *Chatter) error {

			if VERBOSE {
			fmt.Println("Starting handshake sequence")
			fmt.Printf("Initiator identity: %s\n", PrintHandle(&alice.Identity.PublicKey))
			fmt.Printf("Responder identity: %s\n", PrintHandle(&bob.Identity.PublicKey))
		}

			aliceShare, err := alice.InitiateHandshake(&bob.Identity.PublicKey)
			if err != nil {
			t.Logf("Error initiating handshake")
			return err
		}

			if VERBOSE {
			fmt.Printf("Initiator sends ephemeral key: %X\n", aliceShare.Fingerprint())
		}

			bobShare, bobCheck, err := bob.ReturnHandshake(&alice.Identity.PublicKey, aliceShare)
			if err != nil {
			t.Logf("Error responding to handshake")
			return err
		}
			if VERBOSE {
			fmt.Printf("Responder sends ephemeral key: %X\n", bobShare.Fingerprint())
		}

			aliceCheck, err := alice.FinalizeHandshake(&bob.Identity.PublicKey, bobShare)
			if err != nil {
			t.Logf("Error finalizing handshake")
			return err
		}

			if !bytes.Equal(aliceCheck.Key, bobCheck.Key) {
			t.Logf("Handshake participants don't agree on master key")
			return errors.New("Handshake failed")
		}
			if VERBOSE {
			fmt.Printf("Handshake master key hash: %X\n", bobCheck.Key)
		}
			return err
		}*/
}

//a的n次方
//超出uint64的部分会丢失
func exponent(a, n uint64) uint64 {
	result := uint64(1)
	for i := n; i > 0; i >>= 1 {
		if i&1 != 0 {
			result *= a
		}
		a *= a
	}
	return result
}


