package main

import (
	"docoh/db"
	"docoh/files"
	"errors"
	"fmt"

	cli "github.com/urfave/cli/v2"
)

var (
	ErrNeedRuleNumberOrTarget = errors.New("specify rule number or target")
)

func addRule(ctx *cli.Context) error {
	dbpath := ctx.String("db")
	store, err := tryLoad(dbpath)
	if err != nil {
		return cli.Exit(fmt.Errorf("Failed to load db: %w", err), 1)
	}

	if ctx.NArg() != 2 {
		return cli.Exit("Specify target and source", 1)
	}

	target, source := ctx.Args().Get(0), ctx.Args().Get(1)
	err = store.AddRule(target, source)
	if err != nil {
		return err
	}

	err = trySave(store, dbpath)
	if err != nil {
		return cli.Exit(fmt.Errorf("Failed to save db: %w", err), 1)
	}

	return nil
}

func refreshRule(ctx *cli.Context) error {
	dbpath := ctx.String("db")
	store, err := tryLoad(dbpath)
	if err != nil {
		return cli.Exit(fmt.Errorf("Failed to load db: %w", err), 1)
	}

	rule, err := findRule(ctx, store)
	if err != nil {
		return cli.Exit(fmt.Errorf("Failed to get rule: %w", err), 1)
	}

	rule.Entries, err = buildEntriesForRule(rule)
	if err != nil {
		return cli.Exit(fmt.Errorf("Failed to build entries: %w", err), 1)
	}

	err = store.UpdateRule(rule)
	if err != nil {
		return cli.Exit(fmt.Errorf("Failed to update rule: %w", err), 1)
	}

	err = trySave(store, dbpath)
	if err != nil {
		return cli.Exit(fmt.Errorf("Failed to save db: %w", err), 1)
	}

	return nil
}

func reportRules(ctx *cli.Context) error {
	dbpath := ctx.String("db")
	store, err := tryLoad(dbpath)
	if err != nil {
		return cli.Exit(fmt.Errorf("Failed to load db: %w", err), 1)
	}

	var rules []db.Rule
	rule, err := findRule(ctx, store)
	switch {
	case err == nil:
		rules = []db.Rule{rule}
	case errors.Is(err, ErrNeedRuleNumberOrTarget):
		rules = store.Rules
	default:
		return cli.Exit(fmt.Errorf("Failed to get rule: %w", err), 1)
	}

	for _, rule := range rules {
		fmt.Printf("Rule #%d, Target: %s\n", rule.ID, rule.Target)

		entries, err := buildEntriesForRule(rule)
		if err != nil {
			return cli.Exit(fmt.Errorf("Failed to build entries: %w", err), 1)
		}

		for _, rs := range rule.Compare(entries) {
			switch rs.Result {
			case db.CompareNew:
				fmt.Printf("+ %s\n", rs.Path)
			case db.CompareDeleted:
				fmt.Printf("- %s\n", rs.Path)
			case db.CompareChanged:
				fmt.Printf("%s\n", rs.Path)
				fmt.Printf("\t-%s\n", rs.OldChecksum)
				fmt.Printf("\t+%s\n", rs.NewChecksum)
			}
		}
	}

	return nil
}

func findRule(ctx *cli.Context, store *db.DB) (db.Rule, error) {
	ruleID := ctx.Int("n")
	target := ctx.String("t")
	if ruleID == 0 && target == "" {
		return db.Rule{}, ErrNeedRuleNumberOrTarget
	}

	var rule db.Rule
	var err error
	switch {
	case ruleID > 0:
		rule, err = store.RuleByID(ruleID)
	case target != "":
		rule, err = store.RuleByTarget(target)
	}
	if err != nil {
		return db.Rule{}, err
	}

	return rule, nil
}

func buildEntriesForRule(rule db.Rule) ([]db.Entry, error) {
	filepaths, err := files.FilesForPatterns(rule.Source)
	if err != nil {
		return nil, err
	}

	entries := make([]db.Entry, 0, len(filepaths))
	for _, fp := range filepaths {
		cs, err := files.HashFile(fp)
		if err != nil {
			fmt.Printf("Failed to calculate file '%s' hash: %s\n", fp, err.Error())
			continue
		}
		entries = append(entries, db.Entry{
			Path:     fp,
			Checksum: cs,
		})
	}

	return entries, nil
}
