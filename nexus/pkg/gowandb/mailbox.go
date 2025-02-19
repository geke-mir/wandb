package gowandb

import (
	"strings"

	"github.com/wandb/wandb/nexus/internal/shared"
	"github.com/wandb/wandb/nexus/pkg/service"
)

type MailboxHandle struct {
	responseChan chan *service.Result
}

type Mailbox struct {
	handles map[string]*MailboxHandle
}

func NewMailbox() *Mailbox {
	mailbox := &Mailbox{}
	mailbox.handles = make(map[string]*MailboxHandle)
	return mailbox
}

func NewMailboxHandle() *MailboxHandle {
	mbh := &MailboxHandle{responseChan: make(chan *service.Result)}
	return mbh
}

func (mbh *MailboxHandle) wait() *service.Result {
	got := <-mbh.responseChan
	return got
}

func (mb *Mailbox) Deliver(rec *service.Record) *MailboxHandle {
	uuid := "nexus:" + shared.ShortID(12)
	rec.Control = &service.Control{MailboxSlot: uuid}
	handle := NewMailboxHandle()
	mb.handles[uuid] = handle
	return handle
}

func (mb *Mailbox) Respond(result *service.Result) bool {
	slot := result.GetControl().MailboxSlot
	if !strings.HasPrefix(slot, "nexus:") {
		return false
	}
	handle, ok := mb.handles[slot]
	if ok {
		handle.responseChan <- result
		// clean up after thyself?
		delete(mb.handles, slot)
	}
	return ok
}
