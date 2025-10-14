# Console Installation

This guide provides instructions for installing the Convox Console on the latest EKS-based (v3) Rack. 

> **Note**: If you need to install your Console on an ECS-based v2 Rack, please use the v2-specific release tag (e.g., `3.0.17-v2`) which contains v2-compatible instructions.

## Prerequisites

Before beginning the Console installation, ensure you have the following tools installed:

- **Convox CLI**: Follow the [installation instructions](https://docs.convox.com/installation/cli)
- **Terraform**: Required for rack installation ([download here](https://www.terraform.io/downloads))
- **AWS CLI**: Required for AWS operations ([installation guide](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html))
- **jq**: Required for JSON processing ([download here](https://stedolan.github.io/jq/))

Verify your installations:
```bash
$ convox version
$ terraform version
$ aws --version
```

## Rack Installation

### Standard AWS Installation

For standard AWS installations, install a rack using the Convox CLI. First, ensure you're logged into AWS CLI:

```bash
$ aws sts get-caller-identity
```

Then install the rack with recommended parameters:

```bash
$ convox rack install aws console-rack region=us-east-1 \
    node_type=c6i.large \
    build_node_enabled=true \
    build_node_type=c6i.large
```

> **Note**: 
> - `console-rack` is the name for your new rack - you can change this to any name you prefer
> - Adjust the `region` parameter to your desired AWS region (e.g., `us-west-2`, `eu-west-1`)
> - The `c6i.large` instance type is our recommended cost-effective node size for Console hosting racks
> - You can use smaller instance types (e.g., `t3.medium`) for additional cost savings in development environments
> - Monitor your rack's resource utilization after deployment - you can tune down to smaller nodes later if they are underutilized
> - These parameters (except region) can be changed at any time using `convox rack params set`
> - Rack installation typically takes 20-30 minutes. Upon completion, you'll see output similar to:
> ```
> api = <sensitive>
> provider = "aws"
> release = "3.22.3"
> ```

### AWS GovCloud Installation

For AWS GovCloud deployments, create a new rack locally using your AWS GovCloud credentials:

```bash
$ convox rack install aws gov-console-rack region=us-gov-east-1 \
    node_type=c6i.large \
    build_node_enabled=true \
    build_node_type=c6i.large
```

> **Important**: 
> - `gov-console-rack` is the name for your new rack - you can change this to any name you prefer
> - Specify the appropriate GovCloud region (e.g., `us-gov-east-1` or `us-gov-west-1`)
> - The `c6i.large` instance type is our recommended cost-effective node size for Console hosting racks
> - You can use smaller instance types (e.g., `t3.medium`) for additional cost savings in development environments
> - Monitor your rack's resource utilization after deployment - you can tune down to smaller nodes later if they are underutilized
> - These parameters (except region) can be changed at any time using `convox rack params set`
## Application Setup

### Clone the Repository

```bash
$ git clone https://github.com/convox/console-app && cd console-app && git checkout master-v3
```

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

You have two options for SSL certificates:

**Option 1: Generate with Let's Encrypt (Recommended for initial setup)**
```bash
$ convox certs generate console.example.org
```

This will send a certificate validation email to the DNS administrator. Accept the email to complete validation.

**Option 2: Import existing certificate**
```bash
$ convox certs import cert.pem key.pem -a console
```

> **Note**: You can configure DNS-01 challenge with Route53 for automated certificate renewal later. See the [documentation](https://docs.convox.com/deployment/ssl#advanced-ssl-configuration-lets-encrypt-dns01-challenge-with-route53) for details.

### Configure DNS

Create a CNAME record pointing your custom domain to your Rack's router:

1. Get your Rack's router address:
```bash
$ convox rack
Name      console-rack
Provider  aws
Router    router.0a1b2c3d4e5f.convox.cloud
Status    running
Version   3.22.3
```

2. Create a CNAME record:
```
console.example.org → router.0a1b2c3d4e5f.convox.cloud
```

> **Note**: Use a simple routing policy for the CNAME record.

### Configure HOST Environment Variable

```bash
$ convox env set -a console HOST=console.example.org
```

### (Optional) Internal Mode

> **Advanced Configuration**: This is an advanced security configuration. We highly recommend completing the initial Console setup and verification before enabling internal mode. This allows you to ensure everything is working correctly while the Console is still publicly accessible.

To make the Console only accessible within your VPC:

```bash
$ convox env set -a console INTERNAL=true
```

> **Important**: When enabling internal mode, the Console will not be accessible from the public internet. You will need one of the following to access it:
> - **AWS VPN**: Set up a Client VPN or Site-to-Site VPN connection to your VPC
> - **Bastion Host**: Deploy a bastion/jump host in a public subnet to tunnel through
> - **AWS Systems Manager Session Manager**: Use Session Manager to access instances within the VPC
> - **Direct Connect**: If you have AWS Direct Connect established to your VPC
> - **VPC Peering**: If accessing from another peered VPC with appropriate routing
> 
> Ensure you have one of these access methods configured before enabling internal mode, or you will lose access to the Console UI.

## Deploy the Console

Deploy the Console application:

```bash
$ convox deploy -a console
```

> **Note**: The first deployment typically takes 5-10 minutes as it needs to create the Redis cache resource. You will see repeated messages stating `Waiting on dependent resources to be completed` - this is normal while the Redis instance is being provisioned.

## Verification

After deployment completes, you should be able to access your Console at your configured domain (e.g., https://console.example.org).

## Configure Redis Cache

After deployment, configure the Redis cache by setting the CACHE_REDIS_ADDR environment variable.

**Automatic configuration:**
```bash
$ convox env set -a console CACHE_REDIS_ADDR=$(convox resources -a console | sed -n 's/.*elasticache-redis.*redis:\/\/\(.*\)\/0/\1/p') && convox releases promote -a console
```

> **Note**: This command sets the Redis cache address and immediately promotes the release to apply the configuration.

**Manual configuration:**
1. Get the Redis URL:
```bash
$ convox resources -a console
NAME   TYPE               URL
cache  elasticache-redis  redis://cache-console-a1b2c3d4.e5f6g7.ng.0001.use1.cache.amazonaws.com:6379/0
```

2. Set only the URI and port (exclude the redis:// prefix and /0 suffix):
```bash
$ convox env set -a console CACHE_REDIS_ADDR=cache-console-a1b2c3d4.e5f6g7.ng.0001.use1.cache.amazonaws.com:6379
```

3. Promote the release:
```bash
$ convox releases promote -a console
```

## Moving the Rack into the Console

At this point, you have a CLI-managed rack hosting your Console application. To enable team management and full Console features, you need to transfer ownership of this rack from your local CLI to the Console application itself. This allows the Console you just deployed to manage its own hosting infrastructure:

### Move the Local Rack to the Console

1. Navigate to your Console URL (e.g., https://console.example.org)
2. Register a new user and create your organization
3. Click on the user account icon in the upper right corner of the Console
4. In the Account Settings panel, click the refresh/reset button next to "CLI Token" to generate a new CLI key
5. Copy and run the provided login command in your terminal
6. Move the rack to your organization:

```bash
$ convox rack mv console-rack orgName/console-rack
```

> **Note**: Replace `console-rack` with your rack name if different, and `orgName` with your actual organization name you set at user registration.

7. After moving the rack, you'll need to switch to it again as it now has a new name:

```bash
$ convox switch console-rack
```

8. Verify the rack is now managed by the Console:

```bash
$ convox rack
```

You should see the rack details showing it's now accessible through the Console.

9. Confirm the rack appears with its organization prefix:

```bash
$ convox racks
NAME                    PROVIDER  STATUS
orgName/console-rack    aws       running
```

10. Verify the rack appears in the Console UI by navigating to the "Racks" tab (may take up to 30 seconds after `convox rack mv` command)

### Create an AWS Runtime Integration

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

### Configure Console Permissions

**Critical Step**: This step is required for the Console to properly manage the EKS cluster. It grants the Console's IAM role permission to access your rack's Kubernetes cluster.

1. First, get your kubectl configuration for the rack:

```bash
$ convox rack kubeconfig > ~/.kube/config
```

2. Verify you're connected to the correct cluster:

```bash
$ kubectl get ns
```

You should see namespaces including your rack name and app namespaces.

3. Find the Console's IAM role ARN in AWS:
   - Go to your AWS IAM Console
   - Search for a role named `convox-YOURORGID-ConsoleRole-XXXXXXXXXXXX`
   - Copy the full ARN (e.g., `arn:aws:iam::123456789012:role/convox-abc123-ConsoleRole-1234567890`)
   - **Important**: If the ARN contains `/convox/` between `role/` and the role name, remove that part

4. Edit the aws-auth ConfigMap to grant the Console access:

```bash
$ kubectl edit configmap/aws-auth -n kube-system
```

5. Add the following entry to the `mapRoles` section (be careful with YAML formatting):

```yaml
    - rolearn: arn:aws:iam::YOURACCOUNTID:role/convox-YOURORGID-ConsoleRole-XXXXXXXXXXXX
      username: convox-console
      groups:
      - system:masters
```

> **Note**: Replace the `rolearn` with the actual ARN you copied in step 3.

6. Save and exit the editor (`:wq` in vim, `Ctrl+X` then `Y` in nano)

7. Verify the change was applied:

```bash
$ kubectl get configmap aws-auth -n kube-system -o yaml | grep -A2 -B1 convox-console
```

You should see your newly added role mapping in the output.

8. For a complete verification, check the entire ConfigMap:

```bash
$ kubectl get configmap aws-auth -n kube-system -o yaml
```

The output should look similar to this (with your specific values):

```yaml
apiVersion: v1
data:
  mapRoles: |
    - rolearn: arn:aws:iam::123456789012:role/console-rack-nodes
      groups:
      - system:bootstrappers
      - system:nodes
      username: system:node:{{EC2PrivateDNSName}}
    - rolearn: arn:aws:iam::123456789012:role/convox-abc123-ConsoleRole-XXXXXXXXXXXX
      username: convox-console
      groups:
      - system:masters
kind: ConfigMap
metadata:
  name: aws-auth
  namespace: kube-system
```

> **Important**: Be very careful when editing the aws-auth ConfigMap. Incorrect changes can lock you out of your cluster. If you're not comfortable with kubectl commands or have any concerns, please reach out to Convox support for assistance.
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

## Troubleshooting

- If certificate generation fails, ensure your DNS is properly configured and the domain is accessible
- For Redis connection issues, verify the CACHE_REDIS_ADDR contains only the hostname and port
- For GovCloud installations, ensure all region parameters match your GovCloud region
- Check application logs: `convox logs -a console`

## Support

For additional support, please contact Convox support with any error messages you encounter.