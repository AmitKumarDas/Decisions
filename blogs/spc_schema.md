```yaml
apiVersion: openebs.io/v1alpha1
kind: StoragePoolClaim
metadata:
  name: cstor-sparse-pool-striped
spec:
  nodeList:
  - name: node-xyz
    diskGroupList:
    - name: group-1
      diskList:
      - name: sparse-554270d60a8e1d8768c63062da8cdb40
        otherDetailsSpec:
          key1: val1
          key2: val2
      - name: sparse-742b4c20b1047b40d1497a98d4bba97b
        otherDetailsSpec:
          key1: val1
          key2: val2
    - name: group-2
      diskList:
      - name: sparse-554270d60a8e1d8768c63062da8cdb40
        otherDetailsSpec:
          key1: val1
          key2: val2
      - name: sparse-742b4c20b1047b40d1497a98d4bba97b
        otherDetailsSpec:
          key1: val1
          key2: val2
  - name: node-abc
    diskGroupList:
    - name: group-1
      diskList:
      - name: sparse-554270d60a8e1d8768c63062da8cdb40
        otherDetailsSpec:
          key1: val1
          key2: val2
      - name: sparse-742b4c20b1047b40d1497a98d4bba97b
        otherDetailsSpec:
          key1: val1
          key2: val2
    - name: group-2
      diskList:
      - name: sparse-554270d60a8e1d8768c63062da8cdb40
        otherDetailsSpec:
          key1: val1
          key2: val2
      - name: sparse-742b4c20b1047b40d1497a98d4bba97b
        otherDetailsSpec:
          key1: val1
          key2: val2
  pool:
    cacheFile: ""
    overProvisioning: false
    type: striped
  type: sparse
  ```
