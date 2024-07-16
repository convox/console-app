# Console Installation

This guide assumes you are installing the Console on the latest, EKS-based (v3) Rack.  If you want to install your Console on an ECS-based v2 Rack, there are a couple of notes below of extra steps to perform.

## Rack Installation

This step is only necessary if you're using AWS Gov Cloud.
Regular AWS Cloud installation doesn't require a local rack specifically.

Create a new rack locally using your AWS Gov Cloud credentials

You can pass the same [parameters](https://docs.convox.com/installation/production-rack/aws/) as in the UI using the form key=value

    $ convox rack install aws gov-rack region=us-gov-east-1

don't forget to pass the gov cloud region

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

You can do this easily via your AWS Web Console, uploading the `formation.json` at the appropriate stage, or using the `aws cli`:

    $ aws cloudformation create-stack --stack-name console-resources --capabilities CAPABILITY_IAM --template-body file://formation.json

If you are using AWS GovCloud, you have to set the `AwsArn` parameter as `aws-us-gov`:

    $ aws cloudformation create-stack --stack-name console-resources --capabilities CAPABILITY_IAM --parameters ParameterKey=AwsArn,ParameterValue=aws-us-gov --template-body file://formation.json

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

### Configure New Console environment

For v3 rack

    $ bin/export-env-v3 console | convox env set -a console

For v2 rack

    $ bin/export-env-v2 console | convox env set -a console

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

## Steps below are only necessary if you're installing in Gov Cloud

### Move the local rack to the Console

First, go to the console UI, register a new user and name your organization.

Then go to the Account Tab and click "Reset CLI Key" - run the output command in your terminal.

Now, run the following command

    $ convox rack mv gov-rack orgName/gov-rack

Where `gov-rack` is the name of the rack you created and `orgName` the name of the organization you created in the console.

Go to the "Racks" tab and confirm your rack was moved.


### Create an AWS Runtime Integration

Note that this step is only necessary once.

In the Console, go to "Integrations" and click the "+" next to "Runtime".

Select "AWS Gov Web Services" and click "Launch Stack". This will take you to the AWS CloudFormation UI.

Before creating the stack, make sure to check the box near "I acknowledge that AWS CloudFormation might create IAM resources.".

This stack is creating an IAM Role that the console will use to make changes to rack resources on your behalf.


Wait 1-2 minutes and refresh the Integrations UI to confirm the integration was installed

You can now use this Runtime to create additional racks directly from the console.

### Assign the Integration to the moved rack

This step is only necessary for racks you move into the console.

Go to the "Racks" section, click the blue cog icon near your moved rack and in the "Runtime" dropdown, select your newly created Integration. Click "Apply Changes"

### Assign proper console permissions to your moved rack

This step is only necessary for racks you move into the console.

Follow the steps in our [docs](https://docs.convox.com/management/console-rack-management/) - Moving an AWS Rack section to give the console role permission to access the EKS rack cluster.

## Updating to the New Console (Console3)

To update the rack to a new version, you must first update the console-resources CloudFormation stack.
Due to AWS limitations, this will need to be done over three commands as you cannot add multiple Global Secondary Indexes (GSIs) in a single run.

The following update commands assume you have followed this guide and your Console Resources stack is named console-resources. If you installed the Convox Console with a different CloudFormation stack name, you should adjust the `--stack-name` option accordingly.
Each update will take about a minute or less, and you can check the status of the CloudFormation stack in AWS directly or with the following CLI command:

```sh
aws cloudformation describe-stacks --stack-name console-resources --query 'Stacks[0].[StackName, StackStatus]' --output text
```

If you used a custom value for any of the stack parameters, make sure to include them in the CloudFormation update command.
You can check the "Parameters" section of the CloudFormation stack within the AWS Management Console if you're not sure, or with this command:

```sh
aws cloudformation describe-stacks --stack-name console-resources --query 'Stacks[0].Parameters' --output table
```

A default installed stack will produce this output:
-------------------------------------
|          DescribeStacks           |
+---------------+-------------------+
| ParameterKey  |  ParameterValue   |
+---------------+-------------------+
|  SseEnabled   |  false            |
|  TtlEnabled   |  false            |
|  AwsArn       |  aws              |
|  TablePrefix  |  console-private  |
+---------------+-------------------+

E.g.: if you have used a custom value in TablePrefix:

```sh
aws cloudformation update-stack --stack-name console-resources --capabilities CAPABILITY_IAM --parameters ParameterKey=TablePrefix,ParameterValue=CustomValue --template-body file://formation-update-1.json
```

For multiple custom parameters, the format would be like this:
```sh
aws cloudformation update-stack --stack-name console-resources --capabilities CAPABILITY_IAM --parameters ParameterKey=TablePrefix,ParameterValue=CustomValue1 ParameterKey=AnotherParameter,ParameterValue=CustomValue2 ParameterKey=YetAnotherParameter,ParameterValue=CustomValue3 --template-body file://formation-update-1.json
```

Once you've verified that you're ready to update the CloudFormation stack, please run the following commands to create the GSIs:

```sh
aws cloudformation update-stack --stack-name console-resources  --capabilities CAPABILITY_IAM  --parameters ParameterKey=TablePrefix,UsePreviousValue=true --template-body file://formation-update-1.json
```

```sh
aws cloudformation update-stack --stack-name console-resources  --capabilities CAPABILITY_IAM  --parameters ParameterKey=TablePrefix,UsePreviousValue=true --template-body file://formation-update-2.json
```

```sh
aws cloudformation update-stack --stack-name console-resources  --capabilities CAPABILITY_IAM  --parameters ParameterKey=TablePrefix,UsePreviousValue=true --template-body file://formation-update-3.json
```

After the stack updates successfully, export the updated ENV to your console app with the command:

```sh
bin/export-env console-resources | convox env set -a console
```

Then run the appropriate command depending on your Console Rack's engine.

v2
```sh
bin/export-env-v1 console | convox env set -a console
```

or

v3
```sh
bin/export-env-v3 console | convox env set -a console
```

Finally, deploy the app to update the Convox Console:

```sh
convox deploy -a console
```

**Warning**

If running a very old version of the console, it's possible that the `bin/export-env` script will output `RACK_KEY` and `SESSION_KEY` values different from the ones you already set (you can get the current set values by running `convox env -a console`). If that's the case, just remove the variables from the `bin/export-env` output and set the other ones. You can also roll back the release if you forgot to remove them.

```sh
bin/export-env console-resources | convox env set -a console
```

Then run

v2
```sh
bin/export-env-v1 console | convox env set -a console
```

or

v3
```sh
bin/export-env-v3 console | convox env set -a console
```

```sh
convox deploy -a console
```
