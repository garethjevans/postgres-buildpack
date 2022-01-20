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

package helper_test

import (
	"github.com/buildpacks/libcnb"
	"github.com/garethjevans/postgres-buildpack/helper"
	"os"
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sclevine/spec"
)

func testProperties(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		p helper.Properties
	)

	it("BPL_POSTGRES_PROPERTIES_ENABLE=false", func() {

		it("uses configured module", func() {
			Expect(p.Execute()).To(Equal(map[string]string{
				// empty
			}))
		})

		it("does not contribute credentials if 'postgresql' binding doesn't exist", func() {

			p.Bindings = libcnb.Bindings{
				{
					Name: "my-postgresql-service",
					Path: "/test/path/my-postgresql-service",
					Type: "postgresql",
					Secret: map[string]string{
						"url":      "jdbc:postgresql://localhost/petclinic",
						"username": "petclinic-user",
						"password": "petclinic-pass",
					},
				},
			}

			Expect(p.Execute()).To(Equal(map[string]string{
				// empty
			}))
		})

	})

	context("BPL_POSTGRES_PROPERTIES_ENABLE=true", func() {

		p.FileReader = func(in string) (string, error) {
			switch in {
			case "/test/path/my-postgresql-service/url":
				return "jdbc:postgresql://localhost/petclinic", nil
			case "/test/path/my-postgresql-service/username":
				return "petclinic-user", nil
			case "/test/path/my-postgresql-service/password":
				return "petclinic-pass", nil
			}
			return in, nil
		}
		
		os.Setenv("BPL_POSTGRES_PROPERTIES_ENABLE", "true")

		it("contributes credentials if 'postgresql' binding exists", func() {

			p.Bindings = libcnb.Bindings{
				{
					Name: "my-postgresql-service",
					Path: "/test/path/my-postgresql-service",
					Type: "postgresql",
					Secret: map[string]string{
						"url":      "jdbc:postgresql://localhost/petclinic",
						"username": "petclinic-user",
						"password": "petclinic-pass",
					},
				},
			}

			Expect(p.Execute()).To(Equal(map[string]string{
				"POSTGRES_URL":  "jdbc:postgresql://localhost/petclinic",
				"POSTGRES_USER": "petclinic-user",
				"POSTGRES_PASS": "petclinic-pass",
			}))
		})
	})
}
