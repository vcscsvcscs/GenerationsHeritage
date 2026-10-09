package main

import (
	"errors"
	"fmt"
	"io"

	"filippo.io/age"
)

const ageSuffix = ".age"

type nopWriteCloser struct{ io.Writer }

func (nopWriteCloser) Close() error { return nil }

// recipient returns the age recipient configured for encryption, or nil when encryption is disabled.
func (c *config) recipient() (age.Recipient, error) {
	switch {
	case c.ageRecipient != "":
		r, err := age.ParseX25519Recipient(c.ageRecipient)
		if err != nil {
			return nil, fmt.Errorf("BACKUP_AGE_RECIPIENT: %w", err)
		}

		return r, nil
	case c.passphrase != "":
		return age.NewScryptRecipient(c.passphrase)
	}

	return nil, nil //nolint:nilnil // nil recipient means encryption is disabled
}

func (c *config) identity() (age.Identity, error) {
	switch {
	case c.ageIdentity != "":
		i, err := age.ParseX25519Identity(c.ageIdentity)
		if err != nil {
			return nil, fmt.Errorf("BACKUP_AGE_IDENTITY: %w", err)
		}

		return i, nil
	case c.passphrase != "":
		return age.NewScryptIdentity(c.passphrase)
	}

	return nil, errors.New("object is encrypted: set BACKUP_AGE_IDENTITY or BACKUP_ENCRYPTION_PASSPHRASE")
}

func encryptWriter(w io.Writer, r age.Recipient) (io.WriteCloser, error) {
	if r == nil {
		return nopWriteCloser{w}, nil
	}

	return age.Encrypt(w, r)
}
