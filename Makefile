gofumpt:
	gofumpt -l -w .

import:
	gci write --skip-generated -s standard -s default -s "prefix(github.com/ovechkin-dm/mockio)" -s blank -s dot -s alias .
