package pools

import (
	"sync"

	"github.com/DimTur/chat-websocket-go/internal/domain"
)

var Users = userPool{
	pool: make(map[domain.ID]chan interface{}),
}

type userPool struct {
	sync.Mutex
	// key - user id
	pool map[domain.ID]chan interface{}
}

func (p *userPool) Send(userID domain.ID, msg interface{}) {
	p.Lock()
	defer p.Unlock()

	ch, ok := p.pool[userID]
	if !ok {
		return
	}

	// на подумать: буферизированный или горутину
	ch <- msg
}

func (p *userPool) New(userID domain.ID) <-chan interface{} {
	p.Lock()
	ch := make(chan interface{})
	p.pool[userID] = ch
	p.Unlock()

	return ch
}

func (p *userPool) Delete(userID domain.ID) bool {
	p.Lock()
	defer p.Unlock()

	ch, ok := p.pool[userID]
	if !ok {
		return ok
	}

	delete(p.pool, userID)
	close(ch)
	return ok
}
