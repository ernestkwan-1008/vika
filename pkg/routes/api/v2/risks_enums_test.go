// Vikunja is a to-do list application to facilitate your life.
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

package apiv2

import (
	"reflect"
	"strings"
	"testing"

	"code.vikunja.io/api/pkg/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// The enum tags of the risk operations are literals, and Huma validates them before the model
// runs. These tests pin them to the model's own lists, so adding a status, a rating or a sort key
// in one place and not the other fails here instead of rejecting valid requests at runtime.

func tagValues(t *testing.T, typ reflect.Type, field, tag string) []string {
	t.Helper()
	f, ok := typ.FieldByName(field)
	require.True(t, ok, "%s has no field %s", typ, field)
	value := f.Tag.Get(tag)
	require.NotEmpty(t, value, "%s.%s has no %s tag", typ, field, tag)
	return strings.Split(value, ",")
}

func TestRiskEnumTagsMatchTheModel(t *testing.T) {
	listParams := reflect.TypeOf(RiskListQueryParams{})

	assert.ElementsMatch(t, models.RiskStatuses(), tagValues(t, listParams, "Status", "enum"), "status filter")
	assert.ElementsMatch(t, models.RiskStatuses(), tagValues(t, reflect.TypeOf(riskStatusRequest{}), "Status", "enum"), "status operation")
	assert.ElementsMatch(t, models.RiskStatuses(), tagValues(t, reflect.TypeOf(models.Risk{}), "Status", "enum"), "risk status field")

	assert.ElementsMatch(t, models.RiskRatings(), tagValues(t, listParams, "Rating", "enum"), "rating filter")
	assert.ElementsMatch(t, models.RiskRatings(), tagValues(t, reflect.TypeOf(models.Risk{}), "Rating", "enum"), "risk rating field")

	assert.ElementsMatch(t, models.RiskSortKeys(), tagValues(t, listParams, "SortBy", "enum"), "sort keys")
}

func TestRiskScoreRangesInDocsMatchTheBands(t *testing.T) {
	want := models.RiskScoreRangeDoc()
	require.Equal(t, "low 1-4, medium 5-9, high 10-16, critical 17-25", want)

	for typ, field := range map[reflect.Type]string{
		reflect.TypeOf(RiskListQueryParams{}): "Rating",
		reflect.TypeOf(models.Risk{}):         "Rating",
	} {
		f, ok := typ.FieldByName(field)
		require.True(t, ok)
		assert.Contains(t, f.Tag.Get("doc"), want, "%s.%s doc", typ, field)
	}
}
