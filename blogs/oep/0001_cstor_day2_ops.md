### Motivation
Storage as an infrastructure requires proper maintainenance to serve the needs of higher order applications. We shall try to 
cover these activities as Day 2 operations. Day 2 as the name suggests are the activites that a storage admin would like to 
perform after initial storage related provisioning is done and applications have started to make use of this storage.

### UseCases
#### Add or Remove Disks
_As a storage admin I would like to [1] add new disks to increase capacity or [2] replace bad disks with new ones to keep storage up and running without any side-effects._

_As a storage admin I would like to add disk to a specific pool i.e. pool on a specific node._

_As a storage admin I would like to use kubectl to add or remove disk from pool._

_As a storage admin I would like to decomission nodes and create new nodes with automation. I would like to avoid the manual
steps such as: fetch disk list, update pool from the helm chart, push it to git, update the catalog and then update the pool._

### Custom Resources
#### SPC


#### CSPC


### TestCases
_NOTE: These test code will be available at below paths:_
  - github.com/openebs/maya/tests

- [ ] Removing a disk from pool should remove this disk from the storage
  - maya/tests/
- [ ] Replacing an old disk with a new disk should remove the old disk & add the new disk into the storage
  - maya/tests/
- [ ] Removing a storage node from pool should wipe off the storage associated with the node
  - maya/tests/
- [ ] Commissioning a new storage node should automatically update the pool with relevant disks from this node
  - maya/tests/
- [ ] Decommissioning an existing storage node should automatically update the pool by removing relevant disks 
  - maya/tests/
