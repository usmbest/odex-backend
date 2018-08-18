package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	validation "github.com/go-ozzo/ozzo-validation"
	"gopkg.in/mgo.v2/bson"
)

// Pair struct is used to model the pair data in the system and DB
type Pair struct {
	ID                bson.ObjectId  `json:"id" bson:"_id"`
	Name              string         `json:"name" bson:"name"`
	BaseTokenID       bson.ObjectId  `json:"baseTokenId" bson:"baseTokenId"`
	BaseTokenSymbol   string         `json:"baseTokenSymbol" bson:"baseTokenSymbol"`
	BaseTokenAddress  common.Address `json:"baseTokenAddress" bson:"baseTokenAddress"`
	QuoteTokenID      bson.ObjectId  `json:"quoteTokenId" bson:"quoteTokenId"`
	QuoteTokenSymbol  string         `json:"quoteTokenSymbol" bson:"quoteTokenSymbol"`
	QuoteTokenAddress common.Address `json:"quoteTokenAddress" bson:"quoteTokenAddress"`

	Active  bool     `json:"active" bson:"active"`
	MakeFee *big.Int `json:"makeFee" bson:"makeFee"`
	TakeFee *big.Int `json:"takeFee" bson:"takeFee"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type PairSubDoc struct {
	Name       string         `json:"name" bson:"name"`
	BaseToken  common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken common.Address `json:"quoteToken" bson:"quoteToken"`
}

type PairRecord struct {
	ID                bson.ObjectId `json:"id" bson:"_id"`
	Name              string        `json:"name" bson:"name"`
	BaseTokenID       bson.ObjectId `json:"baseTokenId" bson:"baseTokenId"`
	BaseTokenSymbol   string        `json:"baseTokenSymbol" bson:"baseTokenSymbol"`
	BaseTokenAddress  string        `json:"baseTokenAddress" bson:"baseTokenAddress"`
	QuoteTokenID      bson.ObjectId `json:"quoteTokenId" bson:"quoteTokenId"`
	QuoteTokenSymbol  string        `json:"quoteTokenSymbol" bson:"quoteTokenSymbol"`
	QuoteTokenAddress string        `json:"quoteTokenAddress" bson:"quoteTokenAddress"`

	Active  bool   `json:"active" bson:"active"`
	MakeFee string `json:"makeFee" bson:"makeFee"`
	TakeFee string `json:"takeFee" bson:"takeFee"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (p *Pair) SetBSON(raw bson.Raw) error {
	decoded := &PairRecord{}

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	makeFee := new(big.Int)
	makeFee, _ = makeFee.SetString(decoded.MakeFee, 10)
	takeFee := new(big.Int)
	takeFee, _ = takeFee.SetString(decoded.TakeFee, 10)

	p.ID = decoded.ID
	p.Name = decoded.Name
	p.BaseTokenID = decoded.BaseTokenID
	p.BaseTokenSymbol = decoded.BaseTokenSymbol
	p.BaseTokenAddress = common.HexToAddress(decoded.BaseTokenAddress)
	p.QuoteTokenID = decoded.QuoteTokenID
	p.QuoteTokenSymbol = decoded.QuoteTokenSymbol
	p.QuoteTokenAddress = common.HexToAddress(decoded.QuoteTokenAddress)

	p.Active = decoded.Active
	p.MakeFee = makeFee
	p.TakeFee = takeFee

	p.CreatedAt = decoded.CreatedAt
	p.UpdatedAt = decoded.UpdatedAt

	return nil
}

func (p *Pair) GetBSON() (interface{}, error) {
	return &PairRecord{
		ID:                p.ID,
		Name:              p.Name,
		BaseTokenID:       p.BaseTokenID,
		BaseTokenSymbol:   p.BaseTokenSymbol,
		BaseTokenAddress:  p.BaseTokenAddress.Hex(),
		QuoteTokenID:      p.QuoteTokenID,
		QuoteTokenSymbol:  p.QuoteTokenSymbol,
		QuoteTokenAddress: p.QuoteTokenAddress.Hex(),
		Active:            p.Active,
		MakeFee:           p.MakeFee.String(),
		TakeFee:           p.TakeFee.String(),
		CreatedAt:         p.CreatedAt,
		UpdatedAt:         p.UpdatedAt,
	}, nil
}

// Validate function is used to verify if an instance of
// struct satisfies all the conditions for a valid instance
func (p Pair) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.BaseTokenAddress, validation.Required),
		validation.Field(&p.QuoteTokenAddress, validation.Required),
	)
}

// GetOrderBookKeys returns the orderbook price point keys for corresponding pair
// It is used to fetch the orderbook from redis of a pair
func (p *Pair) GetOrderBookKeys() (sell, buy string) {
	return p.BaseTokenAddress.Hex() + "::" + p.QuoteTokenAddress.Hex() + "::sell", p.BaseTokenAddress.Hex() + "::" + p.QuoteTokenAddress.Hex() + "::buy"
}
