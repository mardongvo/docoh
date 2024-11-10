package db

import (
	"fmt"
	"io"
	"slices"

	"gopkg.in/yaml.v3"
)

func NewDB() *DB {
	return &DB{}
}

func (db *DB) Load(r io.Reader) error {
	dec := yaml.NewDecoder(r)
	err := dec.Decode(db)
	if err != nil {
		return fmt.Errorf("failed to decode db: %w", err)
	}
	return nil
}

func (db *DB) Save(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	defer enc.Close()

	enc.SetIndent(2)

	err := enc.Encode(db)
	if err != nil {
		return fmt.Errorf("failed to encode db: %w", err)
	}

	return nil
}

func (db *DB) AddRule(target, source string) error {
	id := 0
	for i, r := range db.Rules {
		id = max(id, r.ID)
		if r.Target != target {
			continue
		}
		for _, src := range r.Source {
			if src == source {
				return ErrDuplicateRule
			}
		}
		db.Rules[i].Source = append(db.Rules[i].Source, source)
		return nil
	}
	id = id + 1
	db.Rules = append(db.Rules, NewRule(id, target, source))
	return nil
}

func (db *DB) DeleteRule(id int) {
	db.Rules = slices.DeleteFunc(db.Rules, func(e Rule) bool {
		return e.ID == id
	})
}

func (db *DB) RuleByID(id int) (Rule, error) {
	idx := slices.IndexFunc(db.Rules, func(e Rule) bool {
		return e.ID == id
	})
	if idx == -1 {
		return Rule{}, ErrFoundNoRule
	}
	return db.Rules[idx], nil
}

func (db *DB) RuleByTarget(target string) (Rule, error) {
	idx := slices.IndexFunc(db.Rules, func(e Rule) bool {
		return e.Target == target
	})
	if idx == -1 {
		return Rule{}, ErrFoundNoRule
	}
	return db.Rules[idx], nil
}

func (db *DB) UpdateRule(rule Rule) error {
	idx := slices.IndexFunc(db.Rules, func(e Rule) bool {
		return e.ID == rule.ID
	})
	if idx == -1 {
		return ErrFoundNoRule
	}

	db.Rules[idx] = rule
	return nil
}
