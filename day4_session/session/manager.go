package session

import (
	__template_and_file "geektime-go/day3_template_and_file"
	"github.com/google/uuid"
)

// 语法糖|胶水，非核心能力，帮助用户便捷操作
type Manager struct {
	Store
	Propagator
	SessionKey string
}

func (m *Manager) GetSession(c *__template_and_file.Context) (Session, error) {
	id, err := m.Extract(c.R)
	if err != nil {
		return nil, err
	}

	return m.Get(c.R.Context(), id)
}

func (m *Manager) InitSession(c *__template_and_file.Context) (Session, error) {
	id := uuid.New().String()
	sess, err := m.Generator(c.R.Context(), id)
	if err != nil {
		return nil, err
	}
	err = m.Inject(sess.ID(), c.W)
	if err != nil {
		return nil, err
	}
	return sess, nil
}

func (m *Manager) RemoveSession(c *__template_and_file.Context) error {
	sess, err := m.GetSession(c)
	if err != nil {
		return err
	}
	err = m.Store.Remove(c.R.Context(), sess.ID())
	if err != nil {
		return err
	}
	err = m.Propagator.Remove(c.W)
	if err != nil {
		return err
	}
	return nil
}
