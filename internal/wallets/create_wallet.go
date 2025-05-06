package wallets

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/ecdsa"
	"crypto/rand"
	"fmt"
	"github.com/knstch/subtrack-libs/tracing"

	"github.com/ethereum/go-ethereum/crypto"
)

func (svc *ServiceImpl) CreateWallet(ctx context.Context, userID uint) error {
	ctx, span := tracing.StartSpan(ctx, "service: CreateWallet")
	defer span.End()

	wallet, err := generateWallet()
	if err != nil {
		return fmt.Errorf("generateWallet: %w", err)
	}

	if err = wallet.encryptPrivateKey(svc.walletSecret); err != nil {
		return fmt.Errorf("encryptPrivateKey: %w", err)
	}

	if err = svc.repo.CreateWallet(ctx, userID, wallet.PublicKey, wallet.PrivateKey); err != nil {
		return fmt.Errorf("repo.CreateWallet: %w", err)
	}

	return nil
}

func generateWallet() (Wallet, error) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		return Wallet{}, fmt.Errorf("crypto.GenerateKey: %w", err)
	}

	privBytes := crypto.FromECDSA(privateKey)

	publicKey := privateKey.Public().(*ecdsa.PublicKey)

	address := crypto.PubkeyToAddress(*publicKey).Hex()

	return Wallet{
		PublicKey:  address,
		PrivateKey: privBytes,
	}, nil
}

func (w *Wallet) encryptPrivateKey(key string) error {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return fmt.Errorf("wallet.encryptPrivateKey: %w", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return fmt.Errorf("wallet.encryptPrivateKey: %w", err)
	}
	nonce := make([]byte, aesGCM.NonceSize())
	_, _ = rand.Read(nonce)

	ciphertext := aesGCM.Seal(nonce, nonce, w.PrivateKey, nil)

	w.PrivateKey = ciphertext

	return nil
}

func (w *Wallet) decryptPrivateKey(key string) (*ecdsa.PrivateKey, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, fmt.Errorf("aes.NewCiphe: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGC: %w", err)
	}

	nonceSize := aesGCM.NonceSize()
	if len(w.PrivateKey) < nonceSize {
		return nil, fmt.Errorf("decryptPrivateKey: ciphertext too short")
	}

	nonce, ciphertext := w.PrivateKey[:nonceSize], w.PrivateKey[nonceSize:]

	plainPrivateKeyBytes, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("aesGCM.Ope: %w", err)
	}

	privateKey, err := crypto.ToECDSA(plainPrivateKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA: %w", err)
	}

	return privateKey, nil
}
