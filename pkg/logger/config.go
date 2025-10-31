// SPDX-FileCopyrightText: 2025 openstor contributors
// SPDX-FileCopyrightText: 2015-2025 MinIO, Inc.
// SPDX-License-Identifier: AGPL-3.0-or-later

package logger

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/openstor/console/pkg/logger/config"
	"github.com/openstor/console/pkg/logger/target/http"
	"github.com/openstor/pkg/v3/env"
)

// NewConfig - initialize new logger config.
func NewConfig() Config {
	cfg := Config{
		HTTP:         make(map[string]http.Config),
		AuditWebhook: make(map[string]http.Config),
	}

	return cfg
}

func lookupLoggerWebhookConfig() (Config, error) {
	cfg := NewConfig()
	envs := env.List(EnvLoggerWebhookEndpoint)
	var loggerTargets []string
	for _, k := range envs {
		target := strings.TrimPrefix(k, EnvLoggerWebhookEndpoint+config.Default)
		if target == EnvLoggerWebhookEndpoint {
			target = config.Default
		}
		loggerTargets = append(loggerTargets, target)
	}

	// Load HTTP logger from the environment if found
	for _, target := range loggerTargets {
		if v, ok := cfg.HTTP[target]; ok && v.Enabled {
			// This target is already enabled using the
			// legacy environment variables, ignore.
			continue
		}
		enableEnv := EnvLoggerWebhookEnable
		if target != config.Default {
			enableEnv = EnvLoggerWebhookEnable + config.Default + target
		}
		enable, err := config.ParseBool(env.Get(enableEnv, ""))
		if err != nil || !enable {
			continue
		}
		endpointEnv := EnvLoggerWebhookEndpoint
		if target != config.Default {
			endpointEnv = EnvLoggerWebhookEndpoint + config.Default + target
		}
		authTokenEnv := EnvLoggerWebhookAuthToken
		if target != config.Default {
			authTokenEnv = EnvLoggerWebhookAuthToken + config.Default + target
		}
		clientCertEnv := EnvLoggerWebhookClientCert
		if target != config.Default {
			clientCertEnv = EnvLoggerWebhookClientCert + config.Default + target
		}
		clientKeyEnv := EnvLoggerWebhookClientKey
		if target != config.Default {
			clientKeyEnv = EnvLoggerWebhookClientKey + config.Default + target
		}
		err = config.EnsureCertAndKey(env.Get(clientCertEnv, ""), env.Get(clientKeyEnv, ""))
		if err != nil {
			return cfg, err
		}
		queueSizeEnv := EnvLoggerWebhookQueueSize
		if target != config.Default {
			queueSizeEnv = EnvLoggerWebhookQueueSize + config.Default + target
		}
		queueSize, err := strconv.Atoi(env.Get(queueSizeEnv, "100000"))
		if err != nil {
			return cfg, err
		}
		if queueSize <= 0 {
			return cfg, errors.New("invalid queue_size value")
		}
		cfg.HTTP[target] = http.Config{
			Enabled:    true,
			Endpoint:   env.Get(endpointEnv, ""),
			AuthToken:  env.Get(authTokenEnv, ""),
			ClientCert: env.Get(clientCertEnv, ""),
			ClientKey:  env.Get(clientKeyEnv, ""),
			QueueSize:  queueSize,
		}
	}

	return cfg, nil
}

func lookupAuditWebhookConfig() (Config, error) {
	cfg := NewConfig()
	var loggerAuditTargets []string
	envs := env.List(EnvAuditWebhookEndpoint)
	for _, k := range envs {
		target := strings.TrimPrefix(k, EnvAuditWebhookEndpoint+config.Default)
		if target == EnvAuditWebhookEndpoint {
			target = config.Default
		}
		loggerAuditTargets = append(loggerAuditTargets, target)
	}

	for _, target := range loggerAuditTargets {
		if v, ok := cfg.AuditWebhook[target]; ok && v.Enabled {
			// This target is already enabled using the
			// legacy environment variables, ignore.
			continue
		}
		enableEnv := EnvAuditWebhookEnable
		if target != config.Default {
			enableEnv = EnvAuditWebhookEnable + config.Default + target
		}
		enable, err := config.ParseBool(env.Get(enableEnv, ""))
		if err != nil || !enable {
			continue
		}
		endpointEnv := EnvAuditWebhookEndpoint
		if target != config.Default {
			endpointEnv = EnvAuditWebhookEndpoint + config.Default + target
		}
		authTokenEnv := EnvAuditWebhookAuthToken
		if target != config.Default {
			authTokenEnv = EnvAuditWebhookAuthToken + config.Default + target
		}
		clientCertEnv := EnvAuditWebhookClientCert
		if target != config.Default {
			clientCertEnv = EnvAuditWebhookClientCert + config.Default + target
		}
		clientKeyEnv := EnvAuditWebhookClientKey
		if target != config.Default {
			clientKeyEnv = EnvAuditWebhookClientKey + config.Default + target
		}
		err = config.EnsureCertAndKey(env.Get(clientCertEnv, ""), env.Get(clientKeyEnv, ""))
		if err != nil {
			return cfg, err
		}
		queueSizeEnv := EnvAuditWebhookQueueSize
		if target != config.Default {
			queueSizeEnv = EnvAuditWebhookQueueSize + config.Default + target
		}
		queueSize, err := strconv.Atoi(env.Get(queueSizeEnv, "100000"))
		if err != nil {
			return cfg, err
		}
		if queueSize <= 0 {
			return cfg, errors.New("invalid queue_size value")
		}
		cfg.AuditWebhook[target] = http.Config{
			Enabled:    true,
			Endpoint:   env.Get(endpointEnv, ""),
			AuthToken:  env.Get(authTokenEnv, ""),
			ClientCert: env.Get(clientCertEnv, ""),
			ClientKey:  env.Get(clientKeyEnv, ""),
			QueueSize:  queueSize,
		}
	}

	return cfg, nil
}

// LookupConfigForSubSys - lookup logger config, override with ENVs if set, for the given sub-system
func LookupConfigForSubSys(subSys string) (cfg Config, err error) {
	switch subSys {
	case config.LoggerWebhookSubSys:
		if cfg, err = lookupLoggerWebhookConfig(); err != nil {
			return cfg, err
		}
	case config.AuditWebhookSubSys:
		if cfg, err = lookupAuditWebhookConfig(); err != nil {
			return cfg, err
		}
	}
	return cfg, nil
}

// GetGlobalDeploymentID :
func GetGlobalDeploymentID() string {
	if globalDeploymentID != "" {
		return globalDeploymentID
	}
	globalDeploymentID = env.Get(EnvGlobalDeploymentID, mustGetUUID())
	return globalDeploymentID
}

// mustGetUUID - get a random UUID.
func mustGetUUID() string {
	u, err := uuid.NewRandom()
	if err != nil {
		CriticalIf(GlobalContext, err)
	}
	return u.String()
}
