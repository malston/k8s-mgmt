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
	"fmt"
	"io"

	kubernetes "k8s.io/client-go/kubernetes/fake"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type FakeClient struct {
	FakeKubeClientset *kubernetes.Clientset
	Context           string
	Stdin             io.Reader
	Stdout            io.Writer
	Stderr            io.Writer
}

func (c *FakeClient) Core() corev1.CoreV1Interface {
	return c.FakeKubeClientset.CoreV1()
}

func (c *FakeClient) CurrentContext() string {
	return c.Context
}

func (c *FakeClient) SetContext(context string) error {
	c.Context = context
	return nil
}

func NewClient() *FakeClient {
	kubeClientset := kubernetes.NewSimpleClientset()

	return &FakeClient{
		FakeKubeClientset: kubeClientset,
	}
}

func (c *FakeClient) Printf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stdout, format, a...)
}

func (c *FakeClient) Eprintf(format string, a ...interface{}) (n int, err error) {
	return fmt.Fprintf(c.Stderr, format, a...)
}
