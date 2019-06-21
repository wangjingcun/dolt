package alterschema

import (
	"context"
	"github.com/liquidata-inc/ld/dolt/go/libraries/doltcore/dtestutils"
	"github.com/liquidata-inc/ld/dolt/go/libraries/doltcore/row"
	"github.com/liquidata-inc/ld/dolt/go/libraries/doltcore/schema"
	"github.com/liquidata-inc/ld/dolt/go/store/go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRenameColumn(t *testing.T) {
	tests := []struct {
		name           string
		colName        string
		newName        string
		expectedSchema schema.Schema
		expectedRows   []row.Row
		expectedErr    string
	}{
		{
			name:    "rename int",
			colName: "age",
			newName: "newAge",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", dtestutils.IdTag, types.UUIDKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("name", dtestutils.NameTag, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("newAge", dtestutils.AgeTag, types.UintKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("is_married", dtestutils.IsMarriedTag, types.BoolKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("title", dtestutils.TitleTag, types.StringKind, false),
			),
			expectedRows: dtestutils.TypedRows,
		},
		{
			name:    "rename string",
			colName: "name",
			newName: "newName",
			expectedSchema: dtestutils.CreateSchema(
				schema.NewColumn("id", dtestutils.IdTag, types.UUIDKind, true, schema.NotNullConstraint{}),
				schema.NewColumn("newName", dtestutils.NameTag, types.StringKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("age", dtestutils.AgeTag, types.UintKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("is_married", dtestutils.IsMarriedTag, types.BoolKind, false, schema.NotNullConstraint{}),
				schema.NewColumn("title", dtestutils.TitleTag, types.StringKind, false),
			),
			expectedRows: dtestutils.TypedRows,
		},
		{
			name:        "column not found",
			colName:     "not found",
			expectedErr: "column not found",
		},
		{
			name:        "name collision",
			colName:     "age",
			newName:     "title",
			expectedErr: "two different columns with the same name exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dEnv := createEnvWithSeedData(t)
			ctx := context.Background()

			root, err := dEnv.WorkingRoot(ctx)
			assert.NoError(t, err)
			tbl, _ := root.GetTable(ctx, tableName)

			updatedTable, err := RenameColumn(ctx, dEnv.DoltDB, tbl, tt.colName, tt.newName)
			if len(tt.expectedErr) > 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
				return
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.expectedSchema, updatedTable.GetSchema(ctx))

			rowData := updatedTable.GetRowData(ctx)
			var foundRows []row.Row
			rowData.Iter(ctx, func(key, value types.Value) (stop bool) {
				foundRows = append(foundRows, row.FromNoms(tt.expectedSchema, key.(types.Tuple), value.(types.Tuple)))
				return false
			})

			assert.Equal(t, tt.expectedRows, foundRows)
		})
	}
}
