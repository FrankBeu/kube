// Package tyk provides the default tykGateway-installation
package tyk

import (
	"fmt"

	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/apps/v1"
	corev1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/core/v1"
	metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"thesym.site/kube/lib/namespace"
)

func CreateTykGateway(ctx *pulumi.Context) error {

	namespaceTykGateway := &namespace.Namespace{
		Name:          "tyk",
		Tier:          namespace.NamespaceTierEdge,
		GlooDiscovery: true,
	}

	_, err := namespace.CreateNamespace(ctx, namespaceTykGateway)

	if err != nil {
		return err
	}

	_, err = corev1.NewConfigMap(ctx, "tyk_gateway_confConfigMap", &corev1.ConfigMapArgs{
		ApiVersion: pulumi.String("v1"),
		Data: pulumi.StringMap{
			"tyk.conf": pulumi.String(fmt.Sprintf("%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v%v", "{\n", "      \"listen_port\": 8080,\n", "      \"secret\": \"352d20ee67be67f6340b4c0605b044b7\",\n", "      \"template_path\": \"/opt/tyk-gateway/templates\",\n", "      \"tyk_js_path\": \"/opt/tyk-gateway/js/tyk.js\",\n", "      \"middleware_path\": \"/opt/tyk-gateway/middleware\",\n", "      \"use_db_app_configs\": false,\n", "      \"app_path\": \"/opt/tyk-gateway/apps/\",\n", "      \"storage\": {\n", "        \"type\": \"redis\",\n", "        \"host\": \"tyk-redis\",\n", "        \"port\": 6379,\n", "        \"username\": \"\",\n", "        \"password\": \"\",\n", "        \"database\": 0,\n", "        \"optimisation_max_idle\": 2000,\n", "        \"optimisation_max_active\": 4000\n", "      },\n", "      \"enable_analytics\": false,\n", "      \"analytics_config\": {\n", "        \"type\": \"csv\",\n", "        \"csv_dir\": \"/tmp\",\n", "        \"mongo_url\": \"\",\n", "        \"mongo_db_name\": \"\",\n", "        \"mongo_collection\": \"\",\n", "        \"purge_delay\": -1,\n", "        \"ignored_ips\": []\n", "      },\n", "      \"health_check\": {\n", "        \"enable_health_checks\": true,\n", "        \"health_check_value_timeouts\": 60\n", "      },\n", "      \"optimisations_use_async_session_write\": true,\n", "      \"enable_non_transactional_rate_limiter\": true,\n", "      \"enable_sentinel_rate_limiter\": false,\n", "      \"enable_redis_rolling_limiter\": false,\n", "      \"allow_master_keys\": false,\n", "      \"policies\": {\n", "        \"policy_source\": \"file\",\n", "        \"policy_record_name\": \"/opt/tyk-gateway/policies/policies.json\"\n", "      },\n", "      \"hash_keys\": true,\n", "      \"close_connections\": false,\n", "      \"http_server_options\": {\n", "        \"enable_websockets\": true\n", "      },\n", "      \"allow_insecure_configs\": true,\n", "      \"coprocess_options\": {\n", "        \"enable_coprocess\": true,\n", "        \"coprocess_grpc_server\": \"\"\n", "      },\n", "      \"enable_bundle_downloader\": true,\n", "      \"bundle_base_url\": \"\",\n", "      \"global_session_lifetime\": 100,\n", "      \"force_global_session_lifetime\": false,\n", "      \"max_idle_connections_per_host\": 500\n", "    }\n")),
		},
		Kind: pulumi.String("ConfigMap"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("tyk-gateway-conf"),
			Namespace: pulumi.String(namespaceTykGateway.Name),
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "tyk_gtwDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("tyk-gtw"),
			Namespace: pulumi.String(namespaceTykGateway.Name),
			Labels: pulumi.StringMap{
				"app": pulumi.String("tyk-gtw"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Replicas: pulumi.Int(1),
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String("tyk-gtw"),
				},
			},
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String("tyk-gtw"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:            pulumi.String("tyk-gtw"),
							Image:           pulumi.String("tykio/tyk-gateway:v3.1.0"),
							ImagePullPolicy: pulumi.String("Always"),
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(8089),
								},
							},
							Env: corev1.EnvVarArray{
								&corev1.EnvVarArgs{
									Name:  pulumi.String("TYK_GW_LISTENPORT"),
									Value: pulumi.String("8080"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("TYK_GW_SECRET"),
									Value: pulumi.String("foo"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("TYK_GW_STORAGE_HOST"),
									Value: pulumi.String("redis"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("TYK_GW_STORAGE_PORT"),
									Value: pulumi.String("6379"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("TYK_GW_STORAGE_PASSWORD"),
									Value: pulumi.String(""),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("TYK_LOGLEVEL"),
									Value: pulumi.String("info"),
								},
								&corev1.EnvVarArgs{
									Name:  pulumi.String("GODEBUG"),
									Value: pulumi.String("netdns=cgo"),
								},
							},
							VolumeMounts: corev1.VolumeMountArray{
								&corev1.VolumeMountArgs{
									Name:      pulumi.String("tyk-gateway-conf"),
									MountPath: pulumi.String("/opt/tyk-gateway/tyk.conf"),
									SubPath:   pulumi.String("tyk.conf"),
								},
							},
							Resources: &corev1.ResourceRequirementsArgs{
								Limits: pulumi.StringMap{
									"memory": pulumi.String("512Mi"),
									"cpu":    pulumi.String("1"),
								},
								Requests: pulumi.StringMap{
									"memory": pulumi.String("256Mi"),
									"cpu":    pulumi.String("0.2"),
								},
							},
						},
					},
					Volumes: corev1.VolumeArray{
						&corev1.VolumeArgs{
							Name: pulumi.String("tyk-gateway-conf"),
							ConfigMap: &corev1.ConfigMapVolumeSourceArgs{
								Name: pulumi.String("tyk-gateway-conf"),
								Items: corev1.KeyToPathArray{
									&corev1.KeyToPathArgs{
										Key:  pulumi.String("tyk.conf"),
										Path: pulumi.String("tyk.conf"),
									},
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewService(ctx, "tyk_svcService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("tyk-svc"),
			Namespace: pulumi.String(namespaceTykGateway.Name),
			Labels: pulumi.StringMap{
				"app": pulumi.String("tyk-gtw"),
			},
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       pulumi.Int(8080),
					TargetPort: pulumi.Int(8080),
				},
			},
			Selector: pulumi.StringMap{
				"app": pulumi.String("tyk-gtw"),
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = appsv1.NewDeployment(ctx, "redisDeployment", &appsv1.DeploymentArgs{
		ApiVersion: pulumi.String("apps/v1"),
		Kind:       pulumi.String("Deployment"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("redis"),
			Namespace: pulumi.String(namespaceTykGateway.Name),
			Labels: pulumi.StringMap{
				"app": pulumi.String("redis"),
			},
		},
		Spec: &appsv1.DeploymentSpecArgs{
			Selector: &metav1.LabelSelectorArgs{
				MatchLabels: pulumi.StringMap{
					"app": pulumi.String("redis"),
				},
			},
			Replicas: pulumi.Int(1),
			Template: &corev1.PodTemplateSpecArgs{
				Metadata: &metav1.ObjectMetaArgs{
					Labels: pulumi.StringMap{
						"app": pulumi.String("redis"),
					},
				},
				Spec: &corev1.PodSpecArgs{
					Containers: corev1.ContainerArray{
						&corev1.ContainerArgs{
							Name:  pulumi.String("master"),
							Image: pulumi.String("k8s.gcr.io/redis:e2e"),
							Resources: &corev1.ResourceRequirementsArgs{
								Limits: pulumi.StringMap{
									"memory": pulumi.String("512Mi"),
									"cpu":    pulumi.String("1"),
								},
								Requests: pulumi.StringMap{
									"memory": pulumi.String("256Mi"),
									"cpu":    pulumi.String("0.2"),
								},
							},
							Ports: corev1.ContainerPortArray{
								&corev1.ContainerPortArgs{
									ContainerPort: pulumi.Int(6379),
								},
							},
						},
					},
				},
			},
		},
	})
	if err != nil {
		return err
	}
	_, err = corev1.NewService(ctx, "redisService", &corev1.ServiceArgs{
		ApiVersion: pulumi.String("v1"),
		Kind:       pulumi.String("Service"),
		Metadata: &metav1.ObjectMetaArgs{
			Name:      pulumi.String("redis"),
			Namespace: pulumi.String(namespaceTykGateway.Name),
			Labels: pulumi.StringMap{
				"app": pulumi.String("redis"),
			},
		},
		Spec: &corev1.ServiceSpecArgs{
			Ports: corev1.ServicePortArray{
				&corev1.ServicePortArgs{
					Port:       pulumi.Int(6379),
					TargetPort: pulumi.Int(6379),
				},
			},
			Selector: pulumi.StringMap{
				"app": pulumi.String("redis"),
			},
		},
	})
	if err != nil {
		return err
	}
	return nil
}
