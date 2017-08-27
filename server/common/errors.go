package common

// Error handling types

type ErrRepoNotInitialized string

func (e ErrRepoNotInitialized) Error() string {
	return string(e)
}

type ErrRepoNotFound string

func (e ErrRepoNotFound) Error() string {
	return string(e)
}
