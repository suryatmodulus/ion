package biz

import (
	"github.com/pion/ion/pkg/grpc/biz"
	"github.com/pion/ion/pkg/grpc/ion"
	"github.com/pion/ion/pkg/util"
)

// Peer represents a peer for client
type Peer struct {
	uid             string
	sid             string
	info            []byte
	lastStreamEvent *ion.StreamEvent
	closed          util.AtomicBool
	sndCh           chan *biz.SignalReply
}

func NewPeer(sid string, uid string, info []byte, senCh chan *biz.SignalReply) *Peer {
	p := &Peer{
		uid:   uid,
		sid:   sid,
		info:  info,
		sndCh: senCh,
	}
	p.closed.Set(false)
	return p
}

// Close peer
func (p *Peer) Close() {
	if p.closed.Get() {
		return
	}
	p.closed.Set(true)
}

// UID return peer uid
func (p *Peer) UID() string {
	return p.uid
}

// SID return session id
func (p *Peer) SID() string {
	return p.sid
}

func (p *Peer) send(data *biz.SignalReply) error {
	go func() {
		p.sndCh <- data
	}()
	return nil
}

func (p *Peer) sendPeerEvent(event *ion.PeerEvent) error {
	data := &biz.SignalReply{
		Payload: &biz.SignalReply_PeerEvent{
			PeerEvent: event,
		},
	}

	return p.send(data)
}

func (p *Peer) sendStreamEvent(event *ion.StreamEvent) error {
	data := &biz.SignalReply{
		Payload: &biz.SignalReply_StreamEvent{
			StreamEvent: event,
		},
	}
	return p.send(data)
}

func (p *Peer) sendMessage(msg *ion.Message) error {
	data := &biz.SignalReply{
		Payload: &biz.SignalReply_Msg{
			Msg: msg,
		},
	}
	return p.send(data)
}
