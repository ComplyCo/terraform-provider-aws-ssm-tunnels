package ssmtunnels

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	pluginSession "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session"
	_ "github.com/aws/session-manager-plugin/src/sessionmanagerplugin/session/portsession"
)

type RemoteTunnelConfig struct {
	Client     *ssm.Client
	Target     string
	Region     string
	RemoteHost string
	RemotePort int
	LocalPort  int
}

func StartRemoteTunnel(ctx context.Context, cfg RemoteTunnelConfig) error {
	if cfg.Target == "" {
		return fmt.Errorf("target must be set")
	}
	if cfg.Region == "" {
		return fmt.Errorf("region must be set")
	}
	if cfg.RemoteHost == "" {
		return fmt.Errorf("remoteHost must be set")
	}
	if cfg.RemotePort == 0 {
		return fmt.Errorf("remotePort must be set")
	}
	if cfg.LocalPort == 0 {
		return fmt.Errorf("localPort must be set")
	}

	startSessionInput := ssm.StartSessionInput{
		Target:       &cfg.Target,
		DocumentName: aws.String("AWS-StartPortForwardingSessionToRemoteHost"),
		Parameters: map[string][]string{
			"host": {
				cfg.RemoteHost,
			},
			"portNumber": {
				strconv.Itoa(cfg.RemotePort),
			},
			"localPortNumber": {
				strconv.Itoa(cfg.LocalPort),
			},
		},
	}

	startSessionOutput, err := cfg.Client.StartSession(ctx, &startSessionInput)
	if err != nil {
		return err
	}

	startSessionOuputJson, err := json.Marshal(startSessionOutput)
	if err != nil {
		return err
	}

	// TODO: Add a way to terminate the session
	// cfg.Client.TerminateSession()

	args := []string{
		"session-manager-plugin",
		string(startSessionOuputJson),
		cfg.Region,
		"StartSession",
		"",
		fmt.Sprintf("{\"Target\": \"%s\"}", cfg.Target),
		fmt.Sprintf("https://ssm.%s.amazonaws.com", cfg.Region),
	}

	// TODO: Run this in a cancelable goroutine

	pluginSession.ValidateInputAndStartSession(args, os.Stdout)

	return nil
}
