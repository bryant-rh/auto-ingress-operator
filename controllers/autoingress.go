package controllers

import (
	"fmt"
	"sync"

	appsv1 "github.com/bryant-rh/auto-ingress-operator/api/v1"
)

type AutoIngressContainer struct {
	mu  sync.Mutex
	set map[string]appsv1.AutoIngress
}

func NewAutoIngressContainer() *AutoIngressContainer {
	return &AutoIngressContainer{
		set: make(map[string]appsv1.AutoIngress),
	}
}

func (c *AutoIngressContainer) Add(ing appsv1.AutoIngress) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.key(ing)
	c.set[key] = ing

}

func (c *AutoIngressContainer) Remove(ing appsv1.AutoIngress) {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := c.key(ing)
	delete(c.set, key)
}

func (c *AutoIngressContainer) List() []appsv1.AutoIngress {
	list := make([]appsv1.AutoIngress, 0)
	for _, v := range c.set {
		list = append(list, v)
	}

	return list
}

func (c *AutoIngressContainer) key(ing appsv1.AutoIngress) string {
	return fmt.Sprintf("%s-%s", ing.Namespace, ing.Name)
}

var autoIngressSet *AutoIngressContainer

func init() {
	autoIngressSet = NewAutoIngressContainer()
}
