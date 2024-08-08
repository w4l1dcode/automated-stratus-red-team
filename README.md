# automated-stratus-red-team
Stratus Red Team is "Atomic Red Team™" for the cloud. This GO program runs automated red team tests on AWS. The code will create the infrastructure needed for the tests, once finished, the infrastructure will also get cleaned up.

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
