# Attestation config API

## Azure SEV-SNP
The version numbers of SEV-SNP are updated as part of [e2e_verify](/.github/actions/e2e_verify/action.yml).
Because the version numbers are not publicly posted by Azure, we observe the versions on Azure VMs and assume a global rollout after a threshold time.

This estimate might make manual intervention necessary when a global rollout didn't happen.

### Manually delete a version
```
COSIGN_PASSWORD=$CPW COSIGN_PRIVATE_KEY="$(cat $PATH_TO_KEY)" AWS_ACCESS_KEY_ID=$ID AWS_ACCESS_KEY=$KEY bazel run //internal/api/attestationconfigapi/cli -- delete azure-sev-snp attestation-report 2025-01-18-09-15
```

### Manually upload a version
```
COSIGN_PASSWORD=$CPW COSIGN_PRIVATE_KEY="$(cat $PATH_TO_KEY)" AWS_ACCESS_KEY_ID=$ID AWS_ACCESS_KEY=$KEY bazel run //internal/api/attestationconfigapi/cli -- upload azure-sev-snp attestation-report 2025-01-18-09-15 --force
```
