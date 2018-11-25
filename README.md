# Console Installation

## Application Setup

### Clone this repository

    $ git clone https://github.com/convox/console-app -b console2

### Create the Console application

You can use any name you like but this document will assume the name `console`.

    $ convox apps create console --wait

### Set up private registry

Convox will provide credentials for a private registry to access the Console images.

Substitute `USERNAME` and `PASSWORD` in this command to add the private registry to your Rack.

    $ convox registries add registry.convox.com USERNAME PASSWORD

## Resource Stack Setup

### Create the stack

Create a new CloudFormation stack using the `formation.json` from this repository.

You can use any name you like but the rest of this document will assume the name `console-resources`.

Wait for this stack to fully complete.

### Configure app environment

    $ bin/export-env console-resources | convox env set -a console

## License Setup

Convox will provide you a license for your Console.

    $ convox env set -a console LICENSE=...

## Custom Domain Setup

Decide on a custom domain for your Console. These instructions will assume `console.example.org`.

### Create a certificate

Create an SSL certificate for your application. You can use `convox certs import` to load a certificate 
that you create manually or you can use the Rack's built-in certificate generator.

    $ convox certs generate console.example.org

If you use the automatic certificate generation you will need to accept the certificate validation email that will be sent to the DNS administrator of the domain.

### Set up DNS

Create a CNAME record for this domain to point at the `Router` attribute shown when you run `convox rack`.

### Configure app environment

    $ convox env set -a console HOST=console.example.org

### Deploy the Console

Deploy the application contained in this repository.

    $ convox deploy -a console --wait

### (OPTIONAL) Integration Setup

If you'd like to use the GitHub, GitLab, or Slack integrations in your private Console you will need to create your own OAuth applications for each service.

Use the following callback URL(s) for each service:

| Provider | Callback URL(s)                                                                        |
|----------|----------------------------------------------------------------------------------------|
| Github   | `https://$host/`                                                                       |
| Gitlab   | `https://$host/integrations/authorize/gitlab`<br>`https://$host/integrations/reauthorize` |
| Slack    | `https://$host/integrations/authorize/slack`                                           |

Once created, set the appropriate environment variables on your Console application:

    $ convox env set -a console GITHUB_CLIENT_ID=... GITHUB_CLIENT_SECRET=...
    $ convox env set -a console GITLAB_CLIENT_ID=... GITLAB_CLIENT_SECRET=...
    $ convox env set -a console SLACK_CLIENT_ID=... SLACK_CLIENT_SECRET=...

If you'd like to use GitHub enterprise, you'll also need to set the endpoint:

    $ convox env set -a console GITHUB_ENDPOINT=https://...

Promote the environment changes

    $ convox releases promote -a console --wait

### (OPTIONAL) LDAP Authentication

You can provide credentials for a secure (TLS) ldap endpoint to use for authentication.

    $ convox env set -a console AUTHENTICATION=ldap
    $ convox env set -a console LDAP_ADDR=auth.example.org:636 LDAP_BASE=dc=example,dc=org

Promote the environment changes

    $ convox releases promote -a console --wait
