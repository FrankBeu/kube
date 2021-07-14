package config

import (
	"context"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func TestDomainNameSuffix(t *testing.T) {
	type args struct {
		domain string
	}
	tests := []struct {
		name                 string
		args                 args
		wantDomainNameSuffix string
	}{
		{
			name:                 "test.domain",
			args:                 args{domain: "test.domain"},
			wantDomainNameSuffix: ".test.domain",
		},
		{
			name:                 "other.domain",
			args:                 args{domain: "other.domain"},
			wantDomainNameSuffix: ".other.domain",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := pulumi.NewContext(context.Background(), pulumi.RunInfo{
				Config: map[string]string{
					":domain": tt.args.domain,
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			if gotDomainNameSuffix := DomainNameSuffix(ctx); gotDomainNameSuffix != tt.wantDomainNameSuffix {
				t.Errorf("DomainNameSuffix() = %v, want %v", gotDomainNameSuffix, tt.wantDomainNameSuffix)
			}
		})
	}
}
