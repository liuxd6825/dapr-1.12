/*
Copyright 2021 The Dapr Authors
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"flag"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/dapr/kit/logger"
	"github.com/liuxd6825/dapr/pkg/config"
	"github.com/liuxd6825/dapr/pkg/config/protocol"
	"github.com/liuxd6825/dapr/pkg/cors"
	"github.com/liuxd6825/dapr/pkg/metrics"
	"github.com/liuxd6825/dapr/pkg/modes"
	"github.com/liuxd6825/dapr/pkg/runtime"
	"github.com/liuxd6825/dapr/pkg/security/consts"
)

type Options struct {
	AppID                        string
	ComponentsPath               string
	ControlPlaneAddress          string
	ControlPlaneTrustDomain      string
	ControlPlaneNamespace        string
	SentryAddress                string
	TrustAnchors                 []byte
	AllowedOrigins               string
	EnableProfiling              bool
	AppMaxConcurrency            int
	EnableMTLS                   bool
	AppSSL                       bool
	DaprHTTPMaxRequestSize       int
	ResourcesPath                stringSliceFlag
	AppProtocol                  string
	EnableAPILogging             *bool
	RuntimeVersion               bool
	BuildInfo                    bool
	WaitCommand                  bool
	DaprHTTPPort                 string
	DaprAPIGRPCPort              string
	ProfilePort                  string
	DaprInternalGRPCPort         string
	DaprPublicPort               string
	AppPort                      string
	DaprGracefulShutdownSeconds  int
	PlacementServiceHostAddr     string
	DaprAPIListenAddresses       string
	AppHealthProbeInterval       int
	AppHealthProbeTimeout        int
	AppHealthThreshold           int
	EnableAppHealthCheck         bool
	Mode                         string
	Config                       stringSliceFlag
	UnixDomainSocket             string
	DaprHTTPReadBufferSize       int
	DisableBuiltinK8sSecretStore bool
	AppHealthCheckPath           string
	AppChannelAddress            string
	Logger                       logger.Options
	Metrics                      *metrics.Options
	LogFile                      string
	LogOutputType                string
}

func New(args []string) *Options {
	opts := Options{
		EnableAPILogging: new(bool),
	}

	flag.StringVar(&opts.Mode, "mode", string(modes.StandaloneMode), "Runtime mode for Dapr")
	flag.StringVar(&opts.DaprHTTPPort, "dapr-http-port", strconv.Itoa(runtime.DefaultDaprHTTPPort), "HTTP port for Dapr API to listen on")
	flag.StringVar(&opts.DaprAPIListenAddresses, "dapr-listen-addresses", runtime.DefaultAPIListenAddress, "One or more addresses for the Dapr API to listen on, CSV limited")
	flag.StringVar(&opts.DaprPublicPort, "dapr-public-port", "", "Public port for Dapr Health and Metadata to listen on")
	flag.StringVar(&opts.DaprAPIGRPCPort, "dapr-grpc-port", strconv.Itoa(runtime.DefaultDaprAPIGRPCPort), "gRPC port for the Dapr API to listen on")
	flag.StringVar(&opts.DaprInternalGRPCPort, "dapr-internal-grpc-port", "", "gRPC port for the Dapr Internal API to listen on")
	flag.StringVar(&opts.AppPort, "app-port", "", "The port the application is listening on")
	flag.StringVar(&opts.ProfilePort, "profile-port", strconv.Itoa(runtime.DefaultProfilePort), "The port for the profile server")
	flag.StringVar(&opts.AppProtocol, "app-protocol", string(protocol.HTTPProtocol), "Protocol for the application: grpc, grpcs, http, https, h2c")
	flag.StringVar(&opts.ComponentsPath, "components-path", "", "Alias for --resources-path [Deprecated, use --resources-path]")
	flag.Var(&opts.ResourcesPath, "resources-path", "Path for resources directory. If not specified, no resources will be loaded. Can be passed multiple times")
	flag.Var(&opts.Config, "config", "Path to config file, or name of a configuration object. In standalone mode, can be passed multiple times")
	flag.StringVar(&opts.AppID, "app-id", "", "A unique ID for Dapr. Used for Service Discovery and state")
	flag.StringVar(&opts.ControlPlaneAddress, "control-plane-address", "", "Address for a Dapr control plane")
	flag.StringVar(&opts.SentryAddress, "sentry-address", "", "Address for the Sentry CA service")
	flag.StringVar(&opts.ControlPlaneTrustDomain, "control-plane-trust-domain", "localhost", "Trust domain of the Dapr control plane")
	flag.StringVar(&opts.ControlPlaneNamespace, "control-plane-namespace", "default", "Namespace of the Dapr control plane")
	flag.StringVar(&opts.PlacementServiceHostAddr, "placement-host-address", "", "Addresses for Dapr Actor Placement servers")
	flag.StringVar(&opts.AllowedOrigins, "allowed-origins", cors.DefaultAllowedOrigins, "Allowed HTTP origins")
	flag.BoolVar(&opts.EnableProfiling, "enable-profiling", false, "Enable profiling")
	flag.BoolVar(&opts.RuntimeVersion, "version", false, "Prints the runtime version")
	flag.BoolVar(&opts.BuildInfo, "build-info", false, "Prints the build info")
	flag.BoolVar(&opts.WaitCommand, "wait", false, "wait for Dapr outbound ready")
	flag.IntVar(&opts.AppMaxConcurrency, "app-max-concurrency", -1, "Controls the concurrency level when forwarding requests to user code; set to -1 for no limits")
	flag.BoolVar(&opts.EnableMTLS, "enable-mtls", false, "Enables automatic mTLS for daprd to daprd communication channels")
	flag.BoolVar(&opts.AppSSL, "app-ssl", false, "Sets the URI scheme of the app to https and attempts a TLS connection [Deprecated, use '--app-protocol https|grpcs']")
	flag.IntVar(&opts.DaprHTTPMaxRequestSize, "dapr-http-max-request-size", runtime.DefaultMaxRequestBodySize, "Increasing max size of request body in MB to handle uploading of big files")
	flag.StringVar(&opts.UnixDomainSocket, "unix-domain-socket", "", "Path to a unix domain socket dir mount. If specified, Dapr API servers will use Unix Domain Sockets")
	flag.IntVar(&opts.DaprHTTPReadBufferSize, "dapr-http-read-buffer-size", runtime.DefaultReadBufferSize, "Increasing max size of read buffer in KB to handle sending multi-KB headers")
	flag.IntVar(&opts.DaprGracefulShutdownSeconds, "dapr-graceful-shutdown-seconds", int(runtime.DefaultGracefulShutdownDuration/time.Second), "Graceful shutdown time in seconds")
	flag.BoolVar(opts.EnableAPILogging, "enable-api-logging", false, "Enable API logging for API calls")
	flag.BoolVar(&opts.DisableBuiltinK8sSecretStore, "disable-builtin-k8s-secret-store", false, "Disable the built-in Kubernetes Secret Store")
	flag.BoolVar(&opts.EnableAppHealthCheck, "enable-app-health-check", false, "Enable health checks for the application using the protocol defined with app-protocol")
	flag.StringVar(&opts.AppHealthCheckPath, "app-health-check-path", runtime.DefaultAppHealthCheckPath, "Path used for health checks; HTTP only")
	flag.IntVar(&opts.AppHealthProbeInterval, "app-health-probe-interval", int(config.AppHealthConfigDefaultProbeInterval/time.Second), "Interval to probe for the health of the app in seconds")
	flag.IntVar(&opts.AppHealthProbeTimeout, "app-health-probe-timeout", int(config.AppHealthConfigDefaultProbeTimeout/time.Millisecond), "Timeout for app health probes in milliseconds")
	flag.IntVar(&opts.AppHealthThreshold, "app-health-threshold", int(config.AppHealthConfigDefaultThreshold), "Number of consecutive failures for the app to be considered unhealthy")
	flag.StringVar(&opts.AppChannelAddress, "app-channel-address", runtime.DefaultChannelAddress, "The network address the application listens on")
	flag.StringVar(&opts.LogFile, "log-file", "", "Save log file name ")
	flag.StringVar(&opts.LogOutputType, "log-output-type", "console", "choose output type (console file, all)")

	opts.Logger = logger.DefaultOptions()
	opts.Logger.AttachCmdFlags(flag.StringVar, flag.BoolVar)

	opts.Metrics = metrics.DefaultMetricOptions()
	opts.Metrics.AttachCmdFlags(flag.StringVar, flag.BoolVar)

	// Ignore errors; CommandLine is set for ExitOnError.
	flag.CommandLine.Parse(args)

	opts.TrustAnchors = []byte(os.Getenv(consts.TrustAnchorsEnvVar))

	// flag.Parse() will always set a value to "enableAPILogging", and it will be false whether it's explicitly set to false or unset
	// For this flag, we need the third state (unset) so we need to do a bit more work here to check if it's unset, then mark "enableAPILogging" as nil
	// It's not the prettiest approach, but…
	if !isFlagPassed("enable-api-logging") {
		opts.EnableAPILogging = nil
	}

	if !isFlagPassed("control-plane-namespace") {
		ns, ok := os.LookupEnv(consts.ControlPlaneNamespaceEnvVar)
		if ok {
			opts.ControlPlaneNamespace = ns
		}
	}

	if !isFlagPassed("control-plane-trust-domain") {
		td, ok := os.LookupEnv(consts.ControlPlaneTrustDomainEnvVar)
		if ok {
			opts.ControlPlaneTrustDomain = td
		}
	}

	initOptions(&opts, args)

	return &opts
}

// initOptions 1:解决路径参数取不到的问题
func initOptions(opts *Options, args []string) {
	// 1:解决路径参数取不到的问题 flag包的bug
	opts.LogFile = getArgValue(opts.LogFile, args, "-log-file")
	opts.ComponentsPath = getArgValue(opts.ComponentsPath, args, "-components-path")
	opts.ResourcesPath = getArgSlice(opts.ResourcesPath, args, "-resources-path")
	opts.Config = getArgSlice(opts.Config, args, "-config")

	// 2:将路径转为绝对路径
	opts.LogFile = getAbsPath(opts.LogFile)
	opts.ComponentsPath = getAbsPath(opts.ComponentsPath)
	opts.Config = sliceAbsPath(opts.Config)
	opts.ResourcesPath = sliceAbsPath(opts.ResourcesPath)
}

func sliceAbsPath(list stringSliceFlag) stringSliceFlag {
	var res stringSliceFlag
	for _, value := range list {
		_ = res.Set(getAbsPath(value))
	}
	return res
}

func getArgSlice(defVal stringSliceFlag, args []string, argName string) stringSliceFlag {
	if len(defVal) == 0 {
		var s stringSliceFlag
		val := getArgValue("", args, argName)
		if val != "" {
			_ = s.Set(val)
		}
		return s
	}
	return defVal
}

func getAbsPath(pathOrFile string) string {
	s := strings.Trim(pathOrFile, " ")
	if s != "" {
		s, _ = filepath.Abs(s)
	}
	return s
}

func getArgValue(defValue string, args []string, argName string) string {
	if defValue != "" {
		return defValue
	}
	argName = strings.ToLower(argName)
	for i, s := range args {
		if s == argName || s == argName+"=" {
			return args[i+1]
		} else if strings.Contains(s, "=") {
			idx := strings.Index(s, "=")
			name := s[:idx]
			if name == argName {
				return strings.Trim(s[idx+1:], "")
			}
		}
	}
	return ""
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// Flag type. Allows passing a flag multiple times to get a slice of strings.
// It implements the flag.Value interface.
type stringSliceFlag []string

// String formats the flag value.
func (f stringSliceFlag) String() string {
	return strings.Join(f, ",")
}

// Set the flag value.
func (f *stringSliceFlag) Set(value string) error {
	if value == "" {
		return nil
	}
	*f = append(*f, value)
	return nil
}
