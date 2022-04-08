package userbiz

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/modules/user/usermodel"
	"math/rand"
	"regexp"
)

type GetNonceStore interface {
	GetNonceDB(address string) (*usermodel.Nonce, error)
	CreateNonce(nonce int, address string) error
}

type getNonceBiz struct {
	store GetNonceStore
}

func NewGetNonceBiz(store GetNonceStore) *getNonceBiz {
	return &getNonceBiz{store: store}
}

func (biz *getNonceBiz) GetNonce(address string) (*usermodel.Nonce, error) {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")

	if !re.MatchString(address) {
		return nil, appCommon.ErrInvalidAddress()
	}

	nonce, err := biz.store.GetNonceDB(address)
	if err != nil {
		if err == appCommon.RecordNotFound {
			nonce = &usermodel.Nonce{
				Address: address,
				Nonce:   rand.Intn(1847483647-1000000000+1) + 1000000000,
			}
			err = biz.store.CreateNonce(nonce.Nonce, address)
			if err != nil {
				return nil, appCommon.ErrCannotCreateEntity(usermodel.NonceEntityName, err)
			}
		} else {
			return nil, appCommon.ErrCannotGetEntity(usermodel.NonceEntityName, err)
		}
	}
	return nonce, nil
}
