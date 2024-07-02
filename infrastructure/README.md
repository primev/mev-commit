# infrastructure

## Secrets Management

All sensitive information, including all keystores with their passwords, are stored in the secrets vault. For your devnet deployments, you can access them on the same IP as Nomad but with port 8200. Use the `root_token` from `~/.vault_init.json` stored on the same machine to sign in. 

For devenv, all the secrets are different every time a new deployment is made, so it does not really matter if the secrets leak or not (this is not the same case for prod, though).
