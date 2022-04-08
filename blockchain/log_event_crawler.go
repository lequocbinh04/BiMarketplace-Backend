package blockchain

import (
	"BiMarketplace/component/asyncjob"
	"context"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	log "github.com/sirupsen/logrus"
	"math/big"
	"strconv"
	"strings"
	"time"
)

const (
	BlockStart    int64 = 16061292
	DefaultRPCURL       = "https://data-seed-prebsc-1-s1.binance.org:8545"
	BiEggAddress        = "0x358A141Af732a7025d5dC3229513B5FAF20eD78d"
)

type BlockTrackerStorage interface {
	GetLatestBlockNumber(ctx context.Context) (int64, error)
	UpdateLatestBlockNumber(ctx context.Context, num int64) error
}
type logCrawler struct {
	rpcURL              string
	currentBlock        int64
	latestBlock         int64
	eggAddress          string
	contractABIFilePath string
	client              *ethclient.Client
	configStorage       BlockTrackerStorage
	logChan             chan types.Log
}

func NewLogCrawler(
	rpcURL string,
	blockStart int64,
	eggAddress string,
	store BlockTrackerStorage,
) *logCrawler {
	newCrawler := &logCrawler{
		rpcURL:        stringOrDefault(rpcURL, DefaultRPCURL),
		currentBlock:  max(blockStart, BlockStart),
		eggAddress:    stringOrDefault(eggAddress, BiEggAddress),
		client:        nil,
		configStorage: store,
		logChan:       make(chan types.Log, 100),
	}

	return newCrawler
}
func (crawler *logCrawler) GetLogChan() chan types.Log {
	return crawler.logChan
}

func (crawler *logCrawler) Start() error {
	client, err := ethclient.Dial(crawler.rpcURL)
	if err != nil {
		return err
	}

	crawler.client = client
	latestBlockNumber, err := crawler.latestBlockNumber()
	if err != nil {
		return err
	}
	latestDBBlockNumber, err := crawler.latestDBBlockNumber()
	if err != nil {
		return err
	}
	crawler.currentBlock = max(latestDBBlockNumber, crawler.currentBlock)
	crawler.latestBlock = latestBlockNumber

	go func() {
		var stepBlockFastScan int64 = 200
		for {
			time.Sleep(time.Second * 2)
			if err := crawler.configStorage.UpdateLatestBlockNumber(context.Background(), crawler.currentBlock); err != nil {
				continue
			}
			if latestBlockNumber > crawler.currentBlock {
				if v := latestBlockNumber - crawler.currentBlock; v < stepBlockFastScan {
					stepBlockFastScan = v
				}
				if err := crawler.scanBlock(crawler.currentBlock, crawler.currentBlock+stepBlockFastScan); err != nil {
					continue
				}
				crawler.currentBlock += stepBlockFastScan + 1
				continue
			}
			latestBlockNumber, err = crawler.latestBlockNumber()
			if err != nil {
				continue
			}
			crawler.latestBlock = latestBlockNumber
			// no new block
			if latestBlockNumber <= crawler.currentBlock {
				continue
			}

			if err := crawler.scanBlock(crawler.currentBlock, crawler.currentBlock); err != nil {
				log.Errorln(err)
				continue
			}
			crawler.currentBlock += 1
		}
	}()

	return nil
}

func (crawler *logCrawler) latestBlockNumber() (int64, error) {
	var result int64 = 1
	job := asyncjob.NewJob(func(ctx context.Context) error {
		latestBlock, err := crawler.client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			return err
		}

		result, err = strconv.ParseInt(latestBlock.Number.String(), 10, 64)
		if err != nil {
			return err
		}
		return nil
	})
	job.SetRetryDurations(time.Second, time.Second*2, time.Second*3)
	if err := asyncjob.NewGroup(false, job).Run(context.Background()); err != nil {
		return 0, err
	}
	return result, nil
}

func (crawler *logCrawler) latestDBBlockNumber() (int64, error) {
	var result int64 = 1
	job := asyncjob.NewJob(func(ctx context.Context) error {
		var err error
		result, err = crawler.configStorage.GetLatestBlockNumber(ctx)
		if err != nil {
			return err
		}
		return nil
	})
	job.SetRetryDurations(time.Second, time.Second*2, time.Second*3)
	if err := asyncjob.NewGroup(false, job).Run(context.Background()); err != nil {
		return 0, err
	}
	return result, nil
}

func (crawler *logCrawler) scanBlock(from, to int64) error {

	logs, err := crawler.client.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: big.NewInt(from),
		ToBlock:   big.NewInt(to),
		Addresses: []common.Address{common.HexToAddress(crawler.eggAddress)},
	})

	if err != nil {
		return err
	}
	for i, l := range logs {
		functionHash := strings.ToLower(l.Topics[0].Hex())
		log.Printf("Block %d - Tx %s - Event %s \n", l.BlockNumber, l.TxHash.Hex(), functionHash)

		crawler.logChan <- logs[i]
	}
	return nil
}

func stringOrDefault(s, d string) string {
	if s == "" {
		return d
	}

	return s
}

func max(n1, n2 int64) int64 {
	if n1 > n2 {
		return n1
	}

	return n2
}
