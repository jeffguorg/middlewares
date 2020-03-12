package signature

type SigningMethod interface {
	Verify(signingString, signature string) error // Returns nil if signature is valid
	Sign(signingString string) (string, error)    // Returns encoded signature or error
}
