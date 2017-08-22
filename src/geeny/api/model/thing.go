package model

// Thing represents the data type of a thing
type Thing struct {
	ID          string     `json:"id"`
	ThingTypeID string     `json:"thingTypeId"`
	PairingCode string     `json:"pairingCode"`
	Attributes  Attributes `json:"attributes"`
	CreatedAt   string     `json:"createdAt"`
	ModifiedAt  string     `json:"modifiedAt"`
}

// Attributes represents the data type of thing attributes
type Attributes struct {
	IOT     IOT         `json:"iot"`
	Owner   Owner       `json:"owner"`
	Creator interface{} `json:"creator"`
}

// IOT represents the data type of an internet of things thing
type IOT struct {
	Certificate Certificate `json:"certificate"`
	IOTEndpoint string      `json:"iotEndpoint"`
}

// Owner represents the data type of a thing owner
type Owner struct {
	IdentityID string `json:"identityId"`
}

// Certificate represents the data type of a credentials certificate
type Certificate struct {
	KeyPair        KeyPair `json:"keyPair"`
	CertificateID  string  `json:"certificateId"`
	CertificatePem string  `json:"certificatePem"`
}

// KeyPair represents a public / private keypair
type KeyPair struct {
	PublicKey  string `json:"PublicKey"`
	PrivateKey string `json:"PrivateKey"`
}

// - ValidationInterface

// IsValid validates the data structure
func (t *Thing) IsValid() bool {
	return ( len(t.ID) > 0 &&
		len(t.ThingTypeID) > 0 &&
		len(t.PairingCode) > 0 &&
		t.Attributes.IsValid() &&
		len(t.CreatedAt) > 0 &&
		len(t.ModifiedAt) > 0)
}

// IsValid validates the data structure
func (a *Attributes) IsValid() bool {
	return (a.IOT.IsValid() &&
		a.Owner.IsValid() &&
		a.Creator != nil)
}

// IsValid validates the data structure
func (i *IOT) IsValid() bool {
	return (i.Certificate.IsValid() &&
		len(i.IOTEndpoint) > 0)
}

// IsValid validates the data structure
func (o *Owner) IsValid() bool {
	return true
}

// IsValid validates the data structure
func (c *Certificate) IsValid() bool {
	return (c.KeyPair.IsValid() &&
		len(c.CertificateID) > 0 &&
		len(c.CertificatePem) > 0)
}

// IsValid validates the data structure
func (k *KeyPair) IsValid() bool {
	return (len(k.PublicKey) > 0 &&
		len(k.PrivateKey) > 0)
}
