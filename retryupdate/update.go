// +build !solution

package retryupdate

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/slon/shad-go/retryupdate/kvapi"
)

type UpdateFnType func(*string) (string, error)

func get(c kvapi.Client, key string, updateFn UpdateFnType) (*kvapi.GetResponse, error) {
	var resp *kvapi.GetResponse
	var eAuth *kvapi.AuthError
	var err error

	oldVal := new(string)

Loop:
	for {
		resp, err = c.Get(&kvapi.GetRequest{
			Key: key,
		})

		switch {
		case err == nil:
			*oldVal = resp.Value
			break Loop

		case errors.Is(err, kvapi.ErrKeyNotFound):
			oldVal = nil
			resp = &kvapi.GetResponse{
				Version: uuid.UUID{},
			}
			break Loop

		case errors.As(err, &eAuth):
			return nil, err

		default:
			continue
		}
	}

	resp.Value, err = updateFn(oldVal)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// UpdateValue ...
func UpdateValue(c kvapi.Client, key string, updateFn UpdateFnType) error {
	for {
		resp, err := get(c, key, updateFn)
		if err != nil {
			return err
		}

		newVersion := uuid.Must(uuid.NewV4())
	Loop:
		for {
			_, err := c.Set(&kvapi.SetRequest{
				Key:        key,
				Value:      resp.Value,
				OldVersion: resp.Version,
				NewVersion: newVersion,
			})

			var eConf *kvapi.ConflictError
			var eAuth *kvapi.AuthError

			switch {
			case err == nil:
				return nil

			case errors.As(err, &eAuth):
				return err

			case errors.As(err, &eConf):
				switch {
				case eConf.ExpectedVersion == newVersion:
					return nil

				case eConf.ExpectedVersion != eConf.ProvidedVersion:
					break Loop
				}

			case errors.Is(err, kvapi.ErrKeyNotFound):
				resp.Value, err = updateFn(nil)
				resp.Version = uuid.UUID{}
				if err != nil {
					return err
				}
			}
		}
	}
}
