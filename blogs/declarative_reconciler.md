### What is Reconcile?
Logic that constantly tries to set the current state of cluster to a desired state is known as reconciliation. In this article
we shall be abstracting the low level reconciliation details to Kubernetes.

I believe reconcile as a strategy has advantages for specific usecases. However, my intention is not to get into the fitting
usecases but would like to discuss about my thoughts on declarative programming logic for reconcile.

### How does this help?
This helps the broader community e.g. non-coders to take advantage of reconcile process abstracted by Kubernetes. These users 
will be able to declaratively specify the reconcile logic without thinking about writing pieces of code.

### Declarative Reconcile - High Level Design
- ReconcileTask - A Kubernetes CustomResource
- Kubernetes watchers, informers & so on will do the job of reconciling the current state of cluster to desired state
- reconcile:0.0.1 - A Docker image. A Kubernetes controller that understands executing ReconcileTask specification.
