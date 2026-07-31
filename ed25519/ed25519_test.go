package ed25519

import (
	"bytes"
	"testing"
)

const pemPriv = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIL6GpkRSmMz3oC5ZQDNozTFwmvOWJncIQlgbzX+hdt0F
-----END PRIVATE KEY-----`

func TestLoadPEMAndSign(t *testing.T) {
	priv, err := LoadPrivateKey([]byte(pemPriv))
	if err != nil {
		t.Fatal(err)
	}
	data := []byte("hola")
	sig := priv.Sign(data)
	pub := priv.Public()
	if !pub.Verify(data, sig) {
		t.Fatal("verify failed")
	}
	if pub.Verify([]byte("otro"), sig) {
		t.Fatal("verify should fail on different data")
	}
}

func TestSaveLoadRoundTrip(t *testing.T) {
	priv, pub, err := Generate()
	if err != nil {
		t.Fatal(err)
	}
	if err := priv.SaveToFile("/tmp/opencode/rt.key"); err != nil {
		t.Fatal(err)
	}
	if err := pub.SaveToFile("/tmp/opencode/rt.pub"); err != nil {
		t.Fatal(err)
	}
	priv2, err := LoadPrivateKeyFromFile("/tmp/opencode/rt.key")
	if err != nil {
		t.Fatal(err)
	}
	pub2, err := LoadPublicKeyFromFile("/tmp/opencode/rt.pub")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(priv2.Bytes(), priv.Bytes()) {
		t.Fatal("private key mismatch")
	}
	if !bytes.Equal(pub2.Bytes(), pub.Bytes()) {
		t.Fatal("public key mismatch")
	}
}
