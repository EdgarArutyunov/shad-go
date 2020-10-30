// +build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

func get(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) (*kvapi.GetResponse, error) {
	var resp *kvapi.GetResponse
	for {
		tmp := ""
		oldVal := &tmp
		var err error
		resp, err = c.Get(&kvapi.GetRequest{
			Key: key,
		})

		var eAuth *kvapi.AuthError

		if err == nil {
			*oldVal = resp.Value
		} else if errors.Is(err, kvapi.ErrKeyNotFound) {
			oldVal = nil
			resp = &kvapi.GetResponse{
				Version: uuid.UUID{},
			}
		} else if errors.As(err, &eAuth) {
			return resp, err
		} else {
			continue
		}

		newVal, err := updateFn(oldVal)
		if err != nil {
			return resp, err
		}
		resp.Value = newVal
		break
	}
	return resp, nil
}

// UpdateValue ...
func UpdateValue(c kvapi.Client, key string, updateFn func(oldValue *string) (newValue string, err error)) error {
	for {
		getResp, err := get(c, key, updateFn)
		tryToSet := uuid.Must(uuid.NewV4())
		if err != nil {
			return err
		}

		for {
			_, err := c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      getResp.Value,
				OldVersion: getResp.Version,
				NewVersion: tryToSet,
			})

			var eConf *kvapi.ConflictError
			var eAuth *kvapi.AuthError

			if err == nil {
				return nil
			} else if errors.As(err, &eConf) {
				if eConf.ExpectedVersion == tryToSet {
					return nil
				}
				if eConf.ExpectedVersion != eConf.ProvidedVersion {
					break
				}
			} else if errors.As(err, &eAuth) {
				return err
			} else if errors.Is(err, kvapi.ErrKeyNotFound) {

				getResp.Value, err = updateFn(nil)
				getResp.Version = uuid.UUID{}
				if err != nil {
					return err
				}
			}
		}
	}
}
