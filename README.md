# docoh
Documentation coherence utility

## Purpose
1. Materialize the links between source code and high-level documentation.
2. Keep track of changes in the source code and give a report
 in case you need to check and correct the documentation.

Files checksum is the same as git object id
(see [Object storage](https://git-scm.com/book/en/v2/Git-Internals-Git-Objects) section).

## Usage

### Storage
All changes will be stored in yaml file (default is .docohdb, use flag `-db`)

### Add rule
`docoh add <target> <filepath | glob pattern>`

Multiple `add` with one target will be saved as one rule.

### Initialize/refresh files checksums
`docoh refresh <-n rulenumber | -t target>`

### Report changes
`docoh report [-n rulenumber | -t target]`