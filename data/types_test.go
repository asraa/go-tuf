package data

import (
	"encoding/json"

	. "gopkg.in/check.v1"
)

const (
	// This public key is from the TUF specs:
	//
	// https://github.com/theupdateframework/specification
	//
	public       = `"72378e5bc588793e58f81c8533da64a2e8f1565c1fcc7f253496394ffc52542c"`
	keyid10      = "1bf1c6e3cdd3d3a8420b19199e27511999850f4b376c4547b2f32fba7e80fca3"
	keyid10algos = "506a349b85945d0d99c7289c3f0f1f6c550218089d1d38a3f64824db31e827ac"
)

type TypesSuite struct{}

var _ = Suite(&TypesSuite{})

func (TypesSuite) TestKeyIDs(c *C) {
	var hexbytes HexBytes
	err := json.Unmarshal([]byte(public), &hexbytes)
	c.Assert(err, IsNil)

	pubBytes, err := json.Marshal(KeyValue{Public: hexbytes})
	c.Assert(err, IsNil)
	key := Key{
		Type:   KeyTypeEd25519,
		Scheme: KeySchemeEd25519,
		Value:  pubBytes,
	}
	c.Assert(key.IDs(), DeepEquals, []string{keyid10})

	key = Key{
		Type:       KeyTypeEd25519,
		Scheme:     KeySchemeEd25519,
		Algorithms: KeyAlgorithms,
		Value:      pubBytes,
	}
	c.Assert(key.IDs(), DeepEquals, []string{keyid10algos})
}

func (TypesSuite) TestRootAddKey(c *C) {
	var hexbytes HexBytes
	err := json.Unmarshal([]byte(public), &hexbytes)
	c.Assert(err, IsNil)

	pubBytes, err := json.Marshal(KeyValue{Public: hexbytes})
	c.Assert(err, IsNil)
	key := &Key{
		Type:   KeyTypeEd25519,
		Scheme: KeySchemeEd25519,
		Value:  pubBytes,
	}

	root := NewRoot()

	c.Assert(root.AddKey(key), Equals, true)
	c.Assert(root.AddKey(key), Equals, false)
}

func (TypesSuite) TestRoleAddKeyIDs(c *C) {
	var hexbytes HexBytes
	err := json.Unmarshal([]byte(public), &hexbytes)
	c.Assert(err, IsNil)

	pubBytes, err := json.Marshal(KeyValue{Public: hexbytes})
	c.Assert(err, IsNil)
	key := &Key{
		Type:   KeyTypeEd25519,
		Scheme: KeySchemeEd25519,
		Value:  pubBytes,
	}

	role := &Role{}
	c.Assert(role.KeyIDs, HasLen, 0)

	c.Assert(role.AddKeyIDs(key.IDs()), Equals, true)
	c.Assert(role.KeyIDs, DeepEquals, []string{keyid10})

	// Adding the key again doesn't modify the array.
	c.Assert(role.AddKeyIDs(key.IDs()), Equals, false)
	c.Assert(role.KeyIDs, DeepEquals, []string{keyid10})

	// Add another key.
	key = &Key{
		Type:       KeyTypeEd25519,
		Scheme:     KeySchemeEd25519,
		Algorithms: KeyAlgorithms,
		Value:      pubBytes,
	}

	// Adding the key again doesn't modify the array.
	c.Assert(role.AddKeyIDs(key.IDs()), Equals, true)
	c.Assert(role.KeyIDs, DeepEquals, []string{keyid10, keyid10algos})
}
