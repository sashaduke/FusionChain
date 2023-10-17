package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyring_IsParty(t *testing.T) {
	type fields struct {
		Id          uint64
		Creator     string
		Description string
		Admins      []string
		Parties     []string
	}
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Address exists in Parties",
			fields: fields{
				Parties: []string{"party1", "party2", "party3"},
			},
			args: args{
				address: "party2",
			},
			want: true,
		},
		{
			name: "Address does not exist in Parties",
			fields: fields{
				Parties: []string{"party1", "party3"},
			},
			args: args{
				address: "party2",
			},
			want: false,
		},
		{
			name: "Empty Parties",
			fields: fields{
				Parties: []string{},
			},
			args: args{
				address: "party2",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Keyring{
				Id:          tt.fields.Id,
				Creator:     tt.fields.Creator,
				Description: tt.fields.Description,
				Admins:      tt.fields.Admins,
				Parties:     tt.fields.Parties,
			}
			assert.Equalf(t, tt.want, k.IsParty(tt.args.address), "IsParty(%v)", tt.args.address)
		})
	}
}

func TestKeyring_AddParty(t *testing.T) {
	type fields struct {
		Id          uint64
		Creator     string
		Description string
		Admins      []string
		Parties     []string
	}
	type args struct {
		address string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []string
	}{
		{
			name: "Add a party to an empty Parties list",
			fields: fields{
				Parties: []string{},
			},
			args: args{
				address: "party1",
			},
			want: []string{"party1"},
		},
		{
			name: "Add a party to a non-empty Parties list",
			fields: fields{
				Parties: []string{"party1", "party2"},
			},
			args: args{
				address: "party3",
			},
			want: []string{"party1", "party2", "party3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &Keyring{
				Id:          tt.fields.Id,
				Creator:     tt.fields.Creator,
				Description: tt.fields.Description,
				Admins:      tt.fields.Admins,
				Parties:     tt.fields.Parties,
			}
			k.AddParty(tt.args.address)
		})
	}
}
