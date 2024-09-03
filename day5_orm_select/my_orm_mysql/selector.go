package my_orm_mysql

import "context"

type selector[T any] struct {
}

func (s *selector[T]) Build() (*Query, error) {
	//TODO implement me
	panic("implement me")
}

func (s *selector[T]) Get(ctx context.Context) (T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
