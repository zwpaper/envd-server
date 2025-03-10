// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package server

import (
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/tensorchord/envd-server/api/types"
)

func (s *Server) environmentCreate(c *gin.Context) {
	var req types.EnvironmentCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(500, err)
		return
	}
	expectedPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.IdentityToken,
			Namespace: "default",
			Labels: map[string]string{
				"name": "test",
			},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "envd",
					Image: req.Image,
					Ports: []v1.ContainerPort{
						{
							Name:          "ssh",
							ContainerPort: 2222,
						},
					},
				},
			},
		},
	}

	created, err := s.client.CoreV1().Pods(
		"default").Create(c, &expectedPod, metav1.CreateOptions{})
	if err != nil {
		c.JSON(500, err)
		return
	}
	expectedService := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.IdentityToken,
			Namespace: "default",
			Labels: map[string]string{
				"name": req.IdentityToken,
			},
		},
		Spec: v1.ServiceSpec{
			Selector: map[string]string{
				"name": req.IdentityToken,
			},
			Type: v1.ServiceTypeClusterIP,
			Ports: []v1.ServicePort{
				{
					Name: "ssh",
					Port: 2222,
				},
			},
		},
	}
	_, err = s.client.CoreV1().Services("default").Create(c, &expectedService, metav1.CreateOptions{})
	if err != nil {
		c.JSON(500, err)
		return
	}

	c.JSON(201, created)
}
