package weaviate

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/weaviate/weaviate-go-client/v5/weaviate/alias"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/auth"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/backup"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/batch"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/classifications"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/cluster"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/connection"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/contextionary"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/data"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/db"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/graphql"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/grpc"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/misc"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/rbac"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/schema"
	"github.com/weaviate/weaviate-go-client/v5/weaviate/users"
)

const defaultTimeout = 60 * time.Second

// Config of the client endpoint
type Config struct {
	// Host of the weaviate instance; this is a mandatory field.
	Host string
	// Scheme of the weaviate instance; this is a mandatory field.
	Scheme string

	// ConnectionClient that will be used to execute http requests to the weaviate instance.
	//  If omitted a default will be used. The default is not able to handle authenticated requests.
	//
	//  To connect with an authenticated weaviate consider using the client from the golang.org/x/oauth2 module.
	// Either this option or AuthConfig can be used
	ConnectionClient *http.Client

	// Configuration for authentication. Either this option or ConnectionClient can be used
	AuthConfig auth.Config

	// Headers added for every request
	Headers map[string]string

	// How long the client should wait for Weaviate to start up
	StartupTimeout time.Duration

	// gRPC configuration
	GrpcConfig *grpc.Config

	// Client connection timeout, defaults to 60s
	Timeout time.Duration
}

func (c Config) getTimeout() time.Duration {
	if c.Timeout == 0 {
		return defaultTimeout
	}
	return c.Timeout
}

// Deprecated: This function is unable to wait for Weaviate to start. Use NewClient() instead and add auth.Config to
// weaviate.Config
func NewConfig(host string, scheme string, authConfig auth.Config,
	headers map[string]string, grpcConfig ...*grpc.Config,
) (*Config, error) {
	var client *http.Client
	var err error
	var additionalHeaders map[string]string
	if authConfig != nil {
		tmpCon := connection.NewConnection(scheme, host, nil, defaultTimeout, headers)
		client, additionalHeaders, err = authConfig.GetAuthInfo(tmpCon)
		if err != nil {
			return nil, err
		}
		if headers == nil {
			headers = map[string]string{}
		}
		for k, v := range additionalHeaders {
			headers[k] = v
		}
	}
	var grpcConf *grpc.Config
	if len(grpcConfig) > 0 && grpcConfig[0] != nil {
		grpcConf = grpcConfig[0]
	}

	return &Config{Host: host, Scheme: scheme, Headers: headers, ConnectionClient: client, GrpcConfig: grpcConf}, nil
}

// Client implementing the weaviate API
// Every function represents one API group of weaviate and provides a set of functions and builders to interact with them.
//
// The client uses the original data models as provided by weaviate itself.
// All these models are provided in the sub module "github.com/weaviate/weaviate/entities/models"
type Client struct {
	connection      *connection.Connection
	grpcClient      *connection.GrpcClient
	misc            *misc.API
	schema          *schema.API
	alias           *alias.API
	data            *data.API
	batch           *batch.API
	c11y            *contextionary.API
	classifications *classifications.API
	backup          *backup.API
	graphQL         *graphql.API
	cluster         *cluster.API
	roles           *rbac.API
	users           *users.API
	experimental    *experimental
}

// experimental contains all experimental client features
type experimental struct {
	grpcClient *connection.GrpcClient
}

// Experimental Search gRPC API group
func (e *experimental) Search() *graphql.Search {
	return graphql.NewSearch(e.grpcClient)
}

func NewClient(config Config) (*Client, error) {
	if config.AuthConfig != nil && config.ConnectionClient != nil {
		return nil, errors.New("only AuthConfig or ConnectionClient can be given in the config")
	}

	// if an authentication config is given, we first need to create a temporary connection to fetch some OIDC
	// infos from Weaviate. This connection is then replaced by the "real" connection
	if config.AuthConfig != nil {
		tmpCon := connection.NewConnection(config.Scheme, config.Host, nil, config.getTimeout(), config.Headers)
		err := tmpCon.WaitForWeaviate(config.StartupTimeout)
		if err != nil {
			return nil, err
		}
		connectionClient, additionalHeaders, err := config.AuthConfig.GetAuthInfo(tmpCon)
		if err != nil {
			return nil, err
		}
		config.ConnectionClient = connectionClient
		if config.Headers == nil {
			config.Headers = map[string]string{}
		}
		for k, v := range additionalHeaders {
			config.Headers[k] = v
		}
		if isWeaviateDomain(config.Host) && config.AuthConfig.ApiKey() != nil {
			config.Headers["X-Weaviate-Api-Key"] = *config.AuthConfig.ApiKey()
			config.Headers["X-Weaviate-Cluster-URL"] = "https://" + config.Host
		}
	}

	con := connection.NewConnection(config.Scheme, config.Host, config.ConnectionClient, config.getTimeout(), config.Headers)

	if err := con.WaitForWeaviate(config.StartupTimeout); err != nil {
		return nil, err
	}

	// some endpoints now require a className namespace.
	// to determine if this new convention is to be used,
	// we must check the weaviate server version
	getVersionFn := func() string {
		meta, err := misc.New(con, nil).MetaGetter().Do(context.Background())
		if err == nil {
			return meta.Version
		}
		return ""
	}

	dbVersionProvider := db.NewVersionProvider(getVersionFn)
	dbVersionSupport := db.NewDBVersionSupport(dbVersionProvider)
	grpcVersionSupport := db.NewGRPCVersionSupport(dbVersionProvider)

	grpcClient, err := createGrpcClient(config, grpcVersionSupport)
	if err != nil {
		return nil, fmt.Errorf("create weaviate client: %w", err)
	}

	client := &Client{
		connection:      con,
		grpcClient:      grpcClient,
		misc:            misc.New(con, dbVersionProvider),
		schema:          schema.New(con),
		alias:           alias.New(con),
		c11y:            contextionary.New(con),
		classifications: classifications.New(con),
		graphQL:         graphql.New(con),
		data:            data.New(con, dbVersionSupport),
		batch:           batch.New(con, grpcClient, dbVersionSupport),
		backup:          backup.New(con),
		cluster:         cluster.New(con),
		roles:           rbac.New(con),
		users:           users.New(con),
		experimental:    &experimental{grpcClient: grpcClient},
	}

	return client, nil
}

// New client from config
// Deprecated: Use NewClient() instead, which returns an error instead of panicking
func New(config Config) *Client {
	if client, err := NewClient(config); err == nil {
		return client
	}
	con := connection.NewConnection(config.Scheme, config.Host, config.ConnectionClient, config.getTimeout(), config.Headers)

	// some endpoints now require a className namespace.
	// to determine if this new convention is to be used,
	// we must check the weaviate server version
	getVersionFn := func() string {
		meta, err := misc.New(con, nil).MetaGetter().Do(context.Background())
		if err == nil {
			return meta.Version
		}
		return ""
	}

	dbVersionProvider := db.NewVersionProvider(getVersionFn)
	dbVersionSupport := db.NewDBVersionSupport(dbVersionProvider)
	gRPCVersionSupport := db.NewGRPCVersionSupport(dbVersionProvider)

	grpcClient, err := createGrpcClient(config, gRPCVersionSupport)
	if err != nil {
		panic(err)
	}

	client := &Client{
		connection:      con,
		grpcClient:      grpcClient,
		misc:            misc.New(con, dbVersionProvider),
		schema:          schema.New(con),
		c11y:            contextionary.New(con),
		classifications: classifications.New(con),
		graphQL:         graphql.New(con),
		data:            data.New(con, dbVersionSupport),
		batch:           batch.New(con, grpcClient, dbVersionSupport),
		backup:          backup.New(con),
		cluster:         cluster.New(con),
		roles:           rbac.New(con),
		users:           users.New(con),
		experimental:    &experimental{grpcClient: grpcClient},
	}

	return client
}

// Waits for Weaviate to start.
func (c *Client) WaitForWeavaite(startupTimeout time.Duration) error {
	return c.connection.WaitForWeaviate(startupTimeout)
}

// Misc collection group for .well_known and root level API commands
func (c *Client) Misc() *misc.API {
	return c.misc
}

// Schema API group
func (c *Client) Schema() *schema.API {
	return c.schema
}

// Alias API group
func (c *Client) Alias() *alias.API {
	return c.alias
}

// Data API group including both things and actions
func (c *Client) Data() *data.API {
	return c.data
}

// Batch loading API group
func (c *Client) Batch() *batch.API {
	return c.batch
}

// C11y (contextionary) API group
func (c *Client) C11y() *contextionary.API {
	return c.c11y
}

// Classifications API group
func (c *Client) Classifications() *classifications.API {
	return c.classifications
}

// GraphQL API group
func (c *Client) GraphQL() *graphql.API {
	return c.graphQL
}

// Backup API group
func (c *Client) Backup() *backup.API {
	return c.backup
}

// Cluster API group
func (c *Client) Cluster() *cluster.API {
	return c.cluster
}

func (c *Client) Roles() *rbac.API {
	return c.roles
}

func (c *Client) Users() *users.API {
	return c.users
}

// Experimental API group
func (c *Client) Experimental() *experimental {
	return c.experimental
}

func createGrpcClient(config Config, gRPCVersionSupport *db.GRPCVersionSupport) (*connection.GrpcClient, error) {
	if config.GrpcConfig != nil {
		return connection.NewGrpcClient(config.GrpcConfig.Host, config.GrpcConfig.Secured, config.Headers, gRPCVersionSupport, config.getTimeout(), config.StartupTimeout)
	}
	return nil, nil
}

func isWeaviateDomain(url string) bool {
	lower := strings.ToLower(url)
	return strings.Contains(lower, "weaviate.io") || strings.Contains(lower, "semi.technology") || strings.Contains(lower, "weaviate.cloud")
}
