package auth_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/pem"
	"errors"
	"os"
	"testing"

	"github.com/rchamarthy/chata/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewIdentity(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	henk := auth.GenerateIdentity()

	pkSize := henk.PublicKey().Size()
	assert.Equal(256, pkSize)

	b, e := henk.MarshalText()
	assert.NotNil(b)
	require.NoError(e)
	require.NoError(henk.UnmarshalText(b))

	x := auth.EmptyIdentity()
	b, e = x.MarshalText()
	assert.Nil(b)
	require.Error(e)

	e = x.UnmarshalText([]byte("blah"))
	require.Error(e)

	henkPub := auth.PubIdentity(henk.PublicKey())

	b, e = henkPub.MarshalText()
	require.NoError(e)
	assert.NotNil(b)
	require.NoError(x.UnmarshalText(b))

	assert.NotEmpty(henkPub.String())

	henkPub2 := henkPub.Public()

	b, e = henkPub2.MarshalText()
	require.NoError(e)
	assert.NotNil(b)
	require.NoError(x.UnmarshalText(b))
}

func TestUnmarshalErrors(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	require := require.New(t)

	badType := pem.EncodeToMemory(&pem.Block{
		Type:    "NULL KEY",
		Headers: nil,
		Bytes:   nil,
	})
	badPriv := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   nil,
	})
	badPub := pem.EncodeToMemory(&pem.Block{
		Type:    "RSA PUBLIC KEY",
		Headers: nil,
		Bytes:   nil,
	})

	assert.Empty(auth.EmptyIdentity().String())

	henk := auth.GenerateIdentity()

	pkSize := henk.PublicKey().Size()
	assert.Equal(256, pkSize)

	x := auth.EmptyIdentity()
	b, e := x.MarshalText()
	assert.Nil(b)
	require.Error(e)

	henkPub := auth.PubIdentity(henk.PublicKey())

	b, e = henkPub.MarshalText()
	require.NoError(e)
	assert.NotNil(b)
	require.Error(x.UnmarshalText(badPub))

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	henkPriv := auth.NewIdentity(priv)

	b, e = henkPriv.MarshalText()
	require.NoError(e)
	assert.NotNil(b)
	require.Error(x.UnmarshalText(badPriv))

	require.Error(x.UnmarshalText(badType))
}

func TestEncrypt(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	henk := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	jaap := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	ingrid := auth.NewIdentity(priv)

	msg := []byte("Arme mensen kunnen niet met geld omgaan: ze geven alles uit aan eten en kleren, " +
		"terwijl rijke mensen het heel verstandig op de bank zetten.")

	// Lets encrypt it using Ingrid's public key.
	henksMessage, err := henk.Encrypt(msg, ingrid.PublicKey())
	require.NoError(err)

	jaapsMessage, err := jaap.Encrypt(msg, ingrid.PublicKey())
	require.NoError(err)

	// Decrypt
	hm, _ := ingrid.Decrypt(henksMessage)
	jm, _ := ingrid.Decrypt(jaapsMessage)

	// Compare the messages of Henk and Jaap, and the original
	assert.True(bytes.Equal(hm, jm))
	assert.True(bytes.Equal(hm, msg))
}

func TestEncryptionNeverTheSame(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	// Even when using the same public key, the encrypted messages are never the same
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	henk := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	jaap := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	joop := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	koos := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	kees := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	erik := auth.NewIdentity(priv)

	// Added a couple of the same Identities at the end, just to prove that the
	// encrypted outcome differs each time.
	identities := []*auth.Identity{henk, jaap, joop, koos, kees, erik, erik, erik, erik}

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	ingrid := auth.NewIdentity(priv)

	msg := []byte("Aan ons land geen polonaise.")
	msgs := make([][]byte, 0, len(identities))

	for _, id := range identities {
		// encrypt the message using Ingrid her public key
		e, _ := id.Encrypt(msg, ingrid.PublicKey())
		msgs = append(msgs, e)
	}

	s := []byte("start")
	for _, m := range msgs {
		assert.False(bytes.Equal(m, s))
	}
}

func TestEncryptDecrypt(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	henk := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	ingrid := auth.NewIdentity(priv)

	// a message from Henk to Ingrid
	msg := []byte("Die uitkeringstrekkers pikken al onze banen in.")
	// Lets encrypt it, we want to sent it to Ingrid, thus, we use her public key.
	encryptedMessage, err := henk.Encrypt(msg, ingrid.PublicKey())
	require.NoError(err)

	// Decrypt Message
	plainTextMessage, err := ingrid.Decrypt(encryptedMessage)
	require.NoError(err)

	assert.True(bytes.Equal(plainTextMessage, msg))
}

func TestEncryptDecryptMyself(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	// If anyone, even you, encrypts (id.e. “locks”) something with your public-key,
	// only you can decrypt it (id.e. “unlock” it) with your secret, private key.
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	henk := auth.NewIdentity(priv)

	// a message from Henk
	msg := []byte("Subsidized, dat is toch iets dat je krijgt als je eigenlijk niet goed genoeg bent?")

	// Lets encrypt it, we want to sent it to self, thus, we need our public key.
	encryptedMessage, err := henk.Encrypt(msg, nil)
	require.NoError(err)

	// Decrypt Message
	plainTextMessage, err := henk.Decrypt(encryptedMessage)
	require.NoError(err)

	assert.True(bytes.Equal(plainTextMessage, msg))
}

func TestSignVerify(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	henk := auth.NewIdentity(priv)

	// A public message from Henk.
	// note that the message is a byte array, not just a string.
	msg := []byte("Wilders doet tenminste iets tegen de politiek.")

	// Henk signs the message with his private key. This will show the recipient
	// proof that this message is indeed from Henk
	sig, _ := henk.Sign(msg)

	// now, if the message msg is public, anyone can read it.
	// the signature sig however, proves this message is from Henk.
	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	ingrid := auth.NewIdentity(priv)

	priv, err = rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	hans := auth.NewIdentity(priv)

	err = ingrid.Verify(msg, sig, henk.PublicKey())
	require.NoError(err)

	err = hans.Verify(msg, sig, henk.PublicKey())
	require.NoError(err)

	// Let's see if we can break the signature verification
	// (1) changing the message
	err = hans.Verify([]byte("Wilders is een opruier"), sig, henk.PublicKey())
	require.Error(err)

	// (2) changing the signature
	err = hans.Verify(msg, []byte("I am not the signature"), henk.PublicKey())
	require.Error(err)

	// (3) changing the public key
	err = hans.Verify(msg, sig, ingrid.PublicKey())
	require.Error(err)

	_, err = ingrid.Public().Sign([]byte("test"))
	assert.Equal(auth.ErrNoPrivateKey, err)

	require.Error(hans.Verify(msg, []byte("whatever"), nil))
}

func TestLoad(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	require := require.New(t)

	id, err := auth.LoadIdentity("")
	assert.Nil(id)
	require.Error(err)

	id, err = auth.LoadIdentity("rsa_test.go")
	assert.Nil(id)
	require.Error(err)

	// Generate a key and save it
	key := auth.GenerateIdentity()
	require.NoError(key.SaveIdentity("./test.key"))

	id, err = auth.LoadIdentity("./test.key")
	require.NoError(err)
	assert.NotNil(id)
	require.NoError(os.RemoveAll("./test.key"))

	require.NoError(id.SaveIdentity("new_copy.key"))
	require.NoError(os.RemoveAll("new_copy.key"))

	id = auth.EmptyIdentity()
	require.Error(id.SaveIdentity("/abcd"))
}

func TestPanicOnError(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	assert.Panics(func() {
		auth.PanicOnError(errors.New("error"))
	})

	assert.NotPanics(func() {
		auth.PanicOnError(nil)
	})
}
