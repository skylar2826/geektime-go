package day4_session

import (
	__template_and_file "geektime-go/day3_template_and_file"
	"github.com/google/uuid"
)

type manager struct {
	Propagator
	Store
	CtxSessionKey string
}

func (m *manager) GetSession(c *__template_and_file.Context) (Session, error) {
	/*
		频繁读取redis中的session => 尝试缓存住数据 => 只能缓存在context中
	*/
	if c.UserValues == nil {
		c.UserValues = make(map[string]any, 1)
	}

	session, ok := c.UserValues[m.CtxSessionKey]
	if ok {
		return session.(Session), nil
	}

	sessionId, err := m.Extract(c.R)
	if err != nil {
		return nil, err
	}

	session, err = m.Get(c.R.Context(), sessionId)
	if err != nil {
		return nil, err
	}
	c.UserValues[m.CtxSessionKey] = session
	return session.(Session), err
}

func (m *manager) InitSession(c *__template_and_file.Context) (Session, error) {
	sessionId := uuid.New().String()
	session, err := m.Generator(c.R.Context(), sessionId)
	if err != nil {
		return nil, err
	}
	err = m.Inject(session.ID(), c.W)
	return session, err
}

func (m *manager) RemoveSession(c *__template_and_file.Context) error {
	//sessionId, err := m.Extract(c.R)
	session, err := m.GetSession(c)
	if err != nil {
		return err
	}
	err = m.Store.Remove(c.R.Context(), session.ID())
	if err != nil {
		return err
	}
	return m.Propagator.Remove(c.W)
}

func (m *manager) RefreshSession(c *__template_and_file.Context) error {
	session, err := m.GetSession(c)
	if err != nil {
		return err
	}
	// 刷新假设sessionId不变
	return m.Refresh(c.R.Context(), session.ID())
}
