package util

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestGetSetChanges(t *testing.T) {
	t.Parallel()

	tt := []struct {
		old    *schema.Set
		new    *schema.Set
		add    *schema.Set
		remove *schema.Set
	}{
		{
			old: schema.NewSet(schema.HashString, []interface{}{
				"foo", "bar", "baz",
			}),
			new: schema.NewSet(schema.HashString, []interface{}{
				"foo", "bar",
			}),
			add: schema.NewSet(schema.HashString, []interface{}{}),
			remove: schema.NewSet(schema.HashString, []interface{}{
				"baz",
			}),
		},
		{
			old: schema.NewSet(schema.HashString, []interface{}{
				"foo", "bar",
			}),
			new: schema.NewSet(schema.HashString, []interface{}{
				"foo", "bar", "baz",
			}),
			add: schema.NewSet(schema.HashString, []interface{}{
				"baz",
			}),
			remove: schema.NewSet(schema.HashString, []interface{}{}),
		},
		{
			old: schema.NewSet(schema.HashString, []interface{}{
				"foo",
			}),
			new: schema.NewSet(schema.HashString, []interface{}{
				"bar",
			}),
			add: schema.NewSet(schema.HashString, []interface{}{
				"bar",
			}),
			remove: schema.NewSet(schema.HashString, []interface{}{
				"foo",
			}),
		},
		{
			old: schema.NewSet(schema.HashString, []interface{}{
				"foo", "bar",
			}),
			new: schema.NewSet(schema.HashString, []interface{}{
				"baz", "qux",
			}),
			add: schema.NewSet(schema.HashString, []interface{}{
				"baz", "qux",
			}),
			remove: schema.NewSet(schema.HashString, []interface{}{
				"foo", "bar",
			}),
		},
		{
			old: schema.NewSet(schema.HashString, []interface{}{
				"foo", "bar", "baz",
			}),
			new: schema.NewSet(schema.HashString, []interface{}{
				"bar", "baz", "qux",
			}),
			add: schema.NewSet(schema.HashString, []interface{}{
				"qux",
			}),
			remove: schema.NewSet(schema.HashString, []interface{}{
				"foo",
			}),
		},
		{
			old: schema.NewSet(schema.HashInt, []interface{}{
				1, 2, 3,
			}),
			new: schema.NewSet(schema.HashInt, []interface{}{
				1, 2,
			}),
			add: schema.NewSet(schema.HashInt, []interface{}{}),
			remove: schema.NewSet(schema.HashInt, []interface{}{
				3,
			}),
		},
		{
			old: schema.NewSet(schema.HashInt, []interface{}{
				1,
			}),
			new: schema.NewSet(schema.HashInt, []interface{}{
				2,
			}),
			add: schema.NewSet(schema.HashInt, []interface{}{
				2,
			}),
			remove: schema.NewSet(schema.HashInt, []interface{}{
				1,
			}),
		},
	}

	for _, item := range tt {
		remove, add := GetSetChanges(item.old, item.new)
		if !remove.Equal(item.remove) || !add.Equal(item.add) {
			t.Errorf("expected add: %#v, remove: %#v; got add: %#v, remove: %#v", add, remove, item.add, item.remove)
		}
	}
}
