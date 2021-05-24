[Back to Table of Contents](../documentation.md)

# Kubernetes Common Use-Cases and Patterns

The following set of examples will help you quickly achieve common tasks in the context of Kubernetes YAML manifests.  

All of these examples are available for your convenience in [examples/kubernetes](../../examples/kubernetes) and are intended to be launched from the root of your local copy of the YAML Overlay Tool repository:

```bash
yot -i examples/kubernetes/< example you wish to run >.yaml -o < desired output path >
```

<br/>


### Example Kubernetes YAML Manifests

Within the [examples/kubernetes/manifests](../../examples/kubernetes/manifests) directory of the YAML Overlay Tool repository, you will find the two example Kubernetes YAML Manifests which we will be manipulating in the following set of example use-cases:

```yaml
# my-app.yaml
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    my.custom.annotation/fake: idk
  labels:
    name: my-web-page
  name: my-web-page
  namespace: my-web-page
spec:
  containers:
    - image: nginx:latest
      name: my-web-page
      ports:
        - containerPort: 443
      resources: {}
  dnsPolicy: ClusterFirst
  restartPolicy: Always

```

```yaml
# my-service.yaml
---
apiVersion: v1
kind: Service
metadata:
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-type: nlb
  labels:
    name: my-web-page
  name: my-service
  namespace: my-web-page
spec:
  ports:
    - name: 8443-443
      port: 8443
      protocol: TCP
      targetPort: 443
  selector:
    name: my-web-page
  type: LoadBalancer

```

<br/>


## Adding Additional Labels and Selectors To All YAML Files In a Directory

In this example we use the `merge` action to add in 2 new labels to a set of YAML files contained in a directory.

```yaml
---
# addLabels.yaml
commonOverlays:
  - name: Add additional labels
    query:
      - metadata.labels
      - spec.template.metadata.labels
      - spec.selector.matchLabels
      - spec.selector
    value:
      app.kubernetes.io/owner: Jeff Smith
      app.kubernetes.io/purpose: static-webpage
    action: merge

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now apply the changes by generating a new set of YAML files:  
>`yot -i ./examples/kubernetes/addLabels.yaml -o /tmp/new`


## Prepend a Private Registry URL to All Container Images

In this example we use the `merge` action to take the images which are currently pointing to the Docker Hub registry (only base imagename:tag), and prepending with a private registry URL.  

We also demonstrate use of the `%v` format marker by inserting the original value into a new line comment.

```yaml
---
# privateContainerRegistry.yaml
commonOverlays:
  - name: Set our private container registry in manifests
    query: ..image
    value: my-private-reg/%v  # old value was: %v
    action: merge

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now apply the changes by generating a new set of YAML files:  
>`yot -i ./examples/kubernetes/privateContainerRegistry.yaml -o /tmp/new`


## Modify the Name of a Label's Key

In this example we will manipulate the `name` label key with `app.kubernetes.io/name` by using the `merge` action and retaining the existing value.  The `~` character in JSONPath+ always returns the value of the key, rather than the value of the key/value pair.

```yaml
---
# formatLabelKey.yaml
commonOverlays:
  - name: Update name label's key to app.kubernetes.io/name
    query: metadata.labels.name~
    value: app.kubernetes.io/%v  # the old key was %v
    action: merge
  - name: 
    query: spec.selector
    value: 
      app.kubernetes.io/%k: "%v"  # the old key was %k
yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now apply the changes by generating a new set of YAML files:  
>`yot -i ./examples/kubernetes/formatLabelKey.yaml -o /tmp/new`

## Replace the Name of a Label's Key

In this example we will replace the `name` label with `my-new-label` by using the `replace` action and retaining the existing value. The `~` character in JSONPath+ always retures the value of the key, rather than the value of the key/value pair.

```yaml
---
# replaceLabelKey.yaml
commonOverlays:
  - name: Replace name label's key to my-new-label
    query: metadata.labels.name~
    value: my-new-label  # old label key was: %v
    action: replace

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now apply the changes by generating a new set of YAML files:  
>`yot -i ./examples/kubernetes/replaceLabelKey.yaml -o /tmp/new`


## Remove All Annotations

Often times annotations are set for certain environments that may not apply to your environment.  To remove all annotations we will use the `delete` action.

```yaml
# removeAnnotations.yaml
commonOverlays:
  - name: Remove all annotations
    query: metadata.annotations
    action: delete

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now apply the changes by generating a new set of YAML files:  
>`yot -i ./examples/kubernetes/removeAnnotations.yaml -o /tmp/new`


## Remove Annotations from Specific Kubernetes Object Types

To build on the previous example, there are often times when you may want to remove annotations from specific Kubernetes types, or a combination of conditions.  To remove these annotations, we will use the `delete` action and a `documentQuery`.

```yaml
# removeAnnotationsWithConditions.yaml
commonOverlays:
  - name: Remove all annotations with conditions
    query: metadata.annotations
    action: delete
    documentQuery:
      - conditions:
          - key: kind
            value: Service
          - key: metadata.namespace
            value: my-web-page
      - conditions:
          - key: metadata.name
            value: my-service

yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now apply the changes by generating a new set of YAML files:  
>`yot -i ./examples/kubernetes/removeAnnotationsWithConditions.yaml -o /tmp/new`


## insert Comments

Sometimes you would like to add some additional comments to denote why you did something.  Additionally, to document what something used to be set to in case you want to restore it to a previous state later.

There are other times where you want to add comments as annotations for another application's consumption.

Yot **can** insert comments!

```yaml
---
# insertComments.yaml
commonOverlays:
  - name: insert line comments via merge
    query:
      - metadata.annotations['my.custom.annotation/fake']
      - metadata.annotations['service.beta.kubernetes.io/aws-load-balancer-type']
    value: "%v" # insert a line comment
    action: merge
  - name: insert a line comment via replace
    query: spec.containers[0].image
    value: new-image:latest # old value was: %v
    action: replace
  - name: insert head, foot, and line comments via merge
    query: metadata.labels
    value:
      # insert a head comment
      app.kubernetes.io/owner: Jeff Smith  # insert a line comment
      app.kubernetes.io/purpose: static-webpage  # insert another line comment
      # insert a foot comment
    action: merge
yamlFiles:
  - name: Set of Kubernetes manifests from upstream
    path: ./examples/kubernetes/manifests
```

Now apply the changes by generating a new set of YAML files:  
>`yot -i ./examples/kubernetes/insertComments.yaml -o /tmp/new`


[Back to Table of Contents](../documentation.md)  
[Next Up: Interactive Tutorials and Learning Paths](tutorials.md)