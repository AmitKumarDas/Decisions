```go
// https://github.com/apptio/kr8
// https://github.com/apptio/kr8/blob/master/cmd/jsonnet.go
```

```go
// imports

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	goyaml "github.com/ghodss/yaml"
	jsonnet "github.com/google/go-jsonnet"
	jsonnetAst "github.com/google/go-jsonnet/ast"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/apimachinery/pkg/util/yaml"
)
```

```go
vm := jsonnet.MakeVM()
vm.Importer(&jsonnet.FileImporter{
  JPaths: jpath,
})
for _, extvar := range extvarfiles {
  kv := strings.SplitN(extvar, "=", 2)
  v, _ := ioutil.ReadFile(kv[1])
  vm.ExtVar(kv[0], string(v))
}
```

- Does a lot of stuff to prepare a jsonnet file
- Adds certain native functions that understands various yaml, json stuff
