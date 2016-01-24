package do

import (
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
)

// SSHKey wraps godo Key.
type SSHKey struct {
	*godo.Key
}

// SSHKeys is a slice of SSHKey
type SSHKeys []SSHKey

// KeysService is the godo KeysService interface.
type KeysService interface {
	List() (SSHKeys, error)
	Get(id string) (*SSHKey, error)
	Create(kcr *godo.KeyCreateRequest) (*SSHKey, error)
	Update(id string, kur *godo.KeyUpdateRequest) (*SSHKey, error)
	Delete(id string) error
}

type keysService struct {
	client *godo.Client
}

var _ KeysService = &keysService{}

// NewKeysService builds an instance of KeysService.
func NewKeysService(client *godo.Client) KeysService {
	return &keysService{
		client: client,
	}
}

func (ks *keysService) List() (SSHKeys, error) {
	f := func(opt *godo.ListOptions) ([]interface{}, *godo.Response, error) {
		list, resp, err := ks.client.Keys.List(opt)
		if err != nil {
			return nil, nil, err
		}

		si := make([]interface{}, len(list))
		for i := range list {
			si[i] = list[i]
		}

		return si, resp, err
	}

	si, err := PaginateResp(f)
	if err != nil {
		return nil, err
	}

	list := make(SSHKeys, len(si))
	for i := range si {
		k := si[i].(godo.Key)
		list[i] = SSHKey{Key: &k}
	}

	return list, nil
}

func (ks *keysService) Get(id string) (*SSHKey, error) {
	var err error
	var k *godo.Key

	if i, aerr := strconv.Atoi(id); aerr == nil {
		k, _, err = ks.client.Keys.GetByID(i)
	} else {
		if len(id) > 0 {
			k, _, err = ks.client.Keys.GetByFingerprint(id)
		} else {
			err = fmt.Errorf("missing key id or fingerprint")
		}
	}

	if err != nil {
		return nil, err
	}

	return &SSHKey{Key: k}, nil
}

func (ks *keysService) Create(kcr *godo.KeyCreateRequest) (*SSHKey, error) {
	k, _, err := ks.client.Keys.Create(kcr)
	if err != nil {
		return nil, err
	}

	return &SSHKey{Key: k}, nil
}

func (ks *keysService) Update(id string, kur *godo.KeyUpdateRequest) (*SSHKey, error) {
	var k *godo.Key
	var err error
	if i, aerr := strconv.Atoi(id); aerr == nil {
		k, _, err = ks.client.Keys.UpdateByID(i, kur)
	} else {
		k, _, err = ks.client.Keys.UpdateByFingerprint(id, kur)
	}

	if err != nil {
		return nil, err
	}

	return &SSHKey{Key: k}, nil
}

func (ks *keysService) Delete(id string) error {
	var err error

	if i, aerr := strconv.Atoi(id); aerr == nil {
		_, err = ks.client.Keys.DeleteByID(i)
	} else {
		_, err = ks.client.Keys.DeleteByFingerprint(id)
	}

	return err
}