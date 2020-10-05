package cmd

import (
	"errors"
	"reflect"
	"testing"

	"github.com/hellofresh/kangal/pkg/controller"
)

func TestControllerPopulateCfgFromOpts(t *testing.T) {
	type fields struct {
		kubeConfig           string
		masterURL            string
		namespaceAnnotations []string
		podAnnotations       []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    controller.Config
		wantErr error
	}{
		{
			name:   "test with empty annotations",
			fields: fields{},
			want: controller.Config{
				NamespaceAnnotations: map[string]string{},
				PodAnnotations:       map[string]string{},
			},
		},
		{
			name: "test with aws annotations",
			fields: fields{
				namespaceAnnotations: []string{"iam.amazonaws.com/permitted:.*"},
				podAnnotations:       []string{"iam.amazonaws.com/role:arn:aws:iam::someid:role/some-role-name"},
			},
			want: controller.Config{
				NamespaceAnnotations: map[string]string{"iam.amazonaws.com/permitted": ".*"},
				PodAnnotations:       map[string]string{"iam.amazonaws.com/role": "arn:aws:iam::someid:role/some-role-name"},
			},
		},
		{
			name: `test with some "`,
			fields: fields{
				namespaceAnnotations: []string{`iam.amazonaws.com/permitted:".*"`},
				podAnnotations:       []string{`iam.amazonaws.com/role:arn:aws:iam::"someid:role/some-role-name"`},
			},
			want: controller.Config{
				NamespaceAnnotations: map[string]string{"iam.amazonaws.com/permitted": ".*"},
				PodAnnotations:       map[string]string{"iam.amazonaws.com/role": "arn:aws:iam::someid:role/some-role-name"},
			},
		},
		{
			name: "test with invalid namespace annotations",
			fields: fields{
				namespaceAnnotations: []string{"iam.amazonaws.com/permitted.*"},
			},
			want:    controller.Config{},
			wantErr: errors.New("failed to convert namepsace annotations: Annotation \"iam.amazonaws.com/permitted.*\" is invalid"),
		},
		{
			name: "test with invalid pod annotations",
			fields: fields{
				namespaceAnnotations: []string{"iam.amazonaws.com/permitted:.*"},
				podAnnotations:       []string{"iam.amazonaws.com/role_arn_aws_iam__someid_role/some-role-name"},
			},
			want:    controller.Config{},
			wantErr: errors.New("failed to convert pod annotations: Annotation \"iam.amazonaws.com/role_arn_aws_iam__someid_role/some-role-name\" is invalid"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &controllerCmdOptions{
				kubeConfig:           tt.fields.kubeConfig,
				masterURL:            tt.fields.masterURL,
				namespaceAnnotations: tt.fields.namespaceAnnotations,
				podAnnotations:       tt.fields.podAnnotations,
			}

			got, gotErr := populateCfgFromOpts(controller.Config{}, opts)

			if tt.wantErr != nil && gotErr == nil ||
				tt.wantErr != nil && gotErr != nil && tt.wantErr.Error() != gotErr.Error() {
				t.Errorf("err = %v, want %v", gotErr, tt.wantErr)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("populateCfgFromOpts() = %v, want %v", got, tt.want)
			}
		})
	}
}
