package typedefs

import genesis_v1 "go.taskfleet.io/grpc/gen/go/genesis/v1"

// Provider describes a cloud provider where compute instances may be created.
type Provider string

const (
	// ProviderAmazonWebServices describes the AWS cloud provider.
	ProviderAmazonWebServices Provider = "amazon-web-services"
	// ProviderGoogleCloudPlatform describes the GCP cloud provider.
	ProviderGoogleCloudPlatform Provider = "google-cloud-platform"
)

// CloudProviderUnmarshalProto returns the internal representation for a hermes cloud provider.
func CloudProviderUnmarshalProto(message genesis_v1.Provider) Provider {
	switch message {
	case genesis_v1.Provider_PROVIDER_AMAZON_WEB_SERVICES:
		return ProviderAmazonWebServices
	case genesis_v1.Provider_PROVIDER_GOOGLE_CLOUD_PLATFORM:
		return ProviderGoogleCloudPlatform
	default:
		panic("unknown cloud provider")
	}
}

// MarshalProto returns the hermes enum of the cloud provider.
func (p Provider) MarshalProto() genesis_v1.Provider {
	switch p {
	case ProviderAmazonWebServices:
		return genesis_v1.Provider_PROVIDER_AMAZON_WEB_SERVICES
	case ProviderGoogleCloudPlatform:
		return genesis_v1.Provider_PROVIDER_GOOGLE_CLOUD_PLATFORM
	default:
		panic("unknown cloud provider")
	}
}
