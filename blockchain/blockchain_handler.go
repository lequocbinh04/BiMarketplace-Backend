package blockchain

import (
	"BiMarketplace/appCommon"
	"BiMarketplace/component/asyncjob"
	"BiMarketplace/modules/egg/eggmodel"
	"BiMarketplace/modules/pet/petmodel"
	"BiMarketplace/modules/transaction/transactionmodel"
	"BiMarketplace/modules/user/usermodel"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"math/big"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	EvtBuyEgg  = "0xe616638ab37f5114159a014ebd8dfe16409b3f2621e10a38a77c3f232393c627"
	EvtOpenEgg = "0x59284618120d3344285aec7c8efc349ee04ed47daf976da3f715e2c27a415fe0"
)

var allPetName = []string{
	"Keg",
	"Patch",
	"Ebony",
	"Dorito",
	"Lemming",
	"Wookie",
	"Cadburry",
	"Alfred",
	"Barry",
	"Cecil",
	"Dennis",
	"Edgar",
	"Fernando",
	"Gustav",
	"Hans",
	"Igor",
	"Jens",
	"Karl",
	"Lars",
	"Mikkel",
	"Nils",
	"Ole",
	"Palle",
	"Quentin",
	"Rasmus",
	"Sebastian",
	"Thomas",
	"Ulrich",
	"Viktor",
	"William",
	"Xavier",
	"Yannick",
	"Zachary",
}

var allPetImage = []string{
	"Pets_1.png",
	"Pets_2.png",
	"Pets_3.png",
	"Pets_4.png",
	"Pets_5.png",
	"Pets_6.png",
	"Pets_7.png",
	"Pets_8.png",
	"Pets_9.png",
	"Pets_10.png",
	"Pets_11.png",
	"Pets_12.png",
}

type EggStorage interface {
	Create(data *eggmodel.Egg) error
	GetDataWithCondition(ctx context.Context, condition map[string]interface{}, moreKeys ...string) (*eggmodel.Egg, error)
	DecreaseEgg(id int) error
}

type PetStorage interface {
	CreateNewPet(data *petmodel.Pet) error
	GetDataWithCondition(ctx context.Context, condition map[string]interface{}, moreKeys ...string) (*petmodel.Pet, error)
}

type TransactionStorage interface {
	CreateNewTx(data *transactionmodel.Transaction) error
	GetDataWithCondition(ctx context.Context, condition map[string]interface{}, moreKeys ...string) (*transactionmodel.Transaction, error)
}

type UserStorage interface {
	GetUserByAddress(address string) (*usermodel.User, error)
}

type mkpHdl struct {
	eggStore  EggStorage
	petStore  PetStorage
	txStore   TransactionStorage
	userStore UserStorage
	eggAbiObj abi.ABI
}

func NewMkpHdl(eggStore EggStorage, petStore PetStorage, txStore TransactionStorage, userStore UserStorage, eggAbi string) *mkpHdl {
	f, err := os.Open(eggAbi)
	eggABI, err := abi.JSON(f)
	if err != nil {
		log.Fatalln(err)
	}

	return &mkpHdl{
		eggStore:  eggStore,
		petStore:  petStore,
		txStore:   txStore,
		userStore: userStore,
		eggAbiObj: eggABI,
	}
}

func (p *mkpHdl) Run(ctx context.Context, queue <-chan types.Log) {
	for l := range queue {
		functionHash := strings.ToLower(l.Topics[0].Hex())
		var job asyncjob.Job

		switch functionHash {
		case EvtBuyEgg:
			job = asyncjob.NewJob(func(ctx context.Context) error {
				return p.handleBuyEgg(ctx, l)
			})
		case EvtOpenEgg:
			job = asyncjob.NewJob(func(ctx context.Context) error {
				return p.handleOpenEgg(ctx, l)
			})
		default:
			log.Printf("Block %d - Tx %s - Event %s \n", l.BlockNumber, l.TxHash.Hex(), functionHash)
			continue
		}

		job.SetRetryDurations(time.Second, time.Second, time.Second, time.Second) // 4 times (1s each)

		if err := job.Execute(ctx); err != nil {
			log.Errorln(err)
		}
	}
}

func (p *mkpHdl) handleOpenEgg(ctx context.Context, l types.Log) error {
	tx, err := p.txStore.GetDataWithCondition(ctx, map[string]interface{}{
		"tx_hash": l.TxHash.Hex(),
	})
	if tx != nil {
		return nil
	}
	if err != appCommon.RecordNotFound {
		return err
	}
	receiver := strings.ToLower(common.HexToAddress(l.Topics[1].Hex()).Hex())
	nftId, _ := strconv.ParseInt(strings.ReplaceAll(l.Topics[2].Hex(), "0x", ""), 16, 64)

	user, err := p.userStore.GetUserByAddress(receiver)
	if err != nil {
		return err
	}

	newTx := transactionmodel.Transaction{
		TxHash:       l.TxHash.Hex(),
		PetId:        nftId,
		OwnerId:      int64(user.Id),
		BlockNumber:  int64(l.BlockNumber),
		BuyerAddress: receiver,
	}
	newTx.Status = "open_egg"

	if err := p.txStore.CreateNewTx(&newTx); err != nil {
		return err
	}

	newPet := petmodel.Pet{
		OwnerId: int64(user.Id),
		Hp:      100,
		Attack:  100,
		Speed:   100,
		NftId:   nftId,
		Name:    allPetName[rand.Intn(len(allPetName))],
		Image:   allPetImage[rand.Intn(len(allPetImage))],
	}

	if err := p.petStore.CreateNewPet(&newPet); err != nil {
		return err
	}

	egg, err := p.eggStore.GetDataWithCondition(ctx, map[string]interface{}{
		"owner_address": receiver,
		"status":        "not_open",
	})

	if err != nil {
		return err
	}

	if err := p.eggStore.DecreaseEgg(egg.Id); err != nil {
		return err
	}

	return nil
}

func (p *mkpHdl) handleBuyEgg(ctx context.Context, l types.Log) error {
	tx, err := p.eggStore.GetDataWithCondition(ctx, map[string]interface{}{
		"tx_hash": l.TxHash.Hex(),
	})
	if tx != nil {
		return nil
	}
	if err != appCommon.RecordNotFound {
		return err
	}
	event := struct {
		Receiver     string
		PaymentToken string
		Price        *big.Int
	}{}
	event.Receiver = strings.ToLower(common.HexToAddress(l.Topics[1].Hex()).Hex())

	ii := new(big.Int)
	ii.SetString(strings.ReplaceAll(l.Topics[3].Hex(), "0x", ""), 16)
	idiv := new(big.Int)
	idiv.SetString("1000000000000000000", 10)
	ii = ii.Div(ii, idiv)

	event.Price = ii
	event.PaymentToken = strings.ToLower(common.HexToAddress(l.Topics[2].Hex()).Hex())
	price := decimal.NewNullDecimal(decimal.NewFromBigInt(event.Price, 0))

	newEgg := eggmodel.Egg{
		TxHash:       l.TxHash.Hex(),
		OwnerAddress: event.Receiver,
		PayableToken: event.PaymentToken,
		Price:        price,
		BlockNumber:  int64(l.BlockNumber),
	}
	newEgg.Status = "not_open"
	if err := p.eggStore.Create(&newEgg); err != nil {
		return err
	}
	return nil
}
