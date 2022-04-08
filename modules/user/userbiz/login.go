package userbiz

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/component/tokenprovider"
	"BiMarketplace/modules/user/usermodel"
	"context"
	"math/rand"
	"strconv"
)

type LoginStore interface {
	FindUser(ctx context.Context, Address string) (*usermodel.User, error)
	GetNonceDB(address string) (*usermodel.Nonce, error)
	ChangeNonce(nonce int, userData *usermodel.User) error
}

type LoginBiz struct {
	store         LoginStore
	tokenProvider tokenprovider.Provider
	expiry        int
}

func NewLoginBiz(store LoginStore, tokenProvider tokenprovider.Provider, expiry int) *LoginBiz {
	return &LoginBiz{
		store:         store,
		tokenProvider: tokenProvider,
		expiry:        expiry,
	}
}

func (biz *LoginBiz) Login(ctx context.Context, loginData usermodel.UserLogin) (*tokenprovider.Token, error) {
	nonce, err := biz.store.GetNonceDB(loginData.Address)
	if err != nil {
		if err == appCommon.RecordNotFound {
			return nil, appCommon.ErrEntityNotFound(usermodel.NonceEntityName, err)
		}
		return nil, appCommon.ErrCannotGetEntity(usermodel.NonceEntityName, err)
	}
	nonceValue := nonce.Nonce
	if !appCommon.VerifySig(loginData.Address, loginData.Signature, []byte(strconv.Itoa(nonceValue))) {
		return nil, appCommon.ErrInvalidSignature()
	}

	user, err := biz.store.FindUser(ctx, loginData.Address)
	if err != nil {
		if err == appCommon.RecordNotFound {
			return nil, appCommon.ErrEntityNotFound(usermodel.EntityName, err)
		}
		return nil, appCommon.ErrCannotGetEntity(usermodel.EntityName, err)
	}

	payload := tokenprovider.TokenPayload{
		UserId:  user.Id,
		Address: user.WalletAddress,
	}

	accessToken, err := biz.tokenProvider.Generate(payload, biz.expiry)
	if err != nil {
		return nil, appCommon.ErrInternal(err)
	}

	if err := biz.store.ChangeNonce(rand.Intn(1847483647-1000000000+1)+1000000000, user); err != nil {
		return nil, appCommon.ErrInternal(err)
	}

	return accessToken, nil
}
