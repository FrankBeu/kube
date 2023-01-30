package kubeconfig

import (
	"context"
	"testing"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// TODO how to test secrets
// func TestDomainNameSuffix(t *testing.T) {
// 	type args struct {
// 		domain string
// 	}
// 	tests := []struct {
// 		name                 string
// 		args                 args
// 		wantDomainNameSuffix string
// 	}{
// 		{
// 			name:                 "test.domain",
// 			args:                 args{domain: "test.domain"},
// 			wantDomainNameSuffix: ".test.domain",
// 		},
// 		{
// 			name:                 "other.domain",
// 			args:                 args{domain: "other.domain"},
// 			wantDomainNameSuffix: ".other.domain",
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			ctx, err := pulumi.NewContext(context.Background(), pulumi.RunInfo{
// 				Config: map[string]string{
// 					// ":domain": tt.args.domain,
// 					// ":domain.subkey": subdomain,
// 					":domainSecret.domain": tt.args.domain,
// 				},
// 			})
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			if gotDomainNameSuffix := DomainNameSuffix(ctx); gotDomainNameSuffix != tt.wantDomainNameSuffix {
// 				t.Errorf("DomainNameSuffix() = %v, want %v", gotDomainNameSuffix, tt.wantDomainNameSuffix)
// 			}
// 		})
// 	}
// }

func TestEnv(t *testing.T) {
	type args struct {
		env string
	}
	tests := []struct {
		name    string
		args    args
		wantEnv string
	}{
		{
			name:    "dev",
			args:    args{env: "dev"},
			wantEnv: "dev",
		},
		{
			name:    "prod",
			args:    args{env: "prod"},
			wantEnv: "prod",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, err := pulumi.NewContext(context.Background(), pulumi.RunInfo{
				Config: map[string]string{
					":env": tt.args.env,
				},
			})
			if err != nil {
				t.Fatal(err)
			}
			if gotEnv := Env(ctx); gotEnv != tt.wantEnv {
				t.Errorf("Env() = %v, want %v", gotEnv, tt.wantEnv)
			}
		})
	}
}
