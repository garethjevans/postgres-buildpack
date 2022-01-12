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

package simple

import (
	"github.com/buildpacks/libcnb"
	"github.com/paketo-buildpacks/libpak/bard"
	"github.com/paketo-buildpacks/libpak/sherpa"
)

type Detect struct {
	Logger bard.Logger
}

func (d Detect) Detect(_ libcnb.DetectContext) (libcnb.DetectResult, error) {
	d.Logger.Info("Detect....")
	d.Logger.Info("Resolving BP_SIMPLE_BUILDPACK....")
	enableBuildpack := sherpa.ResolveBool("BP_SIMPLE_BUILDPACK")
	d.Logger.Infof("Found BP_SIMPLE_BUILDPACK=%s", enableBuildpack)

	if !enableBuildpack {
		d.Logger.Infof("not enabled")
		return libcnb.DetectResult{Pass: false}, nil
	}

	d.Logger.Infof("enabled!!!")

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
