package main

import (
	"github.com/pulumi/pulumi-oci/sdk/v3/go/oci/core"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		targetPorts := []HoneyPort{
			{22, "6", "Cowrie SSH Honeypot"},
			{23, "6", "Cowrie Telnet Honeypot"},
			{80, "6", "Web/Log4Shell Honeypot"},
			{443, "6", "Secure Web Honeypot"},
			{445, "6", "Dionaea SMB/Ransomware Honeypot"},
			{5555, "6", "ADBHoney Android Malware"},
			{6379, "6", "Redis Database Exploits"},
			{5900, "6", "VNC Remote Desktop Exploits"},
			{5060, "17", "SentryPeer VoIP Fraud (UDP)"},
		}

		vcn, err := core.NewVcn(ctx, "honeypot-vcn", &core.VcnArgs{
			CidrBlocks: pulumi.StringArray{
			pulumi.String("10.0.0.0/16"),
			}
		})
	
		igw, err := core.NewInternetGateway(ctx, "honeypot-igw", &core.InternetGatewayArgs{
			VcnId: vcn.ID(),
		})
	
		rt, err := core.NewRouteTable(ctx, "honeypot-rt", &core.RouteTableArgs{
			VcnId: vcn.ID(),
			RouteRules: core.RouteTableRouteRuleArray{
				&core.RouteTableRouteRuleArgs{
					Destination: pulumi.String("0.0.0.0/0"),
					NetworkEntityId: igw.ID(),
				}
			}
		})
		
		var ingressRules core.SecurityListIngressSecurityRuleArray

		for _, hp := range targetPorts {
			ruleArgs := &core.SecurityListIngressSecurityRuleArgs{
				Protocol:    pulumi.String(hp.Protocol),
				Source:      pulumi.String("0.0.0.0/0"), 
				Description: pulumi.String(hp.Description),
				Stateless:   pulumi.Bool(false),
			}

			if hp.Protocol == "6" {
				ruleArgs.TcpOptions = &core.SecurityListIngressSecurityRuleTcpOptionsArgs{
					Max: pulumi.Int(hp.Port),
					Min: pulumi.Int(hp.Port),
				}
			} else if hp.Protocol == "17" {
				ruleArgs.UdpOptions = &core.SecurityListIngressSecurityRuleUdpOptionsArgs{
					Max: pulumi.Int(hp.Port),
					Min: pulumi.Int(hp.Port),
				}
			}

			ingressRules = append(ingressRules, ruleArgs)
		}

		_, err := core.NewSecurityList(ctx, "honeypot-firewall", &core.SecurityListArgs{
			DisplayName:          pulumi.String("honeypot-security-list"),
			VcnId:                vcn.ID(),
			IngressSecurityRules: ingressRules,
			EgressSecurityRules: core.SecurityListEgressSecurityRuleArray{
				&core.SecurityListEgressSecurityRuleArgs{
					Destination: pulumi.String("0.0.0.0/0"),
					Protocol:    pulumi.String("all"), /
				},
			},
		})
	}
	if err != nil {
		return err
	} 
	return nil
)}	
