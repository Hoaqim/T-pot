package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/core"
	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/identity"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateHoneypotServer(ctx *pulumi.Context, subnet *core.Subnet, compartmentId string, sshKey string) error {

	ads, err := identity.GetAvailabilityDomains(ctx, &identity.GetAvailabilityDomainsArgs{
		CompartmentId: compartmentId,
	}, nil)
	if err != nil {
		return err
	}
	if len(ads.AvailabilityDomains) == 0 {
		return fmt.Errorf("no availability domains found in this region")
	}
	primaryAD := ads.AvailabilityDomains[1].Name

	images, err := core.GetImages(ctx, &core.GetImagesArgs{
		CompartmentId:   compartmentId,
		OperatingSystem: pulumi.StringRef("Canonical Ubuntu"),
		Shape:           pulumi.StringRef("VM.Standard.A1.Flex"),
		SortBy:          pulumi.StringRef("TIMECREATED"),
		SortOrder:       pulumi.StringRef("DESC"),
	}, nil)
	if err != nil {
		return err
	}
	if len(images.Images) == 0 {
		return fmt.Errorf("no Ubuntu A1 images found")
	}
	imageId := images.Images[0].Id

	scriptContent, err := os.ReadFile("startup.sh")
	if err != nil {
		return err
	}
	encodedScript := base64.StdEncoding.EncodeToString(scriptContent)
	_, err = core.NewInstance(ctx, "honeypot-server", &core.InstanceArgs{
		AvailabilityDomain: pulumi.String(primaryAD),
		CompartmentId:      pulumi.String(compartmentId),
		DisplayName:        pulumi.String("tpot-honeypot-01"),
		Shape:              pulumi.String("VM.Standard.A1.Flex"),

		ShapeConfig: &core.InstanceShapeConfigArgs{
			Ocpus:       pulumi.Float64(4),
			MemoryInGbs: pulumi.Float64(24),
		},

		CreateVnicDetails: &core.InstanceCreateVnicDetailsArgs{
			SubnetId:       subnet.ID(),
			AssignPublicIp: pulumi.String("true"),
			DisplayName:    pulumi.String("honeypot-vnic"),
		},

		SourceDetails: &core.InstanceSourceDetailsArgs{
			SourceType:          pulumi.String("image"),
			SourceId:            pulumi.String(imageId),
			BootVolumeSizeInGbs: pulumi.String("200"),
		},

		Metadata: pulumi.StringMap{
			"ssh_authorized_keys": pulumi.String(sshKey),
			"user_data":           pulumi.String(encodedScript),
		},
	})
	if err != nil {
		return err
	}

	return nil
}
