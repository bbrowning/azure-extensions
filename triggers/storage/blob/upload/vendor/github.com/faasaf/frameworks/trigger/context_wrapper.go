package trigger

import "github.com/faasaf/frameworks/common"

type ContextWrapper interface {
	GetContext() common.Context
	ResC() chan common.Context
	ErrC() chan error
}

type contextWrapper struct {
	ctx   common.Context
	resCh chan common.Context
	errCh chan error
}

func NewContextWrapper(ctx common.Context) ContextWrapper {
	return &contextWrapper{
		ctx:   ctx,
		resCh: make(chan common.Context),
		errCh: make(chan error),
	}
}

func (c *contextWrapper) GetContext() common.Context {
	return c.ctx
}

func (c *contextWrapper) ResC() chan common.Context {
	return c.resCh
}

func (c *contextWrapper) ErrC() chan error {
	return c.errCh
}
