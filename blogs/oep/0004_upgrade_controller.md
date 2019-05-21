#### Why?
- Dedicated controller(s) helps in achieving specific logic per resource kind

#### How?
- Upgrade will take help of some custom resources to perform upgrade
- These custom resources will have individual controllers
  - e.g. CStorPoolUpgrade, CStorVolumeUpgrade, JivaVolumeUpgrade, etc
- These controllers will make use of appropriate **UpgradeRecipe** to perform upgrade
  - Recipes provides a declarative model to code the requirements

#### Why UpgradeRecipe(s)?
- It helps in having only k8s controller related logic inside controllers
- Upgrade/Domain specific logic will reside in Recipes

#### UpgradeRecipe -- Deep Dive
- This will be a Kubernetes CR that is used to hold **static content**
- UpgradeRecipe is a thin abstraction around MayaLang
- It delegates to MayaLang to perform actions required to be at the desired upgrade state
- A recipe in itself does not have a controller
- A recipe is only used to create MayaLang resource
  - MayaLang resource has its controller/watcher

#### Upgrade controller(s)
- Controllers will watch for its respective resources
  - In this case these resources are:
    - CStorPoolUpgrade, 
    - CStorVolumeUpgrade, 
    - JivaVolumeUpgrade
- Controller logic will refer to corresponding recipe to perform following tasks:
  - Build a MayaLang resource based on the recipe
  - Fill in appropriate constants against MayaLang resource
  - Apply MayaLang resource against the K8s cluster
  - Delete MayaLang resource when its status shows completed
  - Update its watched resource's status

```yaml
kind: CStorPoolUpgrade
metadata:
  name: my-csp-upgrade
spec:
  sourceVersion:
  targetVersion:

  # optional; controller can decide the
  # recipe based on source & target
  # versions
  recipe:
    name: 

  # upgrade all cstor pools
  # optional; can be filled up by controller
  all: true

  # or upgrade based on the names
  # is optional 
  # controller sets spec.all to true
  # if this list is nil
  pools:
  - name: csp-pool1
  - name: csp-pool2
  - name: csp-pool3
status:
```

```yaml
kind: UpgradeRecipe
metadata:
  name: cstorpool-upgrade-from-082-to-090
  namespace: openebs
  labels:
    cstorpool.upgrade.openebs.io/sourceversion: 
    cstorpool.upgrade.openebs.io/targetversion:
spec:
  mlang:
    spec:
      go:
        vars:
          SourceVersion: "0.8.2"
          TargetVersion: "0.9.0"

          # PoolNamespace: "<filled-in-at-runtime>"
          # PoolName: "<filled-in-at-runtime>"

        funcs:
        - name: saveCStorPoolDetails
          body: |
            unstructops.New().
              WithStore(store).
              WithGVK("openebs.io", "v1alpha1", "CStorPool").
              WithNamespace(PoolNamespace).
              GetFromKubernetes(PoolName).
              SaveUIDToStoreWithKey("pool.uid").
              SaveToStore(
                unstructops.GetValueFromPath(".spec.device"),
                unstructops.WithStoreKey("pool.device"),
              )          

        - name: setCStorPoolVersion
          body: |
            cspops.New().
              GetFromKubernetes(PoolName).
              SetLabel("openebs.io/version", TargetVersion)
```
