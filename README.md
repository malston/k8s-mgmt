# Kubernetes Management

## Plan Requirements

### Configure Compute Resources for System Daemons

* Kubelet system-reserved
  * Reserve Compute Resources for System Daemons. Enter a comma separated list of parameters e.g. `memory=250Mi`, `cpu=150m`

* Kubelet eviction-hard
  * Hard eviction thresholds set for worker kubelet to kill pods when the set thresholds are reached. e.g. `memory.available=100Mi`, `nodefs.available=10%`, `nodefs.inodesFree=5%`

### Configure Admission Control Plugins

An [admission controller](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) is a piece of code that intercepts requests to the Kubernetes API server prior to persistence of the object, but after the request is authenticated and authorized.

## Cluster Requirements

### Create and Configure Cluster

* Configure RBAC for Users/Groups in a Cluster

* Configure Pod Security Policies

* Update/Resize Cluster

## Namespace Requirements

By default, a Kubernetes cluster will instantiate a default namespace when provisioning the cluster to hold the default set of Pods, Services, and Deployments used by the cluster. Best practice is to subdivide a cluster into multiple namespaces where each team, organization, or application gets its own namespace. You may also want to give each developer a separate namespace ensuring that one developer can not accidentally delete another developers work. Namespaces can also serve as scopes for the deployment of services so that one application's front-end service doesn't interfere with another app's front-end service.

### Create and Configure Namespace

#### Configure RBAC for Users/Groups in a Namespace

Before you can assign a user to a namespace, you have to onboard that user to the Kubernetes cluster itself. To achieve this, there are two options. You can use certificate based authentication to create a new certificate for the user and give them a `kubeconfig` file which they can use to login or you can configure your cluster to use an external identity system (for example Active Directory) to access their cluster.

In general, using an external identity system is a best practice since it doesn't require that you maintain two different sources of identity, but in some cases this isn't possible and certificates need to be used. Fortunately, you can use the Kubernetes certificate API for creating and managing such certificates.

After the certificate has been added to the `kubeconfig` file you will need to apply Kubernetes RBAC for the user to grant them privileges to a namespace otherwise the user has no access privileges.

* Configure Default Memory Requests and Limits for a Namespace

* Configure Default CPU Requests and Limits for a Namespace

* Configure Minimum and Maximum Memory Constraints for a Namespace

* Configure Minimum and Maximum CPU Constraints for a Namespace

* Configure Memory and CPU Quotas for a Namespace

## Flow

Issue a command like `create-namespaces`.

* Loop through files under the config directory and find all the namespace folders and open each `namespace.yml` file.
* Create a new namespace based on contents of `namespace.yml`.

### Example Code

```code
m := config.NewManager("/Users/user/.cfg")
for n := range m.GetNamespaces() {
   n.Create()
}
 ```
