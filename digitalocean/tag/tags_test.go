package tag_test

import (
	"reflect"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestDiffTags(t *testing.T) {
	cases := []struct {
		Old, New       *schema.Set
		Create, Remove map[string]string
	}{
		// Basic add/remove
		{
			Old: schema.NewSet(util.HashStringIgnoreCase, []interface{}{
				"foo",
			}),
			New: schema.NewSet(util.HashStringIgnoreCase, []interface{}{
				"bar",
			}),
			Create: map[string]string{
				"bar": "bar",
			},
			Remove: map[string]string{
				"foo": "foo",
			},
		},

		// Noop
		{
			Old: schema.NewSet(util.HashStringIgnoreCase, []interface{}{
				"foo",
			}),
			New: schema.NewSet(util.HashStringIgnoreCase, []interface{}{
				"foo",
			}),
			Create: map[string]string{},
			Remove: map[string]string{},
		},
	}

	for i, tc := range cases {
		r, c := tag.DiffTags(tag.TagsFromSchema(tc.Old), tag.TagsFromSchema(tc.New))
		if !reflect.DeepEqual(r, tc.Remove) {
			t.Fatalf("%d: bad remove: %#v", i, r)
		}
		if !reflect.DeepEqual(c, tc.Create) {
			t.Fatalf("%d: bad create: %#v", i, c)
		}
	}
}

func TestAccDigitalOceanTag_NameValidation(t *testing.T) {
	cases := []struct {
		Input       string
		ExpectError bool
	}{
		{
			Input:       "",
			ExpectError: true,
		},
		{
			Input:       "foo",
			ExpectError: false,
		},
		{
			Input:       "foo-bar",
			ExpectError: false,
		},
		{
			Input:       "foo:bar",
			ExpectError: false,
		},
		{
			Input:       "foo_bar",
			ExpectError: false,
		},
		{
			Input:       "foo-001",
			ExpectError: false,
		},
		{
			Input:       "foo/bar",
			ExpectError: true,
		},
		{
			Input:       "foo\bar",
			ExpectError: true,
		},
		{
			Input:       "foo.bar",
			ExpectError: true,
		},
		{
			Input:       "foo*",
			ExpectError: true,
		},
		{
			Input:       acctest.RandString(256),
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		_, errors := tag.ValidateTag(tc.Input, tc.Input)

		hasError := len(errors) > 0

		if tc.ExpectError && !hasError {
			t.Fatalf("Expected the DigitalOcean Tag Name to trigger a validation error for '%s'", tc.Input)
		}

		if hasError && !tc.ExpectError {
			t.Fatalf("Unexpected error validating the DigitalOcean Tag Name '%s': %s", tc.Input, errors[0])
		}
	}
}

func TestExpandTags(t *testing.T) {
	tags := []interface{}{"foo", "bar"}

	expandedTags := tag.ExpandTags(tags)

	if len(tags) != len(expandedTags) {
		t.Fatalf("incorrect expected length of expanded tags")
	}
}
func TestFlattenTags(t *testing.T) {
	tags := []string{"foo", "bar"}

	flattenedTags := tag.FlattenTags(tags)

	if len(tags) != flattenedTags.Len() {
		t.Fatalf("incorrect expected length of flattened tags")
	}
}
