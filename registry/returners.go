package registry

import (
	"github.com/ovechkin-dm/mockio/v2/matchers"
)

func ToReturnerSingle[T any](retAll matchers.ReturnerAll) matchers.ReturnerSingle[T] {
	return &returnerSingleImpl[T]{
		all: retAll,
	}
}

func ToReturnerDouble[A any, B any](retAll matchers.ReturnerAll) matchers.ReturnerDouble[A, B] {
	return &returnerDoubleImpl[A, B]{
		all: retAll,
	}
}

type returnerDummyImpl struct{}

func (r *returnerDummyImpl) ThenReturn(values ...any) matchers.ReturnerAll {
	return r
}

func (r *returnerDummyImpl) ThenAnswer(f matchers.Answer) matchers.ReturnerAll {
	return r
}

func (r *returnerDummyImpl) Verify(m matchers.MethodVerifier) {
}

type returnerAllImpl struct {
	methodMatch *methodMatch
	ctx         *mockContext
}

type returnerSingleImpl[T any] struct {
	all matchers.ReturnerAll
}

func (r *returnerSingleImpl[T]) ThenReturn(value T) matchers.ReturnerSingle[T] {
	return r.ThenAnswer(func(args []any) T {
		return value
	})
}

func (r *returnerSingleImpl[T]) ThenAnswer(f func(args []any) T) matchers.ReturnerSingle[T] {
	all := r.all.ThenAnswer(func(args []any) []any {
		return []any{f(args)}
	})
	return &returnerSingleImpl[T]{
		all: all,
	}
}

func (r *returnerSingleImpl[T]) Verify(verifier matchers.MethodVerifier) {
	r.all.Verify(verifier)
}

type returnerDoubleImpl[A any, B any] struct {
	all matchers.ReturnerAll
}

func (r *returnerDoubleImpl[A, B]) ThenReturn(a A, b B) matchers.ReturnerDouble[A, B] {
	return r.ThenAnswer(func(args []any) (A, B) {
		return a, b
	})
}

func (r *returnerDoubleImpl[A, B]) ThenAnswer(f func(args []any) (A, B)) matchers.ReturnerDouble[A, B] {
	all := r.all.ThenAnswer(func(args []any) []any {
		t, e := f(args)
		return []any{t, e}
	})
	return &returnerDoubleImpl[A, B]{
		all: all,
	}
}

func (r *returnerDoubleImpl[A, B]) Verify(verifier matchers.MethodVerifier) {
	r.all.Verify(verifier)
}

func (r *returnerAllImpl) ThenReturn(values ...any) matchers.ReturnerAll {
	return r.ThenAnswer(makeReturnFunc(values))
}

func (r *returnerAllImpl) ThenAnswer(f matchers.Answer) matchers.ReturnerAll {
	wrapper := &answerWrapper{
		ans: f,
	}
	r.methodMatch.addAnswer(wrapper)
	return r
}

func (r *returnerAllImpl) Verify(verifier matchers.MethodVerifier) {
	r.methodMatch.verifiers = append(r.methodMatch.verifiers, verifier)
}

func makeReturnFunc(values []any) matchers.Answer {
	return func(args []any) []interface{} {
		return values
	}
}

func NewReturnerAll(ctx *mockContext, data *methodMatch) matchers.ReturnerAll {
	return &returnerAllImpl{
		methodMatch: data,
		ctx:         ctx,
	}
}

func NewEmptyReturner() matchers.ReturnerAll {
	return &returnerDummyImpl{}
}
