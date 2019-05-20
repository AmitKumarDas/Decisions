### refer
- https://github.com/moby/moby/tree/17.05.x/pkg/mount

### Nice use of maps

```go
var flags = map[string]struct {
	clear bool
	flag  int
}{
	"defaults":      {false, 0},
	"ro":            {false, RDONLY},
	"rw":            {true, RDONLY},
	"suid":          {true, NOSUID},
	"nosuid":        {false, NOSUID},
}

var validFlags = map[string]bool{
	"":          true,
	"size":      true,
	"mode":      true,
	"uid":       true,
	"gid":       true,
	"nr_inodes": true,
	"nr_blocks": true,
	"mpol":      true,
}

var propagationFlags = map[string]bool{
	"bind":        true,
	"rbind":       true,
	"unbindable":  true,
	"runbindable": true,
	"private":     true,
	"rprivate":    true,
	"shared":      true,
	"rshared":     true,
	"slave":       true,
	"rslave":      true,
}
```
