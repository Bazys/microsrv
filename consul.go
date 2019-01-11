package config

import (
	"errors"

	"github.com/hashicorp/consul/api"
	"github.com/labstack/echo"
)

//ServiceCatalog struct
type ServiceCatalog struct {
	Services map[string]*api.AgentService
}

// Fill func
func (s *ServiceCatalog) Fill(address string) error {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = "127.0.0.1:8500"
	consulClient, err := api.NewClient(consulConfig)
	if err != nil {
		return err
	}
	query := &api.QueryOptions{}
	catalog := consulClient.Catalog()
	nodes, _, err := catalog.Nodes(query)
	if err != nil {
		return err
	}
	nodeName := nodes[0].Node
	node, _, err := catalog.Node(nodeName, query)
	if err != nil {
		return err
	}
	delete(node.Services, "consul")
	s.Services = node.Services
	return nil
}

// ConfigMiddleware func
func (s ServiceCatalog) ServiceCatalogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("serviceCatalog", s)
			return next(c)
		}
	}
}

// ConfigFromContext func
func (s *ServiceCatalog) ServiceCatalogFromContext(c echo.Context) error {
	s, ok := c.Get("serviceCatalog").(*ServiceCatalog)
	if !ok {
		return errors.New("No serviceCatalog in context")
	}
	return nil
}
