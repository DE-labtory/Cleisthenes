package bba

import (
	"encoding/json"
	"errors"
	"log"
	"sync"
	"sync/atomic"

	"github.com/golang/protobuf/ptypes"

	"github.com/DE-labtory/cleisthenes/pb"

	"github.com/DE-labtory/cleisthenes"
)

var ErrUndefinedRequestType = errors.New("unexpected request type")

type request struct {
	sender cleisthenes.Member
	data   *pb.Message_Bba
	err    chan error
}

type BBA struct {
	*sync.RWMutex

	// number of network nodes
	n int

	// number of byzantine nodes which can tolerate
	f int

	proposer cleisthenes.Member

	stopFlag int32
	// done is flag whether BBA is terminated or not
	done bool

	// epoch is current honeybadger epoch
	epoch uint64
	// round is value for current Binary Byzantine Agreement round
	round       uint64
	binValueSet *binarySet

	// broadcastedBvalSet is set of bval value instance has sent
	broadcastedBvalSet *binarySet

	// est is estimated value of BBA instance, dec is decision value
	est, dec *cleisthenes.BinaryState

	bvalRepo        cleisthenes.RequestRepository
	auxRepo         cleisthenes.RequestRepository
	incomingReqRepo cleisthenes.IncomingRequestRepository

	reqChan             chan request
	closeChan           chan struct{}
	binValueChan        chan struct{}
	tryoutAgreementChan chan struct{}
	advanceRoundChan    chan struct{}

	broadcaster cleisthenes.Broadcaster

	coinGenerator cleisthenes.CoinGenerator
}

func New(n int) *BBA {
	return &BBA{
		n:                   n,
		est:                 cleisthenes.NewBinaryState(),
		dec:                 cleisthenes.NewBinaryState(),
		reqChan:             make(chan request),
		closeChan:           make(chan struct{}),
		binValueChan:        make(chan struct{}),
		tryoutAgreementChan: make(chan struct{}),
		advanceRoundChan:    make(chan struct{}),
	}
}

// HandleInput will set the given val as the initial value to be proposed in the
// Agreement
func (bba *BBA) HandleInput(sender cleisthenes.Member, msg *pb.Message_Bba) error {
	return nil
}

// HandleMessage will process the given rpc message.
func (b *BBA) HandleMessage(sender cleisthenes.Member, msg *pb.Message_Bba) error {
	return nil
}

func (bba *BBA) Result() bool {
	return false
}

func (bba *BBA) Close() {
	if first := atomic.CompareAndSwapInt32(&bba.stopFlag, int32(0), int32(1)); !first {
		return
	}
	bba.closeChan <- struct{}{}
	<-bba.closeChan
}

func (bba *BBA) muxMessage(sender cleisthenes.Member, msg *pb.Message_Bba) error {
	// TODO: request epoch check
	switch msg.Bba.Type {
	case pb.BBA_BVAL:
		bval := &bvalRequest{}
		if err := json.Unmarshal(msg.Bba.Payload, bval); err != nil {
			return err
		}
		return bba.handleBvalRequest(sender, bval)
	case pb.BBA_AUX:
		aux := &auxRequest{}
		if err := json.Unmarshal(msg.Bba.Payload, aux); err != nil {
			return err
		}
		return bba.handleAuxRequest(sender, aux)
	default:
		return ErrUndefinedRequestType
	}
}

func (bba *BBA) handleBvalRequest(sender cleisthenes.Member, bval *bvalRequest) error {
	if err := bba.saveBvalIfNotExist(sender, bval); err != nil {
		return err
	}
	count := bba.countBvalByValue(bval.Value)
	if count == bba.binValueSetThreshold() {
		bba.binValueSet.union(bval.Value)
		bba.binValueChan <- struct{}{}
		return nil
	}
	if count == bba.bvalBroadcastThreshold() && !bba.broadcastedBvalSet.exist(bval.Value) {
		bba.broadcastedBvalSet.union(bval.Value)
		return bba.broadcast(pb.BBA_BVAL, bval)
	}
	return nil
}

func (bba *BBA) handleAuxRequest(sender cleisthenes.Member, aux *auxRequest) error {
	if err := bba.saveAuxIfNotExist(sender, aux); err != nil {
		return err
	}
	count := bba.countAuxByValue(aux.Value)
	if count < bba.tryoutAgreementThreshold() {
		return nil
	}
	bba.tryoutAgreementChan <- struct{}{}
	return nil
}

func (bba *BBA) tryoutAgreement() {
	if bba.done {
		return
	}
	coin := bba.coinGenerator.Coin()

	binList := bba.binValueSet.toList()
	if len(binList) == 0 {
		log.Fatalf("binary set is empty, but tried agreement")
		return
	}
	if len(binList) > 1 {
		bba.est.Set(cleisthenes.Binary(coin))
		bba.advanceRoundChan <- struct{}{}
		return
	}

	if binList[0] != cleisthenes.Binary(coin) {
		bba.est.Set(binList[0])
		bba.advanceRoundChan <- struct{}{}
		return
	}
	bba.dec.Set(binList[0])
	bba.done = true
	bba.advanceRoundChan <- struct{}{}
}

func (bba *BBA) broadcast(typ pb.BBAType, req cleisthenes.Request) error {
	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}
	bba.broadcaster.ShareMessage(pb.Message{
		Timestamp: ptypes.TimestampNow(),
		Payload: &pb.Message_Bba{
			Bba: &pb.BBA{
				Round:   bba.round,
				Type:    pb.BBA_BVAL,
				Payload: payload,
			},
		},
	})
	return nil
}

func (bba *BBA) toDie() bool {
	return atomic.LoadInt32(&(bba.stopFlag)) == int32(1)
}

func (bba *BBA) run() {
	for !bba.toDie() {
		select {
		case <-bba.closeChan:
			bba.closeChan <- struct{}{}
		case req := <-bba.reqChan:
			req.err <- bba.muxMessage(req.sender, req.data)
		case <-bba.binValueChan:
			// TODO: broadcast AUX message
		case <-bba.tryoutAgreementChan:
			bba.tryoutAgreement()
		case <-bba.advanceRoundChan:
			// TODO: advance round and initialize new round
		}
	}
}

func (bba *BBA) saveBvalIfNotExist(sender cleisthenes.Member, data *bvalRequest) error {
	r, err := bba.bvalRepo.Find(sender.Address)
	if err != nil {
		return err
	}
	if r != nil {
		return nil
	}
	return bba.bvalRepo.Save(sender.Address, data)
}

func (bba *BBA) saveAuxIfNotExist(sender cleisthenes.Member, data *auxRequest) error {
	r, err := bba.auxRepo.Find(sender.Address)
	if err != nil {
		return err
	}
	if r != nil {
		return nil
	}
	return bba.auxRepo.Save(sender.Address, data)
}

func (bba *BBA) countBvalByValue(val cleisthenes.Binary) int {
	bvalList := bba.convToBvalList(bba.bvalRepo.FindAll())
	count := 0
	for _, bval := range bvalList {
		if bval.Value == val {
			count++
		}
	}
	return count
}

func (bba *BBA) countAuxByValue(val cleisthenes.Binary) int {
	auxList := bba.convToAuxList(bba.auxRepo.FindAll())
	count := 0
	for _, aux := range auxList {
		if aux.Value == val {
			count++
		}
	}
	return count
}

func (bba *BBA) convToBvalList(reqList []cleisthenes.Request) []*bvalRequest {
	result := make([]*bvalRequest, 0)
	for _, req := range reqList {
		result = append(result, req.(*bvalRequest))
	}
	return result
}

func (bba *BBA) convToAuxList(reqList []cleisthenes.Request) []*auxRequest {
	result := make([]*auxRequest, 0)
	for _, req := range reqList {
		result = append(result, req.(*auxRequest))
	}
	return result
}

func (bba *BBA) bvalBroadcastThreshold() int {
	return bba.f + 1
}

func (bba *BBA) binValueSetThreshold() int {
	return 2*bba.f + 1
}

func (bba *BBA) tryoutAgreementThreshold() int {
	return bba.n - bba.f
}
