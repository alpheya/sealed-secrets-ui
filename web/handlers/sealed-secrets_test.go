package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyValuePairs(t *testing.T) {
	tcs := []struct {
		name      string
		field     string
		want      map[string]string
		isWantErr bool
	}{
		{
			name:      "empty",
			field:     "",
			isWantErr: true,
		},
		{
			name:  "one oneline value",
			field: `API_TOKEN=adf23456-.,=%!+/\\-_()[]{}+@?,'"#`,
			want: map[string]string{
				"API_TOKEN": `adf23456-.,=%!+/\\-_()[]{}+@?,'"#`,
			},
		},
		{
			name:  "one oneline value with backticked block",
			field: `API_TOKEN=` + "`" + `adf23456-.,=%!+/\\-_()[]{}+@?,'"#` + "`",
			want: map[string]string{
				"API_TOKEN": `adf23456-.,=%!+/\\-_()[]{}+@?,'"#`,
			},
		},
		{
			name:  "one oneline value contains multiple equal signs",
			field: `API_TOKEN=adf23=456==`,
			want: map[string]string{
				"API_TOKEN": "adf23=456==",
			},
		},
		{
			name:  "one oneline value contains escaped backtick",
			field: `API_TOKEN=adf23456-.,=%!+/` + escapedBacktick + `\\-_()[]{}+@?,'"#`,
			want: map[string]string{
				"API_TOKEN": `adf23456-.,=%!+/` + escapedBacktick + `\\-_()[]{}+@?,'"#`,
			},
		},
		{
			name:  "one oneline value contains non escaped backtick",
			field: `API_TOKEN=adf23456-.,=%!+/` + "`" + `\\-_()[]{}+@?,'"#`,
			want: map[string]string{
				"API_TOKEN": `adf23456-.,=%!+/` + "`" + `\\-_()[]{}+@?,'"#`,
			},
			// this succeeds because the parser checks only the last character
			// in the specific line.
		},
		{
			name:  "one oneline value contains non escaped backtick at end",
			field: `API_TOKEN=adf23456-.,=%!+/\\-_()[]{}+@?,'"#` + "`",
			want: map[string]string{
				"API_TOKEN": `adf23456-.,=%!+/\\-_()[]{}+@?,'"#` + "`",
			},
		},
		{
			name:  "one oneline value contains escaped backtick at end",
			field: `API_TOKEN=adf23456-.,=%!+/\\-_()[]{}+@?,'"#` + escapedBacktick,
			want: map[string]string{
				"API_TOKEN": `adf23456-.,=%!+/\\-_()[]{}+@?,'"#` + escapedBacktick,
			},
		},
		{
			name: "one multiline value",
			field: `PRIV_KEY=` + "`" + `--begin
some-val
--end` + "`",
			want: map[string]string{
				"PRIV_KEY": `--begin
some-val
--end`,
			},
		},
		{
			name: "multiline value contains escaped backtick in middle",
			field: `PRIV_KEY=` + "`" + `--begin
some-val` + escapedBacktick + `
--end` + "`",
			want: map[string]string{
				"PRIV_KEY": `--begin
some-val` + "`" + `
--end`,
			},
		},
		{
			name: "multiline value contains non escaped backtick in middle",
			field: `PRIV_KEY=` + "`" + `--begin
some-val` + "`" + `
--end` + "`",
			isWantErr: true,
		},
		{
			name: "multiline value contains equal signs in middle",
			field: `PRIV_KEY=` + "`" + `--begin
some===-val
--end` + "`",
			want: map[string]string{
				"PRIV_KEY": `--begin
some===-val
--end`,
			},
		},
		{
			name: "multiline value contains equal signs in first line",
			field: `PRIV_KEY=` + "`" + `--be==gin==
some-val
--end` + "`",
			want: map[string]string{
				"PRIV_KEY": `--be==gin==
some-val
--end`,
			},
		},
		{
			name: "mixed",
			field: `API_TOKEN=adf23456-.,=%!+/\\-_()[]{}+@?,'"#
PRIV_KEY=` + "`" + `--begin
some-val
--end` + "`" + `
ENV_VAR=some-value=/
PUBLICK=` + "`" + `qwertz
12345
xcvb` + "`",
			want: map[string]string{
				"API_TOKEN": `adf23456-.,=%!+/\\-_()[]{}+@?,'"#`,
				"PRIV_KEY": `--begin
some-val
--end`,
				"ENV_VAR": "some-value=/",
				"PUBLICK": `qwertz
12345
xcvb`,
			},
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			gotResult, gotErr := parseKeyValuePairs(tc.field)
			assert.Equal(t, tc.want, gotResult)
			if tc.isWantErr {
				assert.Error(t, gotErr)
			} else {
				assert.NoError(t, gotErr)
			}
		})
	}
}
