/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package checker

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/onflow/cadence/runtime/sema"
)

func TestArrayElementTypeInference(t *testing.T) {

	t.Parallel()

	t.Run("numeric array", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: [Int8] = [1, 2, 3]
				var y: [Int8]? = [1, 2, 3]
			}
		`)

		require.NoError(t, err)
	})

	t.Run("anystruct array", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: [AnyStruct] = [1, 2, 3]
			}
		`)

		require.NoError(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: [Int8]? = [1, 534, 3]
			}
		`)

		require.Error(t, err)
		assert.IsType(t, &sema.CheckerError{}, err)
	})

	t.Run("anystruct", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: AnyStruct = [1, 534534, 3]
			}
		`)

		require.NoError(t, err)
	})

	t.Run("inferring from rhs", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x = [1, 534534, 3]
				var y: Int = x[2]
			}
		`)

		require.NoError(t, err)
	})

	t.Run("nested array", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: [[[Int8]]] = [[[1, 2, 3], [4]], []]
			}
		`)

		require.NoError(t, err)
	})
}

func TestDictionaryTypeInference(t *testing.T) {

	t.Parallel()

	t.Run("numerics", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: {Int8:Int64} = {0: 6, 1: 5}
				var y: {Int8:Int64?} = {0: 6, 1: 5}
			}
		`)

		require.NoError(t, err)
	})

	t.Run("heterogeneous", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: {Int:AnyStruct} = {0: 6, 1: "hello", 2: nil}
			}
		`)

		require.NoError(t, err)
	})

	t.Run("nested", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: {Int:{Int:{Int:AnyStruct}}} = {0: {0: {0: 6}, 1: {0: 7}}}
			}
		`)

		require.NoError(t, err)
	})

	t.Run("invalid", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: {Int8:Int64} = {23423:6, 1:5}
			}
		`)

		require.Error(t, err)
		assert.IsType(t, &sema.CheckerError{}, err)
	})

	t.Run("nested invalid", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: {Int:{Int:{Int:Int8}}} = {0: {0: {0: "hello"}, 1: {0: 7}}}
			}
		`)

		require.Error(t, err)
		assert.IsType(t, &sema.CheckerError{}, err)
	})
}

func TestReturnTypeInference(t *testing.T) {

	t.Parallel()

	_, err := ParseAndCheck(t, `
		fun test(): [Int8] {
			return [1, 2, 3]
		}
	`)

	require.NoError(t, err)
}

func TestFunctionArgumentTypeInference(t *testing.T) {

	t.Parallel()

	_, err := ParseAndCheck(t, `
		fun test() {
			foo(a: [1, 2, 3])
		}

		fun foo(a: [Int8]) {
		}
	`)

	// Type inferring for function arguments is not supported yet.
	require.Error(t, err)
	assert.IsType(t, &sema.CheckerError{}, err)
}

func TestBinaryExpressionTypeInference(t *testing.T) {

	t.Parallel()

	t.Run("integer add", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: Int8 = 5 + 6
			}
		`)

		require.NoError(t, err)
	})

	t.Run("fixed point add", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: Fix64 = 5.0 + 6.0
			}
		`)

		require.NoError(t, err)
	})
}

func TestUnaryExpressionTypeInference(t *testing.T) {

	t.Parallel()

	t.Run("invalid negate", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var b : Bool =  !"string"
			}
		`)

		require.Error(t, err)
		assert.IsType(t, &sema.CheckerError{}, err)

		checkerErr := err.(*sema.CheckerError)

		// Should only give a single error
		require.Len(t, checkerErr.Errors, 1)

		assert.IsType(t, &sema.InvalidUnaryOperandError{}, checkerErr.Errors[0])
	})
}

func TestForceExpressionTypeInference(t *testing.T) {

	t.Parallel()

	t.Run("array forced", func(t *testing.T) {
		_, err := ParseAndCheck(t, `
			fun test() {
				var x: [Int8] = [5, 7, 2]!
			}
		`)

		require.NoError(t, err)
	})
}