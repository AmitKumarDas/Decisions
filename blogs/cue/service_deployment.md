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

### NOTES
- By replacing the service and deployment name with <Name> we have changed the definition into a template.
- Templates are applied to (are unified with) all entries in the struct in which they are defined:
  - so we need to either strip fields specific to the definition, 
  - generalize them, or
  - remove them.
- One of the labels defined in the Kubernetes metadata seems to be always set to parent directory name. 
  - We enforce this by defining component: string, 
  - meaning that a field of name component must be set to some string value, and then define this later on. 

