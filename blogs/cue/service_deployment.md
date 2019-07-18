```cue
$ cat <<EOF > kube.cue
package kube

service <Name>: {
    apiVersion: "v1"  
    kind:       "Service"
    metadata: {
        name: Name
        labels: {
            app:       Name    // by convention
            domain:    "prod"  // always the same in the given files
            component: string  // varies per directory
        }
    }
    spec: {
        // Any port has the following properties.
        ports: [...{
            port:       int
            protocol:   *"TCP" | "UDP"      // from the Kubernetes definition
            name:       string | *"client"
        }]
        selector: metadata.labels // we want those to be the same
    }
}

deployment <Name>: {
    apiVersion: "extensions/v1beta1"
    kind:       "Deployment"
    metadata name: Name
    spec: {
        // 1 is the default, but we allow any number
        replicas: *1 | int
        template: {
            metadata labels: {
                app:       Name
                domain:    "prod"
                component: string
            }
            // we always have one namesake container
            spec containers: [{ name: Name }]
        }
    }
}
EOF
```
