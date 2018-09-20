// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package keystore

import (
	"github.com/sero-cash/go-sero"
	"github.com/sero-cash/go-sero/accounts"
	"github.com/sero-cash/go-sero/core/types"
	"github.com/sero-cash/go-sero/zero/txs/tx"
	"github.com/sero-cash/go-sero/core/state"
	"github.com/sero-cash/go-sero/zero/txs"
	"github.com/sero-cash/go-sero/common"
	"github.com/sero-cash/go-sero/crypto/sha3"
	"github.com/sero-cash/go-sero/rlp"
	"github.com/sero-cash/go-czero-import/keys"
	"github.com/sero-cash/go-sero/zero/utils"
)

// keystoreWallet implements the accounts.Wallet interface for the original
// keystore.
type keystoreWallet struct {
	account  accounts.Account // Single account contained in this wallet
	keystore *KeyStore        // Keystore where the account originates from
}

// URL implements accounts.Wallet, returning the URL of the account within.
func (w *keystoreWallet) URL() accounts.URL {
	return w.account.URL
}

// Status implements accounts.Wallet, returning whether the account held by the
// keystore wallet is unlocked or not.
func (w *keystoreWallet) Status() (string, error) {
	w.keystore.mu.RLock()
	defer w.keystore.mu.RUnlock()

	if _, ok := w.keystore.unlocked[w.account.Address]; ok {
		return "Unlocked", nil
	}
	return "Locked", nil
}

// Open implements accounts.Wallet, but is a noop for plain wallets since there
// is no connection or decryption step necessary to access the list of accounts.
func (w *keystoreWallet) Open(passphrase string) error { return nil }

// Close implements accounts.Wallet, but is a noop for plain wallets since is no
// meaningful open operation.
func (w *keystoreWallet) Close() error { return nil }

// Accounts implements accounts.Wallet, returning an account list consisting of
// a single account that the plain kestore wallet contains.
func (w *keystoreWallet) Accounts() []accounts.Account {
	return []accounts.Account{w.account}
}

// Contains implements accounts.Wallet, returning whether a particular account is
// or is not wrapped by this wallet instance.
func (w *keystoreWallet) Contains(account accounts.Account) bool {
	return account.Address == w.account.Address && (account.URL == (accounts.URL{}) || account.URL == w.account.URL)
}

// Derive implements accounts.Wallet, but is a noop for plain wallets since there
// is no notion of hierarchical account derivation for plain keystore accounts.
func (w *keystoreWallet) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	return accounts.Account{}, accounts.ErrNotSupported
}

// SelfDerive implements accounts.Wallet, but is a noop for plain wallets since
// there is no notion of hierarchical account derivation for plain keystore accounts.
func (w *keystoreWallet) SelfDerive(base accounts.DerivationPath, chain ethereum.ChainStateReader) {}


//TODO zero delte Sign
/*// SignHash implements accounts.Wallet, attempting to sign the given hash with
// the given account. If the wallet does not wrap this particular account, an
// error is returned to avoid account leakage (even though in theory we may be
// able to sign via our shared keystore backend).
func (w *keystoreWallet) SignHash(account accounts.Account, hash []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if account.Address != w.account.Address {
		return nil, accounts.ErrUnknownAccount
	}
	if account.URL != (accounts.URL{}) && account.URL != w.account.URL {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignHash(account, hash)
}

// SignTx implements accounts.Wallet, attempting to sign the given transaction
// with the given account. If the wallet does not wrap this particular account,
// an error is returned to avoid account leakage (even though in theory we may
// be able to sign via our shared keystore backend).
func (w *keystoreWallet) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Make sure the requested account is contained within
	if account.Address != w.account.Address {
		return nil, accounts.ErrUnknownAccount
	}
	if account.URL != (accounts.URL{}) && account.URL != w.account.URL {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignTx(account, tx, chainID)
}

// SignHashWithPassphrase implements accounts.Wallet, attempting to sign the
// given hash with the given account using passphrase as extra authentication.
func (w *keystoreWallet) SignHashWithPassphrase(account accounts.Account, passphrase string, hash []byte) ([]byte, error) {
	// Make sure the requested account is contained within
	if account.Address != w.account.Address {
		return nil, accounts.ErrUnknownAccount
	}
	if account.URL != (accounts.URL{}) && account.URL != w.account.URL {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignHashWithPassphrase(account, passphrase, hash)
}

// SignTxWithPassphrase implements accounts.Wallet, attempting to sign the given
// transaction with the given account using passphrase as extra authentication.
func (w *keystoreWallet) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	// Make sure the requested account is contained within
	if account.Address != w.account.Address {
		return nil, accounts.ErrUnknownAccount
	}
	if account.URL != (accounts.URL{}) && account.URL != w.account.URL {
		return nil, accounts.ErrUnknownAccount
	}
	// Account seems valid, request the keystore to sign
	return w.keystore.SignTxWithPassphrase(account, passphrase, tx, chainID)
}
*/

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func (w *keystoreWallet) EncryptTx(account accounts.Account, tx *types.Transaction, txt *tx.T,state *state.StateDB) (*types.Transaction, error) {
	// Make sure the requested account is contained within
	if account.Address != w.account.Address {
		return nil, accounts.ErrUnknownAccount
	}
	if account.URL != (accounts.URL{}) && account.URL != w.account.URL {
		return nil, accounts.ErrUnknownAccount
	}
	seed ,err:=w.keystore.GetSeed(account)
	if err != nil{
		return nil, err
	}
	return w.EncryptTxWithSeed(*seed,tx,txt,state)

}

func (w *keystoreWallet) EncryptTxWithSeed(seed common.Seed, btx *types.Transaction, txt *tx.T,state *state.StateDB) (*types.Transaction, error) {

	for i,ctx := range txt.CTxs{
	   tk:=keys.Seed2Tk(seed.SeedToUint256())
		outs,ammount,err := txs.GetRoots(&tk,state.GetZState(),ctx.Cost().ToRef(),&ctx.Currency)
		if err != nil{
			return nil, err
		}
		ins :=[]tx.In{}
		for _,out := range outs{
			ins = append(ins,tx.In{out})
		}
		txt.CTxs[i].Ins =ins

		balance:=ammount
		balance.SubU(ctx.Cost().ToRef())

		if balance.Cmp(&utils.U256_0) >0 {
			selfOut :=tx.Out{
				Addr:keys.Seed2Addr(seed.SeedToUint256()),
				Value:balance,
				Z:tx.TYPE_Z,
			}
			txt.CTxs[i].Outs = append(txt.CTxs[i].Outs,selfOut)
		}
	}

	Ehash :=rlpHash([] interface{}{
		btx.GasPrice(),
		btx.Data(),
		btx.Currency(),
	})
	copy(txt.Ehash[:],Ehash[:])

	stx,err := txs.Gen(seed.SeedToUint256(),txt,state.GetZState())
	if err !=nil {
		return nil,err
	}

	return btx.WithEncrypt(&stx)

}



func (w *keystoreWallet) EncryptTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction,txt *tx.T,state *state.StateDB) (*types.Transaction, error) {
	// Make sure the requested account is contained within
	if account.Address != w.account.Address {
		return nil, accounts.ErrUnknownAccount
	}
	if account.URL != (accounts.URL{}) && account.URL != w.account.URL {
		return nil, accounts.ErrUnknownAccount
	}

	seed ,err:=w.keystore.GetSeedWithPassphrase(account, passphrase)
	if err != nil{
		return nil, err
	}
	return w.EncryptTxWithSeed(*seed,tx,txt,state)

}


/*func (w *keystoreWallet) GetSeed() (*common.Seed, error) {
	// Make sure the requested account is contained within
	seed ,err:=w.keystore.GetSeed(w.account)
	if err != nil{
		return nil, err
	}
	return seed,nil

}*/

func (w *keystoreWallet) IsMine(onceAddress common.Address) (bool) {
	tk:=w.account.Tk.ToUint512()
	return keys.IsMyPKr(tk,onceAddress.ToUint512())

}







