
1. Setup the environment variables.

    ```
    export KL_STATE_DIRECTORY="/tmp/kl"
    export KL_IAAS="azure"
    
    export KL_AZURE_ENVIRONMENT=
    export KL_AZURE_REGION=
    export KL_AZURE_SUBSCRIPTION_ID=
    export KL_AZURE_TENANT_ID=
    export KL_AZURE_CLIENT_ID=
    export KL_AZURE_CLIENT_SECRET=
    ```

1. `go install` under the directory `cmd/kl/`.

1. `kl digest`

1. `kl plan_infra`

1. `kl apply_infra`
