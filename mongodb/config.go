/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package mongodb

import (
	"errors"
	"net"
	"strconv"
	"time"

	"github.com/spf13/viper"

	"github.com/apecloud/mongodb_plugin/constant"
)

const (
	host             = "host"
	username         = "username"
	password         = "password"
	server           = "server"
	databaseName     = "databaseName"
	operationTimeout = "operationTimeout"
	params           = "params"
	adminDatabase    = "admin"

	defaultTimeout = 5 * time.Second
	defaultDBPort  = 27017

	EnvRootUser     = "MONGODB_ROOT_USER"
	EnvRootPassword = "MONGODB_ROOT_PASSWORD"
)

type Config struct {
	Hosts            []string
	Username         string
	Password         string
	ReplSetName      string
	DatabaseName     string
	Params           string
	Direct           bool
	OperationTimeout time.Duration
}

var config *Config

func NewConfig(properties map[string]string) (*Config, error) {
	config = &Config{
		Hosts:            []string{"localhost:27017"},
		Direct:           true,
		Username:         "root",
		OperationTimeout: defaultTimeout,
	}

	if viper.IsSet("KB_SERVICE_PORT") {
		config.Hosts = []string{"localhost:" + viper.GetString("KB_SERVICE_PORT")}
	}

	if viper.IsSet("SERVICE_PORT") {
		config.Hosts = []string{"localhost:" + viper.GetString("SERVICE_PORT")}
	}

	if val, ok := properties[host]; ok && val != "" {
		config.Hosts = []string{val}
	}

	if val, ok := properties[username]; ok && val != "" {
		config.Username = val
	}

	if val, ok := properties[password]; ok && val != "" {
		config.Password = val
	}

	if viper.IsSet(EnvRootUser) {
		config.Username = viper.GetString(EnvRootUser)
	}

	if viper.IsSet(EnvRootPassword) {
		config.Password = viper.GetString(EnvRootPassword)
	}

	if viper.IsSet(constant.KBEnvClusterCompName) {
		config.ReplSetName = viper.GetString(constant.KBEnvClusterCompName)
	}

	config.DatabaseName = adminDatabase
	if val, ok := properties[databaseName]; ok && val != "" {
		config.DatabaseName = val
	}

	if val, ok := properties[params]; ok && val != "" {
		config.Params = val
	}

	var err error
	if val, ok := properties[operationTimeout]; ok && val != "" {
		config.OperationTimeout, err = time.ParseDuration(val)
		if err != nil {
			return nil, errors.New("incorrect operationTimeout field from metadata")
		}
	}

	return config, nil
}

func (config *Config) GetDBPort() int {
	_, portStr, err := net.SplitHostPort(config.Hosts[0])
	if err != nil {
		return defaultDBPort
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return defaultDBPort
	}

	return port
}

func (config *Config) DeepCopy() *Config {
	newConf := *config
	newConf.Hosts = make([]string, len(config.Hosts))
	copy(newConf.Hosts, config.Hosts)
	return &newConf
}

func GetConfig() *Config {
	return config
}
