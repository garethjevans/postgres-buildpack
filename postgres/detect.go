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

package postgres

import (
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/bindings"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) Detect(dc libcnb.DetectContext) (libcnb.DetectResult, error) {

	// check the binding first
	if _, ok, err := bindings.ResolveOne(dc.Platform.Bindings, bindings.OfType("postgresql")); err != nil {
		return libcnb.DetectResult{Pass: false}, nil
	} else if ok {
		return libcnb.DetectResult{
			Pass: true,
			Plans: []libcnb.BuildPlan{
				{
					Provides: []libcnb.BuildPlanProvide{
						{Name: "postgres-buildpack"},
					},
					Requires: []libcnb.BuildPlanRequire{
						{Name: "postgres-buildpack"},
						{Name: "jvm-application"},
					},
				},
			},
		}, nil
	}

	d.Logger.Info("Binding 'postgresql' not found, falling back on env var BP_POSTGRES_ENABLE....")

	if !sherpa.ResolveBool("BP_POSTGRES_ENABLE") {
		return libcnb.DetectResult{Pass: false}, nil
	}

	return libcnb.DetectResult{
		Pass: true,
		Plans: []libcnb.BuildPlan{
			{
				Provides: []libcnb.BuildPlanProvide{
					{Name: "postgres-buildpack"},
				},
				Requires: []libcnb.BuildPlanRequire{
					{Name: "postgres-buildpack"},
					{Name: "jvm-application"},
				},
			},
		},
	}, nil
}
