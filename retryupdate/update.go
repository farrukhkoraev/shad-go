//go:build !solution

package retryupdate

import (
	"errors"
	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func UpdateValue(
	c kvapi.Client,
	key string,
	updateFn func(oldValue *string) (newValue string, err error),
) error {
	var (
		authErr     *kvapi.AuthError
		conflictErr *kvapi.ConflictError
	)

	needFetch := true
	var oldValue *string
	var curVersion uuid.UUID
	var newValue string
	var newVersion uuid.UUID

	for {
		if needFetch {
			// 1. get current value
			resp, err := c.Get(&kvapi.GetRequest{Key: key})
			if err != nil {
				if errors.As(err, &authErr) {
					return err
				}
				if errors.Is(err, kvapi.ErrKeyNotFound) {
					oldValue = nil
					curVersion = uuid.Nil
				} else {
					// temporary error — retry get
					continue
				}
			} else {
				curVersion = resp.Version
				oldValue = &resp.Value
			}

			// 2. compute new value
			var err2 error
			newValue, err2 = updateFn(oldValue)
			if err2 != nil {
				return err2
			}
			newVersion = uuid.Must(uuid.NewV4())
			needFetch = false
		}

		// 3. set value
		_, err := c.Set(&kvapi.SetRequest{
			Key:        key,
			Value:      newValue,
			OldVersion: curVersion,
			NewVersion: newVersion,
		})

		if errors.As(err, &conflictErr) {
			// false conflict: our set already succeeded
			if conflictErr.ExpectedVersion == newVersion {
				return nil
			}
			// real conflict: re-fetch and retry
			needFetch = true
			continue
		}
		if errors.Is(err, kvapi.ErrKeyNotFound) {
			// key vanished: recompute value as new key
			curVersion = uuid.Nil
			var err2 error
			newValue, err2 = updateFn(nil)
			if err2 != nil {
				return err2
			}
			newVersion = uuid.Must(uuid.NewV4())
			continue
		}

		if err == nil {
			return nil
		}
		if errors.As(err, &authErr) {
			return err
		}

		// temporary set error: retry set without re-fetching
	}
}
