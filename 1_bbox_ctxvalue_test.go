package syncs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/mark-ahn/syncs"
)

func Test_ContextValue(t *testing.T) {
	type inner_type struct {
		value string
	}

	ctx := context.Background()

	{
		{
			value, ok := syncs.ValueFrom[int](ctx)
			if ok {
				t.Fatalf("expect false, got true")
			}
			fmt.Println(value, ok)
		}
		ctx = syncs.WithValue(ctx, 100)
		{
			value, ok := syncs.ValueFrom[int](ctx)
			if !ok {
				t.Fatalf("expect true, got false")
			}
			if value != 100 {
				t.Fatalf("expect 100, got %v", value)
			}
			fmt.Println(value, ok)
		}
	}

	{
		{
			value, ok := syncs.ValueFrom[inner_type](ctx)
			if ok {
				t.Fatalf("expect false, got true")
			}
			fmt.Println(value, ok)
		}
		ctx = syncs.WithValue(ctx, inner_type{value: "some"})
		{
			value, ok := syncs.ValueFrom[inner_type](ctx)
			if !ok {
				t.Fatalf("expect true, got false")
			}
			expect := inner_type{value: "some"}
			if value != expect {
				t.Fatalf("expect %v, got %v", expect, value)
			}
			fmt.Println(value, ok)
		}
	}
}
