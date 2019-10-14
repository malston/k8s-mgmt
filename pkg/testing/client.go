/*
Copyright © 2019 Mark Alston <marktalston@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package testing

import (
	kubernetes "k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type FakeKubeClient struct {
	FakeKubeClientset *kubernetes.Clientset
	Context           string
}

func (c *FakeKubeClient) Core() corev1.CoreV1Interface {
	return c.FakeKubeClientset.CoreV1()
}

func (c *FakeKubeClient) CurrentContext() string {
	return c.Context
}

func (c *FakeKubeClient) SetContext(context string) error {
	c.Context = context
	return nil
}

func NewKubeClient() *FakeKubeClient {
	kubeClientset := kubernetes.NewSimpleClientset()

	return &FakeKubeClient{
		FakeKubeClientset: kubeClientset,
	}
}
