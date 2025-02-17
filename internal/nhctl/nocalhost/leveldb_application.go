/*
Copyright 2021 The Nocalhost Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package nocalhost

import (
	"github.com/pkg/errors"
)

func UpdateKey(ns, app string, key string, value string) error {
	db, err := openApplicationLevelDB(ns, app)
	if err != nil {
		return err
	}
	defer db.Close()

	return errors.Wrap(db.Put([]byte(key), []byte(value), nil), "")
}
