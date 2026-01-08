package main

import (
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		projectCfg := config.New(ctx, "T-pot")
		sshKey := projectCfg.Require("sshPublicKey")

		ociCfg := config.New(ctx, "oci")
		rootCompartmentId := ociCfg.Require("tenancyOcid")

		subnet, err := CreateNetwork(ctx, rootCompartmentId)
		if err != nil {
			return err
		}

		err = CreateHoneypotServer(ctx, subnet, rootCompartmentId, sshKey)
		if err != nil {
			return err
		}
		return nil
	})
}
