// Task64 is a to-do list application to facilitate your life.
// Copyright 2018-present Vikunja and contributors. All rights reserved.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package health

import (
	"context"
	"errors"

	"github.com/OrnilioNeto/Task64/pkg/config"
	"github.com/OrnilioNeto/Task64/pkg/db"
	"github.com/OrnilioNeto/Task64/pkg/red"
)

// Check verifies the main service dependencies are reachable.
func Check() error {
	s := db.NewSession()
	defer s.Close()
	if err := s.Ping(); err != nil {
		return err
	}

	if config.RedisEnabled.GetBool() {
		r := red.GetRedis()
		if r == nil {
			return errors.New("redis not initialized")
		}
		if err := r.Ping(context.Background()).Err(); err != nil {
			return err
		}
	}
	return nil
}
