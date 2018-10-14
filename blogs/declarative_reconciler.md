### What is Reconcile?
A path (you may also call it _logic_) that constantly tries to set the current state of an entity to a desired state is known 
as **reconciliation**. In this article we shall be talking about an approach to manage reconciliation in a declarative manner.

_NOTE- Most of the low level stuff associated with reconciliation is abstracted to Kubernetes (to be more precise, its etcd database, its controllers, its client library and other utilities Kubernetes might be using)._

I believe reconcile as a workflow has advantages for specific usecases. However, my intention is not to list those fitting
usecases but would like to discuss about my mental models on declarative programming logic for reconcile.

### How does this help?
This is obviously the first question that boomerangs. In my opinion it does no harm to re-think a particular idea or process 
in the fields / disciplines where it was not meant to be in the first place. For example, why should core programmers (the ones
who know writing programming language syntax) be the only ones meant to use reconciliation? What about enabling users, 
non-technical personas to use this technique. Now one of the solutions, that come to my mind is to expose reconcilation into
simple declarative english like syntaxes. This helps the broader community of the project.

### Declarative Reconcile - High Level Design
We can start by designing a specification which is more english like versus code syntax like. Let us give this specification a
name. Let me call it _**ReconcileTask**_.

#### How does it look like - Take 1
```yaml
kind: ReconcileTask
spec:
  - select: clusterIP
    action: get kube service
    with:
    - "name": my-svc
  - select: name 
    action: create kube deploy
    with: 
    - "clusterip": clusterIP
```
#### How does it look like - Take 2
```yaml
kind: ReconcileTask
spec:
  - select clusterIP from a kube service having name my-svc
  - select name after creating a kube deploy with clusterip $clusterIP
```
#### How does it look like - Take 3
```yaml
kind: ReconcileTask
spec:
  - select clusterIP | get kube service | with "name" "my-svc"
  - select name | create kube deploy | with "clusterip" $clusterIP
```

#### Who understands this specification?
A Docker image `reconcile:0.0.1`. This is the logic (which abstracts Kubernetes reconcile operator logic) and accepts 
ReconcileTask spec as its input / config.
