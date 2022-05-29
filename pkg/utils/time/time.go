/*
 * Copyright 2022 The Furiko Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package time

import (
	"time"
)

func DurationMax(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}

func DurationMin(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}

func Max(a, b time.Time) time.Time {
	if a.Before(b) {
		return b
	}
	return a
}

func Min(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}

func MinNonZero(a, b time.Time) time.Time {
	if a.IsZero() {
		return b
	}
	return Min(a, b)
}
