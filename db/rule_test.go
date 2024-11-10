package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCompareEntries(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		r1 := Rule{}
		require.Len(t, r1.Compare(nil), 0)
	})

	t.Run("equal", func(t *testing.T) {
		entries := []Entry{{"a", "cs1"}, {"b", "cs2"}}
		r1 := Rule{
			Entries: entries,
		}
		require.Len(t, r1.Compare(entries), 0)
	})

	t.Run("change", func(t *testing.T) {
		r1 := Rule{
			Entries: []Entry{{"a", "cs1"}, {"b", "cs2"}},
		}
		entries := []Entry{{"a", "cs3"}, {"b", "cs4"}}
		cmp := r1.Compare(entries)
		require.Len(t, cmp, 2)
		require.Equal(t, cmp[0].Path, "a")
		require.Equal(t, cmp[0].Result, CompareChanged)
		require.Equal(t, cmp[1].Path, "b")
		require.Equal(t, cmp[1].Result, CompareChanged)
	})

	t.Run("new", func(t *testing.T) {
		r1 := Rule{
			Entries: []Entry{{"a", "cs1"}},
		}
		entries := []Entry{{"a", "cs1"}, {"b", "cs4"}}
		cmp := r1.Compare(entries)
		require.Len(t, cmp, 1)
		require.Equal(t, cmp[0].Path, "b")
		require.Equal(t, cmp[0].Result, CompareNew)
	})

	t.Run("deleted", func(t *testing.T) {
		r1 := Rule{
			Entries: []Entry{{"a", "cs1"}, {"b", "cs4"}},
		}
		entries := []Entry{{"a", "cs1"}}
		cmp := r1.Compare(entries)
		require.Len(t, cmp, 1)
		require.Equal(t, cmp[0].Path, "b")
		require.Equal(t, cmp[0].Result, CompareDeleted)
	})
}
