# Console Installation

This guide assumes you are installing the Console on the latest, EKS-based (v3) Rack.  If you want to install your Console on an ECS-based v2 Rack, there are a couple of notes below of extra steps to perform.

## Application Setup

### Clone this repository

    $ git clone https://github.com/convox/console-app

### Create the Console application

You can use any name you like but this document will assume the name `console`.

    $ convox apps create console

### Set up private registry

Convox will provide credentials for a private registry to access the Console images.

Substitute `USERNAME` and `PASSWORD` in this command to add the private registry to your Rack.

    $ convox registries add enterprise.convox.com USERNAME PASSWORD

## Resource Stack Setup

### Create the stack

Create a new CloudFormation stack using the `formation.json` from this repository.

You can use any name you like but the rest of this document will assume the name `console-resources`.

You can do this easily via your AWS Web Console, uploading the `formation.json` at the appropriate stage, or using the `aws cli`

    $ aws cloudformation create-stack --stack-name console-resources --capabilities CAPABILITY_IAM --template-body file://formation.json

Wait for this stack to fully complete (can take ~10 minutes to complete depending on AWS).

### Configure Console environment

You will need to have [jq](https://stedolan.github.io/jq/) installed if you don't already.

    $ bin/export-env console-resources | convox env set -a console

## License Setup

Convox will provide you a license key for your Console.

    $ convox env set -a console LICENSE_KEY=...

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

### (OPTIONAL) Internal Mode

To make the Console only accessible inside the VPC, you will need to set Internal mode.

    $ convox env set -a console INTERNAL=true

(If you are deploying your Console app to an older v2 Rack, this will require your Rack to have the parameter `Internal=Yes` set)

## Deploy the Console

Deploy the application contained in this repository.

    $ convox deploy -a console

### Configure Console parameters (only required for older v2 Racks)

    $ convox apps params set RackUrl=Yes -a console

## (OPTIONAL) Integration Setup

If you'd like to use the GitHub, GitLab, or Slack integrations in your private Console you will need to create your own OAuth applications for each service.

Use the following callback URL(s) for each service:

| Provider | Callback URL(s)                                                                           |
|----------|-------------------------------------------------------------------------------------------|
| Github   | `https://$host/`                                                                          |
| Gitlab   | `https://$host/integrations/authorize/gitlab`<br>`https://$host/integrations/reauthorize` |
| Slack    | `https://$host/integrations/authorize/slack`                                              |

Once created, set the appropriate environment variables on your Console application:

    $ convox env set -a console GITHUB_CLIENT_ID=... GITHUB_CLIENT_SECRET=...
    $ convox env set -a console GITLAB_CLIENT_ID=... GITLAB_CLIENT_SECRET=...
    $ convox env set -a console SLACK_CLIENT_ID=... SLACK_CLIENT_SECRET=...

If you're using GitHub, you'll need to set a random webhook secret:

    $ convox env set -a console GITHUB_WEBHOOK_SECRET=...

If you'd like to use GitHub Enterprise, you'll also need to set the host:

    $ convox env set -a console GITHUB_ENTERPRISE_CLIENT_ID=... GITHUB_ENTERPRISE_CLIENT_SECRET=... GITHUB_ENTERPRISE_HOST=github.mycompany.org

Promote the environment changes

    $ convox releases promote -a console

## (OPTIONAL) LDAP Authentication

You can provide credentials for a secure (TLS) LDAP endpoint to use for authentication.

    $ convox env set -a console AUTHENTICATION=ldap
    $ convox env set -a console LDAP_ADDR=auth.example.org:636 LDAP_BIND=uid=%s,dc=example,dc=org

Set `LDAP_BIND` to a full bind string where `%s` will be substituted for the user's email address.

If your LDAP server does not have a valid certificate issued by a known CA, you can disable certificate validation:

    $ convox env set -a console LDAP_VERIFY=no

Promote the environment changes

    $ convox releases promote -a console

## (OPTIONAL) SAML Authentication

You can provide configuration details to use SAML for authentication.

    $ convox env set -a console AUTHENTICATION=saml
    $ convox env set -a console SAML_METADATA=https://login.microsoftonline.com/common/FederationMetadata/2007-06/FederationMetadata.xml

`SAML_METADATA` should be set to the metadata endpoint for your SAML Identity Provider.  This varies from provider to provider so please check your documentation from them.

Promote the environment changes

    $ convox releases promote -a console
