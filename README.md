# cloud-threat-emulation
This Go program automates red team testing on AWS and Kubernetes using Stratus Red Team. It sets up the required infrastructure, executes the red team tests, and cleans up afterward. Each run randomly chooses a tactic and performs all associated red team tests for that tactic, with a 20-minute delay between each attack.
## Prerequisites
- A sandbox account
- A kubernetes cluster
- A persistent access key and secret
- an aws-auth config map.
  - Firs create a yaml file, `aws-auth.yaml`:
    ```yaml
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: aws-auth
      namespace: kube-system
    data:
      mapRoles: |
        # Your existing roles here
      mapUsers: |
        - userarn: arn:aws:iam::<ACCOUNT_ID>:user/<USERNAME>
          username: <USERNAME>
          groups:
            - system:masters
    ```
  - Apply  the changes with:
    ```shell
    kubectl apply -f aws-auth.yaml
    ```

## Running
First create a yaml file, such as `config.yml`:
```yaml
log:
  level: DEBUG

aws:
  access_key_id: ""
  secret_access_key: ""
  aws_region: ""

kubernetes:
  cluster_name: ""
  k8s_region: ""
```

And now run the program from source code:
```shell
% make
go run ./cmd/... -config=config.yml
INFO[0000] set log level                                 fields.level=debug
INFO[0000] Starting AWS tests...  
INFO[0000] Executing tests of tactic: Impact         
INFO[0000] Number of ttps found: 3                      
INFO[0000] Executing ttp: aws.impact.s3-ransomware-batch-deletion   
```

## Building

```shell
% make build
```

## Reference
[Stratus Red Team](https://github.com/DataDog/stratus-red-team) by DataDog
