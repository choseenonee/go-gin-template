package customerr

const (
	UserRollbackError = Error("error on rollback in user repo")
	UserScanError     = Error("error on scan in user repo")
	UserCommitError   = Error("error on committing tx in user repo")
)
