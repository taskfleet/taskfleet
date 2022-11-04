package gcptestutils

import computepb "google.golang.org/genproto/googleapis/cloud/compute/v1"

type fakeNetworkServer struct {
	computepb.UnimplementedNetworkServiceServer
}

type fakeSubnetworkServer struct{}
