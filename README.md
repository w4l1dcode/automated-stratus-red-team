# automated-stratus-red-team
Stratus Red Team is "Atomic Red Team™" for the cloud. This Go program automates red team tests on AWS by creating the necessary infrastructure, executing the tests, and then cleaning up the infrastructure afterward. For each run, the program randomly selects a tactic and executes all associated (available) red team tests, with a 20-minute delay between each attack.

## Prerequisites
A sandbox environment to run the red team tests on. You will need admin permissions on that environment.

## Running
First create a yaml file, such as `config.yml`:
```yaml
log:
  level: DEBUG

aws:
  access_key_id: ""
  secret_access_key: ""
  session_token: ""
```

And now run the program from source code:
```shell
% make
go run ./cmd/... -config=dev.yml
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
