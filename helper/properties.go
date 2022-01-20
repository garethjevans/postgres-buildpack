/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package helper

import (
	"fmt"
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/bindings"
	"github.com/paketo-buildpacks/libpak/sherpa"
	"io/ioutil"
)

var (
	BindingName = "postgresql"
)

type ContentsGetter func(string) (string, error)

type Properties struct {
	Bindings   libcnb.Bindings
	Logger     bard.Logger
	FileReader ContentsGetter
}

func (c Properties) Execute() (map[string]string, error) {
	enabled := sherpa.ResolveBool("BPL_POSTGRES_PROPERTIES_ENABLE")
	environment := make(map[string]string)

	if enabled {
		c.Logger.Info("Configuring 'postgresql' properties")

		if b, ok, err := bindings.ResolveOne(c.Bindings, bindings.OfType(BindingName)); err != nil {
			return nil, fmt.Errorf("unable to resolve single binding %s\n%w", BindingName, err)
		} else if ok {
			c.Logger.Infof("Found binding %+v", b)
			if p, ok := b.SecretFilePath("url"); ok {
				c.Logger.Info("Configuring POSTGRES_URL")
				contents, err := c.GetContents(p)
				if err != nil {
					return nil, err
				}
				environment["POSTGRES_URL"] = contents
			}

			if p, ok := b.SecretFilePath("username"); ok {
				c.Logger.Info("Configuring POSTGRES_USER")
				contents, err := c.GetContents(p)
				if err != nil {
					return nil, err
				}
				environment["POSTGRES_USER"] = contents
			}

			if p, ok := b.SecretFilePath("password"); ok {
				c.Logger.Info("Configuring POSTGRES_PASS")
				contents, err := c.GetContents(p)
				if err != nil {
					return nil, err
				}
				environment["POSTGRES_PASS"] = contents
			}
		}

		c.Logger.Infof("Environment = %s", environment)
	} else {
		c.Logger.Info("Skipping configurin 'postgresql' properties")
	}

	return environment, nil
}

func (c Properties) GetContents(path string) (string, error) {
	if c.FileReader == nil {
		c.Logger.Infof("getting secret from path %s", path)
		b, err := ioutil.ReadFile(path)
		if err != nil {
			return "", fmt.Errorf("unable to read secret file %s\n%w", path, err)
		}
		return string(b), nil
	} else {
		return c.FileReader(path)
	}
}
