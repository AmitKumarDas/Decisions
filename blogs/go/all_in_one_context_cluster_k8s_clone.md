refer - https://github.com/kubermatic/kubeone/blob/master/pkg/util/context.go

```go
// Context hold together currently test flags and parsed info, along with
// utilities like logger
type Context struct {
	Cluster                   *kubeoneapi.KubeOneCluster
	Logger                    logrus.FieldLogger
	Connector                 *ssh.Connector
	Configuration             *Configuration
	Runner                    *Runner
	WorkDir                   string
	JoinCommand               string
	JoinToken                 string
	RESTConfig                *rest.Config
	DynamicClient             dynclient.Client
	Verbose                   bool
	BackupFile                string
	DestroyWorkers            bool
	ForceUpgrade              bool
	UpgradeMachineDeployments bool
}

// Clone returns a shallow copy of the context.
func (c *Context) Clone() *Context {
	newCtx := *c
	return &newCtx
}
```
