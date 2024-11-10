package db

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDB(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		db := NewDB()
		require.Len(t, db.Rules, 0)
	})

	t.Run("add", func(t *testing.T) {
		db := NewDB()
		err := db.AddRule("t1", "s1")
		require.Len(t, db.Rules, 1)
		require.NoError(t, err)

		err = db.AddRule("t1", "s1")
		require.Len(t, db.Rules, 1)
		require.ErrorIs(t, err, ErrDuplicateRule)

		err = db.AddRule("t2", "s2")
		require.Len(t, db.Rules, 2)
		require.NoError(t, err)
	})

	t.Run("add for one target", func(t *testing.T) {
		db := NewDB()
		err := db.AddRule("t1", "s1")
		require.Len(t, db.Rules, 1)
		require.NoError(t, err)

		err = db.AddRule("t1", "s2")
		require.Len(t, db.Rules, 1)
		require.NoError(t, err)

		rule := db.Rules[0]
		require.Len(t, rule.Source, 2)
	})

	t.Run("del", func(t *testing.T) {
		const n = 100

		db := NewDB()

		for i := range n {
			err := db.AddRule(fmt.Sprintf("t%d", i), "s")
			require.NoError(t, err)
		}
		require.Len(t, db.Rules, n)

		for i := range n {
			db.DeleteRule(i + 1)
			require.Len(t, db.Rules, n-i-1)
		}
		require.Len(t, db.Rules, 0)
	})

	t.Run("save-load", func(t *testing.T) {
		db := NewDB()
		db.Rules = append(db.Rules, Rule{
			Target: "t1",
			Source: []string{"s1"},
			Entries: []Entry{
				{"e1", "cs1"},
				{"e2", "cs3"},
			},
		})
		db.Rules = append(db.Rules, Rule{
			Target: "t2",
			Source: []string{"s2"},
			Entries: []Entry{
				{"E1", "Cs1"},
				{"E2", "Cs3"},
			},
		})

		var buf bytes.Buffer
		err := db.Save(&buf)
		require.NoError(t, err)

		db2 := NewDB()
		err = db2.Load(&buf)
		require.NoError(t, err)

		require.Equal(t, db, db2)
	})

}
