# Console Installation

This guide provides instructions for installing the Convox Console on an ECS-based (v2) Rack.

> **Note**: If you are installing your Console on an EKS-based v3 Rack, use the `master-v3` branch, which contains v3-compatible instructions.

## Prerequisites

Before beginning the Console installation, ensure you have the following tools installed:

- **Convox CLI**: Follow the [installation instructions](https://docsv2.convox.com/installation/cli)
- **AWS CLI**: Required for AWS operations ([installation guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html))
- **jq**: Required for JSON processing ([download here](https://stedolan.github.io/jq/))

Verify your installations:
```bash
$ convox version
$ aws --version
$ jq --version
```

## Rack Installation

### Standard AWS Installation

For standard AWS installations, install a rack using the Convox CLI. First, ensure you're logged into AWS CLI and that your region is set:

```bash
$ export AWS_REGION=us-east-1
$ aws sts get-caller-identity
```

Then install the rack with recommended parameters:

```bash
$ convox rack install aws -n console-rack \
    InstanceType=c6i.large \
    BuildInstance=c6i.large
```

> **Note**: 
> - `console-rack` is the name for your new rack - you can change this to any name you prefer
> - A v2 Rack is created in the region of your current AWS credentials. Set `AWS_REGION` before installing; there is no region parameter
> - The `c6i.large` instance type is our recommended cost-effective node size for Console hosting racks
> - You can use smaller instance types (e.g., `t3.medium`) for additional cost savings in smaller environments
> - Monitor your rack's resource utilization after deployment - you can tune down to smaller nodes later if they are underutilized
> - These parameters can be changed at any time using `convox rack params set`, with the exception of `HighAvailability`, which is only accepted during installation
> - Rack installation typically takes 20-30 minutes and streams CloudFormation events, ending with:
> ```
> 2026-01-15T12:10:00Z system/cloudformation aws/cfm console-rack CREATE_COMPLETE AWS::CloudFormation::Stack
> Rack console-rack installed successfully
> ```

### AWS GovCloud Installation

For AWS GovCloud deployments, create a new rack locally using your AWS GovCloud credentials:

```bash
$ export AWS_REGION=us-gov-east-1
$ convox rack install aws -n gov-console-rack \
    InstanceType=c6i.large \
    BuildInstance=c6i.large
```

> **Important**: 
> - `gov-console-rack` is the name for your new rack - you can change this to any name you prefer
> - Set `AWS_REGION` to the appropriate GovCloud region (e.g., `us-gov-east-1` or `us-gov-west-1`) before installing
> - The `c6i.large` instance type is our recommended cost-effective node size for Console hosting racks
> - You can use smaller instance types (e.g., `t3.medium`) for additional cost savings in smaller environments
> - Monitor your rack's resource utilization after deployment - you can tune down to smaller nodes later if they are underutilized
> - These parameters can be changed at any time using `convox rack params set`

## Application Setup

### Clone the Repository

Clone the Console application repository and switch to the master-v2 branch:

```bash
$ git clone https://github.com/convox/console-app && cd console-app && git checkout master-v2
```

> **Note**: The `master-v2` branch contains the latest Console version compatible with v2 (ECS-based) racks.

### Create the Console Application

First, switch to your newly created rack:

```bash
$ convox switch console-rack
```

> **Note**: If you chose a different rack name during installation, use that name instead. Alternatively, you can skip this step and append `-r rackName` to all subsequent commands.

Create the Console application with your chosen name. This guide uses `console` as the example:

```bash
$ convox apps create console
```

> **Note**: To use a different name, replace `-a console` with `-a yourAppName` in all subsequent commands.

### Configure Private Registry

Convox will provide credentials for accessing Console images. Add the private registry to your Rack:

```bash
$ convox registries add enterprise.convox.com USERNAME PASSWORD
```

## Resource Stack Setup

> **Upgrading from a version before 3.0.26?** The `formation.json` in this release adds new DynamoDB tables required by Console 3.0.26+. You must update your existing CloudFormation stack before deploying the new Console version. See [Updating an Existing Resource Stack](#updating-an-existing-resource-stack) below. New installations can skip this notice.

### Create the CloudFormation Stack

Create a new CloudFormation stack using the `formation.json` from this repository.

**For Standard AWS:**
```bash
$ aws cloudformation create-stack \
    --stack-name console-resources \
    --capabilities CAPABILITY_IAM \
    --template-body file://formation.json \
    --region us-east-1 \
    --output text \
    --no-cli-pager
```

**For AWS GovCloud:**
```bash
$ aws cloudformation create-stack \
    --stack-name console-resources \
    --capabilities CAPABILITY_IAM \
    --parameters ParameterKey=AwsArn,ParameterValue=aws-us-gov \
    --template-body file://formation.json \
    --region us-gov-east-1 \
    --output text \
    --no-cli-pager
```

> **Note**: If you installed your rack in a different region, adjust the `--region` parameter accordingly.

Wait for the stack to complete (approximately 10 minutes). You can check the status with:

```bash
$ aws cloudformation describe-stacks \
    --stack-name console-resources \
    --query 'Stacks[0].StackStatus' \
    --region us-east-1 \
    --output text \
    --no-cli-pager
```

> **Note**: You can also monitor the CloudFormation stack progress in the AWS Console under CloudFormation → Stacks → console-resources.

### Updating an Existing Resource Stack

If you are upgrading from a Console version before 3.0.26, update your existing `console-resources` stack with the latest `formation.json` to add the new required DynamoDB tables. This is a safe, additive operation that creates new resources without modifying existing tables or data.

**For Standard AWS:**
```bash
$ aws cloudformation update-stack \
    --stack-name console-resources \
    --capabilities CAPABILITY_IAM \
    --template-body file://formation.json \
    --region us-east-1 \
    --output text \
    --no-cli-pager
```

**For AWS GovCloud:**
```bash
$ aws cloudformation update-stack \
    --stack-name console-resources \
    --capabilities CAPABILITY_IAM \
    --parameters ParameterKey=AwsArn,ParameterValue=aws-us-gov \
    --template-body file://formation.json \
    --region us-gov-east-1 \
    --output text \
    --no-cli-pager
```

Wait for the stack update to complete. You can check the status with:

```bash
$ aws cloudformation describe-stacks \
    --stack-name console-resources \
    --query 'Stacks[0].StackStatus' \
    --region us-east-1 \
    --output text \
    --no-cli-pager
```

The stack status should transition from `UPDATE_IN_PROGRESS` to `UPDATE_COMPLETE`. Once complete, proceed with deploying the new Console version.

### Configure Console Environment

Once the `console-resources` CloudFormation stack shows `CREATE_COMPLETE` status, export the stack outputs as environment variables:

```bash
$ bin/export-env console-resources | convox env set -a console
```

> **Note**: Ensure the CloudFormation stack has fully completed before running this command, as it needs to read the stack outputs.

You can verify the environment variables were set correctly:

```bash
$ convox env -a console
```

## License Setup

Convox will provide you with a license key for your Console:

```bash
$ convox env set -a console LICENSE_KEY=your-license-key-here
```

## Custom Domain Setup

### Create SSL Certificate

On a v2 Rack, certificates are issued by AWS Certificate Manager and are managed at the Rack level rather than per app. You have two options:

**Option 1: Request a certificate through ACM**
```bash
$ convox certs generate console.example.org
```

ACM sends a validation email to the registered contacts for the domain's apex. Accept the email to complete validation.

**Option 2: Import an existing certificate**
```bash
$ convox certs import cert.pem key.pem
```

> **Note**: If no existing certificate matches the domain you set in `HOST`, deploying the Console requests one through ACM automatically. That request is also email-validated, and the CloudFormation stack waits on it, so approve the validation email promptly or the deploy will appear to hang.

### Configure DNS

Create a CNAME record pointing your custom domain to your Rack's router:

1. Get your Rack's router address:
```bash
$ convox rack
Name      console-rack
Provider  aws
Region    us-east-1
Router    console-rack-Route-1A2B3C4D5E6F-123456789.us-east-1.elb.amazonaws.com
Status    running
Version   20260826164715
```

2. Create a CNAME record:
```
console.example.org → console-rack-Route-1A2B3C4D5E6F-123456789.us-east-1.elb.amazonaws.com
```

> **Note**: Use a simple routing policy for the CNAME record.

### Configure HOST Environment Variable

```bash
$ convox env set -a console HOST=console.example.org
```

## Deploy the Console

Deploy the Console application:

```bash
$ convox deploy -a console
```

> **Note**: The first deployment typically takes 5-10 minutes as it needs to create the Redis cache resource. You will see repeated messages stating `Waiting on dependent resources to be completed` - this is normal while the Redis instance is being provisioned.

### Configure Console Parameters

**Required**: grant the Console access to the API of the Rack hosting it.

```bash
$ convox apps params set RackUrl=Yes -a console
```

This adds `RACK_URL` to the Console's task definitions, which is how the Console reaches its own Rack to run Rack installs, updates and workflow jobs. The command triggers a CloudFormation update of the app stack and takes a few minutes.

> **Important**: Without this parameter, the Console starts and serves the UI normally, but every Rack install started from the Console will create the Rack record and then stall with the install terminal showing `Loading ...`. Setting `RACK_URL` with `convox env set` does not work as a substitute, because the Rack only passes through environment variables declared in `convox.yml`.

Confirm it applied:

```bash
$ convox apps params -a console | grep RackUrl
RackUrl                                Yes
```

## Verification

After deployment completes, you should be able to access your Console at your configured domain (e.g., https://console.example.org).

## Configure Redis Cache

After deployment, configure the Redis cache by setting the CACHE_REDIS_ADDR environment variable.

**Automatic configuration:**
```bash
$ convox env set -a console CACHE_REDIS_ADDR=$(convox resources -a console | sed -n 's/.*redis.*redis:\/\/\(.*\)\/0/\1/p') && convox releases promote -a console
```

> **Note**: This command sets the Redis cache address and immediately promotes the release to apply the configuration.

**Manual configuration:**
1. Get the Redis URL:
```bash
$ convox resources -a console
NAME   TYPE   URL
cache  redis  redis://cache-console-a1b2c3d4.e5f6g7.ng.0001.use1.cache.amazonaws.com:6379/0
```

2. Set only the URI and port (exclude the redis:// prefix and /0 suffix):
```bash
$ convox env set -a console CACHE_REDIS_ADDR=cache-console-a1b2c3d4.e5f6g7.ng.0001.use1.cache.amazonaws.com:6379
```

3. Promote the release:
```bash
$ convox releases promote -a console
```

## Importing the Rack into the Console

At this point, you have a CLI-managed rack hosting your Console application. To enable team management and full Console features, import the rack into the Console so it can manage its own hosting infrastructure.

### Import the Rack

1. Navigate to your Console URL (e.g., https://console.example.org)
2. Register a new user and create your organization
3. Collect the rack's hostname and password:
   - **hostname**: AWS CloudFormation → Stacks → filter for `console-rack` → Outputs tab → the value of the `Dashboard` key
   - **password**: AWS Systems Manager → Parameter Store → filter for `console-rackRackApiSecret` → the parameter value
4. From the "Racks" page in the Console, click "+ Import"
5. Fill in the exact rack name, the hostname, and the password
6. Click "Add Rack"

> **Note**: Replace `console-rack` with your rack name if different. The SSM parameter name is your rack name with `RackApiSecret` appended, with no separator.

7. Log in to the Console from your CLI so it addresses the imported rack:
   - Click the user account icon in the upper right corner of the Console
   - In the Account Settings panel, click the refresh/reset button next to "CLI Token" to generate a new CLI key
   - Copy and run the provided login command in your terminal

8. Switch to the rack, which now has your organization name attached:

```bash
$ convox switch console-rack
```

9. Verify the rack is addressed and running:

```bash
$ convox rack
```

10. Confirm the rack appears with its organization prefix:

```bash
$ convox racks
NAME                    PROVIDER  STATUS
orgName/console-rack    aws       running
```

### Create an AWS Runtime Integration

> **Note**: Ensure you're authenticated to the AWS Console on the account where you installed the rack before proceeding.

1. Go to the "Integrations" page in the Console
2. Click the "+" button in the Runtime tab
3. Select "AWS or AWS Gov Web Services"
4. Click "Launch Stack" to open AWS CloudFormation
5. Check "I acknowledge that AWS CloudFormation might create IAM resources"
6. Create the stack
7. Wait 1-2 minutes and refresh the Integrations page to confirm installation

### Assign the Integration to Your Rack

1. Navigate to "Racks" in the Console
2. Select your rack
3. Go to "Rack Settings" → "General Settings"
4. In the "Runtime" dropdown, select your newly created integration
5. Click "Save Changes"

## Optional Integrations

### GitHub, GitLab, and Slack

Create OAuth applications for each service you want to integrate:

| Provider | Callback URL(s) |
|----------|-----------------|
| GitHub | `https://console.example.org/` |
| GitLab | `https://console.example.org/integrations/authorize/gitlab` and `https://console.example.org/integrations/reauthorize` |
| Slack | `https://console.example.org/integrations/authorize/slack` |

Set the environment variables:

```bash
$ convox env set -a console \
    GITHUB_CLIENT_ID=... \
    GITHUB_CLIENT_SECRET=... \
    GITHUB_WEBHOOK_SECRET=...

$ convox env set -a console \
    GITLAB_CLIENT_ID=... \
    GITLAB_CLIENT_SECRET=...

$ convox env set -a console \
    SLACK_CLIENT_ID=... \
    SLACK_CLIENT_SECRET=...
```

For GitHub Enterprise:

```bash
$ convox env set -a console \
    GITHUB_ENTERPRISE_CLIENT_ID=... \
    GITHUB_ENTERPRISE_CLIENT_SECRET=... \
    GITHUB_ENTERPRISE_HOST=github.mycompany.org
```

### LDAP Authentication

Configure LDAP authentication:

```bash
$ convox env set -a console \
    AUTHENTICATION=ldap \
    LDAP_ADDR=auth.example.org:636 \
    LDAP_BIND=uid=%s,dc=example,dc=org
```

To disable certificate validation (if needed):

```bash
$ convox env set -a console LDAP_VERIFY=no
```

### SAML Authentication

**Standard SAML (e.g., Okta):**

```bash
$ convox env set -a console \
    AUTHENTICATION=saml \
    SAML_METADATA=https://dev-12345678.okta.com/app/exk1a2b3c4d5e6f7g8/sso/saml/metadata
```

> **Note**: The metadata URL format varies by provider:
> - **Okta**: `https://{your-okta-domain}/app/{app-id}/sso/saml/metadata`
> - **Azure AD**: `https://login.microsoftonline.com/{tenant-id}/FederationMetadata/2007-06/FederationMetadata.xml`
> - Check your SAML provider's documentation for the correct metadata endpoint

**Google SAML:**

1. Set up a SAML app in Google Admin:
   - ACS URL: `https://console.example.org/saml`
   - Entity ID: `https://console.example.org`
   - Name ID: Basic Information → Primary email

2. Configure attribute mapping for Primary email to NameID

3. Download the metadata XML file and host it publicly (e.g., S3 bucket, GitHub Pages)

4. Set environment variables:

```bash
$ convox env set -a console \
    AUTHENTICATION=saml-google \
    SAML_METADATA=https://your-hosted-metadata-url.com/metadata.xml
```

### Apply Configuration Changes

After setting any optional configurations, promote the release:

```bash
$ convox releases promote -a console
```

## Verification for Authentication Methods

If you configured SAML or LDAP authentication:

1. Navigate to your Console URL (e.g., https://console.example.org)
2. You should be redirected to your identity provider's login page
3. After successful authentication, you'll be redirected back to the Console
4. For SAML: Ensure your user attributes are properly mapped (email, name, etc.)
5. For LDAP: Verify your bind DN format is working with your credentials

> **Troubleshooting Auth Issues**: 
> - Check logs with `convox logs -a console` for authentication errors
> - Verify your metadata URL (SAML) or LDAP server is accessible from the Console
> - Ensure callback URLs are correctly configured in your identity provider

## Optional Tuning

These environment variables have safe defaults and do not need to be set for normal operation. They are available for operators who want to adjust logging, retention, or performance characteristics.

> **Note**: A v2 Rack only passes through environment variables declared in `convox.yml`. Anything you set with `convox env set` that is not listed there is stored but never reaches the container, with no error. If you need a variable this guide does not cover, add it to both services in `convox.yml` and redeploy.

| Variable | Type | Default | Valid Range | Description |
|----------|------|---------|-------------|-------------|
| `LOG_LEVEL` | string | `info` | `info`, `verbose`, `debug` | Controls Console logging verbosity. `info` logs errors only. `verbose` adds DynamoDB operation timing. `debug` adds Redis cache hit/miss logging. |
| `APP_EVENT_TTL_DAYS` | integer | `90` | 0+ (0 disables) | Retention horizon in days for app lifecycle events in DynamoDB. Table-level TTL pruning must be enabled separately via AWS CLI. |
| `CONSOLE_COST_ROLLUP_RACK_CONCURRENCY` | integer | `8` | 1-24 | Maximum concurrent rack fan-out during cost rollup queries. Lower values reduce peak DynamoDB load at the cost of slower cost page rendering. |

To set any of these:

```bash
$ convox env set -a console LOG_LEVEL=verbose
$ convox releases promote -a console
```

## Troubleshooting

- If certificate generation stalls, check the certificate's status in AWS Certificate Manager. ACM validation emails go to the registered contacts for the domain's apex and expire after 72 hours
- If a Rack install started from the Console stalls at `Loading ...`, confirm `convox apps params -a console` reports `RackUrl  Yes`
- For Redis connection issues, verify the CACHE_REDIS_ADDR contains only the hostname and port
- For GovCloud installations, ensure `AWS_REGION` matches your GovCloud region
- Check application logs: `convox logs -a console`

## Support

For additional support, please contact Convox support with any error messages you encounter.