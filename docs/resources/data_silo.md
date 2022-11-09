---
page_title: "transcend_data_silo Resource - terraform-provider-transcend"
subcategory: ""
description: |-
  
---

# transcend_data_silo (Resource)



## Example Usages

### Connecting an API Key integration

You can connect your data silos with API Key or other similar configuration through Transcend. In this flow, Transcend still does not ever see your API keys, but encrypts them through your internal Sombra before storing the encrypted values in Transcend's backend.

Before we configure the data silo, it's worth noting that for self-hosted sombras you will need to add an authentication token to your internal sombra. This can be done by adding the `internal_sombra_key` field in the provider or the `TRANSCEND_INTERNAL_SOMBRA_KEY` environment variable with the value of the internal key you used when setting up your sombra service. If you are using Transcend-hosted sombra as your encryption gateway, you can skip this step as just your API key can authenticate you.

In the data silo, add the `secret_context` values for each field you want to specify. Here's an example of fully connecting a Datadog silo:

```terraform
variable "dd_api_key" { sensitive = true }
variable "dd_app_key" { sensitive = true }

resource "transcend_data_silo" "datadog" {
  type            = "datadog"
  skip_connecting = false

  secret_context {
    name  = "apiKey"
    value = var.dd_api_key
  }
  secret_context {
    name  = "applicationKey"
    value = var.dd_app_key
  }
  secret_context {
    name  = "queryTemplate"
    value = "service:programmatic-remote-seeding AND @email:{{identifier}}"
  }
}
```

### Connecting an AWS Silo

Connecting Amazon to Transcend is done through [AWS IAM Roles](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles.html). In any AWS Account you want us to have access to audit, you need to create an IAM Role allowing our AWS organization access to it. This is the recommended pattern from Amazon [documented here](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_roles_create_for-user_externalid.html). This is done in a few steps:

- You create an AWS data silo in Transcend.
- We provide you with an external id through a resource output
- You create an IAM Policy for what permissions Transcend can take in your organization
- You create an IAM Role that only allows Transcend to access it, and only when using the given external ID

```terraform
resource "transcend_data_silo" "aws" {
  type        = "amazonWebServices"
  description = "Amazon Web Services (AWS) provides information technology infrastructure services to businesses in the form of web services."

  # Normally, Data Silos are connected in this resource. But for AWS, we want to delay connecting until after
  # we create the IAM Role, which must use the `aws_external_id` output from this resource. So instead, we set
  # `skip_connecting` to `true` here and use a `transcend_data_silo_connection` resource below
  skip_connecting = true
  lifecycle { ignore_changes = [plaintext_context] }
}

resource "aws_iam_role" "iam_role" {
  name        = "TranscendAWSIntegrationRole2"
  description = "Policy to allow Transcend access to this AWS Account"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Action = "sts:AssumeRole"
        Effect = "Allow"
        // 829095311197 is the AWS Organization for Transcend that will try to assume role into your organization
        Principal = { AWS = "arn:aws:iam::829095311197:root" }
        Condition = { StringEquals = { "sts:ExternalId" : transcend_data_silo.aws.aws_external_id } }
      },
    ]
  })

  inline_policy {
    name = "TranscendPermissions"
    policy = jsonencode({
      Version = "2012-10-17"
      Statement = [
        {
          Action = [
            "dynamodb:ListTables",
            "dynamodb:DescribeTable",
            "rds:DescribeDBInstances",
            "s3:ListAllMyBuckets"
          ]
          Effect   = "Allow"
          Resource = "*"
        },
      ]
    })
  }
}

# Give AWS Time to become consistent with the new IAM Role permissions
resource "time_sleep" "pause" {
  depends_on      = [aws_iam_role.iam_role]
  create_duration = "10s"
}

data "aws_caller_identity" "current" {}
resource "transcend_data_silo_connection" "connection" {
  data_silo_id = transcend_data_silo.aws.id

  plaintext_context {
    name  = "role"
    value = aws_iam_role.iam_role.name
  }

  plaintext_context {
    name  = "accountId"
    value = data.aws_caller_identity.current.account_id
  }

  depends_on = [time_sleep.pause]
}
```

The above example completes this entire flow. The `transcend_data_silo_connection` resource ensures that the IAM Role is assumable by Transcend. 

### Creating an Automated Vendor Coordination Silo

```terraform
resource "transcend_data_silo" "avc" {
  type                 = "promptAPerson"
  outer_type           = "coupa"
  notify_email_address = "dpo@coupa.com"
  description          = "Coupa is a cloud platform for business spend that offers a fully unified suite of financial applications for business spend management"
  is_live              = true
}
```

In most cases, Transcend aims to provide an API-based integration into your SaaS tools and internal systems. Sometimes, the API based approach may not be possible due to restrictions such as:

- The SaaS vendor has no API for a certain type(s) of privacy requests, and requests you to send an email in a specific template.
- The SaaS vendor only provides a self-serve dashboard to submit data privacy requests, and someone on your team would need to log in and submit the request through their browser.
- You want to notify an internal team to perform a manual process against a database or internal tool.

In these cases, you can configure your Transcend instance to automate the sending of an email template whenever a particular type of data subject request is made. These emails can be sent to the SaaS vendor directly, or to an individual in your organization.

### Connecting a database silo

In this example, we will create a Postgres relational database using Amazon's Relational Database Service that is in a private network. We will then deploy a sombra encryption gateway inside the private network that can talk to the database. With that connection, we can then create a data silo in Transcend that will use sombra as a middle man to communicate to the database, ensuring your data remains private from both Transcend and the outside internet.

```terraform
locals {
  subdomain = "https-test"
  # You should pick a hosted zone that is in your AWS Account
  parent_domain = "sombra.dev.trancsend.com"
  # Org URI found on https://app.transcend.io/infrastructure/sombra
  organization_uri = "wizard"
}

######################################################################################
# Create a private network to put our database in with the sombra encryption gateway #
######################################################################################

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 2.18.0"

  name = "sombra-example-https-test-vpc"
  cidr = "10.0.0.0/16"
  azs  = ["us-east-1a", "us-east-1b"]

  private_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]
  public_subnets   = ["10.0.201.0/24", "10.0.202.0/24"]
  database_subnets = ["10.0.103.0/24", "10.0.104.0/24"]

  enable_nat_gateway                 = true
  enable_dns_hostnames               = true
  enable_dns_support                 = true
  create_database_subnet_group       = true
  create_database_subnet_route_table = true
}

#######################################################################
# Deploy a Sombra encryption gateway and register it to a domain name #
#######################################################################

data "aws_route53_zone" "this" {
  name = local.parent_domain
}

module "acm" {
  source      = "terraform-aws-modules/acm/aws"
  version     = "~> 2.0"
  zone_id     = data.aws_route53_zone.this.id
  domain_name = "${local.subdomain}.${local.parent_domain}"
}

variable "tls_cert" {}
variable "tls_key" {}
variable "jwt_ecdsa_key" {}
variable "internal_key_hash" {}
module "sombra" {
  source  = "transcend-io/sombra/aws"
  version = "1.4.1"

  # General Settings
  deploy_env       = "example"
  project_id       = "example-https"
  organization_uri = local.organization_uri

  # This should not be done in production, but allows testing the external endpoints during development
  transcend_backend_ips = ["0.0.0.0/0"]

  # VPC settings
  vpc_id                      = module.vpc.vpc_id
  public_subnet_ids           = module.vpc.public_subnets
  private_subnet_ids          = module.vpc.private_subnets
  private_subnets_cidr_blocks = module.vpc.private_subnets_cidr_blocks
  aws_region                  = "us-east-1"
  use_private_load_balancer   = false

  # DNS Settings
  subdomain       = local.subdomain
  root_domain     = local.parent_domain
  zone_id         = data.aws_route53_zone.this.id
  certificate_arn = module.acm.this_acm_certificate_arn

  # App settings
  data_subject_auth_methods = ["transcend", "session"]
  employee_auth_methods     = ["transcend", "session"]

  # HTTPS Configuration
  desired_count = 1
  tls_config = {
    passphrase = "unsecurePasswordAsAnExample"
    cert       = var.tls_cert
    key        = var.tls_key
  }
  transcend_backend_url = "https://api.dev.trancsend.com:443"

  # The root secrets that you should generate yourself and keep secret
  # See https://docs.transcend.io/docs/security/end-to-end-encryption/deploying-sombra#6.-cycle-your-keys for information on how to generate these values
  jwt_ecdsa_key     = var.jwt_ecdsa_key
  internal_key_hash = var.internal_key_hash

  tags = {}
}

######################################################################
# Create a security group that allows Sombra to talk to the database #
######################################################################

module "security_group" {
  source  = "terraform-aws-modules/security-group/aws"
  version = "~> 4.0"

  name   = "database-ingress"
  vpc_id = module.vpc.vpc_id

  # ingress
  ingress_with_cidr_blocks = [
    {
      from_port   = 5432
      to_port     = 5432
      protocol    = "tcp"
      description = "PostgreSQL access from private subnets within VPC (which includes sombra)"
      cidr_blocks = join(",", module.vpc.private_subnets_cidr_blocks)
    },
  ]
}

###################################################
# Create a sample postgres database using AWS RDS #
###################################################

module "postgresDb" {
  source  = "terraform-aws-modules/rds/aws"
  version = "~> 5.0"

  allocated_storage    = 5
  engine               = "postgres"
  engine_version       = "11.14"
  family               = "postgres11"
  major_engine_version = "11"
  instance_class       = "db.t3.micro"

  multi_az               = true
  db_subnet_group_name   = module.vpc.database_subnet_group
  vpc_security_group_ids = [module.security_group.security_group_id]
  skip_final_snapshot    = true
  deletion_protection    = false
  apply_immediately      = true

  identifier = "some-postgres-db"
  username   = "someUsername"
  db_name    = "somePostgresDb"
}

#######################################################
# As Sombra can talk to the database, we can create a #
# data silo using the private connection information. #
#######################################################

resource "transcend_data_silo" "database" {
  type = "database"

  schema_discovery_plugin {
    enabled                    = true
    schedule_frequency_minutes = 1440 # 1 day
    schedule_start_at          = "2022-09-06T17:51:13.000Z"
  }

  secret_context {
    name  = "driver"
    value = "PostgreSQL Unicode"
  }
  secret_context {
    name = "connectionString"
    value = join(";", [
      "Server=${module.postgresDb.db_instance_address}",
      "Database=${module.postgresDb.db_instance_name}",
      "UID=${module.postgresDb.db_instance_username}",
      "PWD=${module.postgresDb.db_instance_password}",
      "Port=${module.postgresDb.db_instance_port}"
    ])
  }
}
```

This type of setup works for many types of relational databases including Postgres, Amazon Redshift, Snowflake, MySQL, Google BigQuery, MariaDb, and Microsoft SQL. We are always adding support for more relational database types, please reach out to support if you are wanting to connect a type you don't see existing documentation for.

### Adding a Data Silo Plugin

Sometimes when you connect a data silo, you don't actually have data in that silo, but may want to connect that silo so that Transcend can find other silos from that tool. An example might be connecting Okta to find SaaS tools you use, or connecting AWS to find your S3 buckets and databases.

To do so, add a `data_silo_discovery_plugin` block inside any data silo that supports silo discovery like so:

```terraform
resource "transcend_data_silo" "aws" {
  type        = "amazonWebServices"
  description = "Amazon Web Services (AWS) provides information technology infrastructure services to businesses in the form of web services."

  data_silo_discovery_plugin {
    enabled                    = true
    schedule_frequency_minutes = 1440 # 1 day
    schedule_start_at          = "2022-09-06T17:51:13.000Z"
  }

  # ...other fields...
}

# ...other resources...
```

### Adding a Data Point Plugin

Sometimes when you connect a data silo, you don't know what data is inside that silo, and want Transcend to try to classify the data for you. This is called data point discovery, and can be performed by certain integrations such as Salesforce, Amazon S3/Dynamodb, Mongodb, Snowflake, and more.

To do so, add the `schema_discovery_plugin` and `content_classification_plugin` blocks inside any data silo that supports data point discovery like so:

```terraform
resource "transcend_data_silo" "aws" {
  type        = "amazonWebServices"
  description = "Amazon Web Services (AWS) provides information technology infrastructure services to businesses in the form of web services."

  schema_discovery_plugin {
    enabled                    = true
    schedule_frequency_minutes = 1440 # 1 day
    schedule_start_at          = "2022-09-06T17:51:13.000Z"
  }

  content_classification_plugin {
    enabled                    = true
    schedule_frequency_minutes = 1440 # 1 day
    schedule_start_at          = "2022-09-06T17:51:13.000Z"
  }

  # ...other fields...
}

# ...other resources...
```

## Looking up Data Silo metadata

If you are wondering what integration names Transcend supports or what fields are available on those integrations, you can lookup all data silo metadata via our GraphQL API.

Go to [our GraphQL Playground](https://api.transcend.io/graphql) and enter a query like

```gql
query {
  searchCatalogs(input: { text: "slack", limit: 25 }) {
    catalogs {
      integrationName
      description
      formConfigs {
        passportName
        type
        formItems {
          name
          type
          isPlaintext
        }
      }
      promptEmailTemplateId
      promptAVendorEmailAddress
      isPromptAVendorCompatible
      dataPointsCustomizable
      allowedActions
    }
  }
}
```

to search for integration metadata based on a title substring. Make sure you are logged into [your Organization's admin-dashboard](https://app.transcend.io/login) to have credentials on the GraphQL Playground.

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `type` (String) Type of silo

### Optional

- `content_classification_plugin` (Block List, Max: 1) Configuration for the Content Classification plugin for data silos. To be used in conjunction with the Schema Discovery plugin. (see [below for nested schema](#nestedblock--content_classification_plugin))
- `data_point_discovery_plugin` (Block List, Max: 1) [DEPRECATED] Configuration for the Data Point discovery plugin for data silos. (see [below for nested schema](#nestedblock--data_point_discovery_plugin))
- `data_silo_discovery_plugin` (Block List, Max: 1) Configuration for the Data Silo discovery plugin for data silos. (see [below for nested schema](#nestedblock--data_silo_discovery_plugin))
- `description` (String) The description of the data silo
- `headers` (Block List) Custom headers to include in outbound webhook (see [below for nested schema](#nestedblock--headers))
- `is_live` (Boolean) Whether the data silo should be live
- `notify_email_address` (String) The email address that should be notified whenever new requests are made
- `outer_type` (String) The catalog name responsible for the cosmetics of the integration (name, description, logo, email fields)
- `owner_emails` (List of String) The emails of the users to assign as owners of this data silo. These emails must have matching users on Transcend.
- `plaintext_context` (Block Set) This is where you put non-secretive values that go in the form when connecting a data silo (see [below for nested schema](#nestedblock--plaintext_context))
- `schema_discovery_plugin` (Block List, Max: 1) Configuration for the Schema Discovery plugin for data silos. (see [below for nested schema](#nestedblock--schema_discovery_plugin))
- `secret_context` (Block Set) This is where you put values that go in the form when connecting a data silo. In general, most form values are secret context. (see [below for nested schema](#nestedblock--secret_context))
- `skip_connecting` (Boolean) If true, the data silo will be left unconnected. When false, the provided credentials will be tested against a live environment
- `title` (String) The title of the data silo
- `url` (String) The URL of the server to post to if a server silo

### Read-Only

- `aws_external_id` (String) The external ID for the AWS IAM Role for AWS data silos
- `connection_state` (String) The current state of the integration
- `has_avc_functionality` (Boolean) Whether the data silo supports automated vendor coordination
- `id` (String) The ID of this resource.
- `link` (String) The link to the data silo

<a id="nestedblock--content_classification_plugin"></a>
### Nested Schema for `content_classification_plugin`

Optional:

- `enabled` (Boolean) State to toggle plugin to
- `schedule_frequency_minutes` (Number) The updated frequency with which we should schedule this plugin, in milliseconds
- `schedule_start_at` (String) The updated start time when we should start scheduling this plugin, in ISO format

Read-Only:

- `id` (String) The ID of this resource.
- `last_enabled_at` (String) The date at which this data silo was last enabled


<a id="nestedblock--data_point_discovery_plugin"></a>
### Nested Schema for `data_point_discovery_plugin`

Optional:

- `enabled` (Boolean) State to toggle plugin to
- `schedule_frequency_minutes` (Number) The updated frequency with which we should schedule this plugin, in milliseconds
- `schedule_start_at` (String) The updated start time when we should start scheduling this plugin, in ISO format

Read-Only:

- `id` (String) The ID of this resource.
- `last_enabled_at` (String) The date at which this data silo was last enabled


<a id="nestedblock--data_silo_discovery_plugin"></a>
### Nested Schema for `data_silo_discovery_plugin`

Optional:

- `enabled` (Boolean) State to toggle plugin to
- `schedule_frequency_minutes` (Number) The updated frequency with which we should schedule this plugin, in milliseconds
- `schedule_start_at` (String) The updated start time when we should start scheduling this plugin, in ISO format

Read-Only:

- `id` (String) The ID of this resource.
- `last_enabled_at` (String) The date at which this data silo was last enabled


<a id="nestedblock--headers"></a>
### Nested Schema for `headers`

Required:

- `name` (String) The name of the custom header
- `value` (String, Sensitive) The value of the custom header

Optional:

- `is_secret` (Boolean) When true, the value of this header will be considered sensitive


<a id="nestedblock--plaintext_context"></a>
### Nested Schema for `plaintext_context`

Required:

- `name` (String) The name of the plaintext input
- `value` (String) The value of the plaintext input


<a id="nestedblock--schema_discovery_plugin"></a>
### Nested Schema for `schema_discovery_plugin`

Optional:

- `enabled` (Boolean) State to toggle plugin to
- `schedule_frequency_minutes` (Number) The updated frequency with which we should schedule this plugin, in milliseconds
- `schedule_start_at` (String) The updated start time when we should start scheduling this plugin, in ISO format

Read-Only:

- `id` (String) The ID of this resource.
- `last_enabled_at` (String) The date at which this data silo was last enabled


<a id="nestedblock--secret_context"></a>
### Nested Schema for `secret_context`

Required:

- `name` (String) The name of the input
- `value` (String, Sensitive) The value of the input in plaintext

## Import

Import is supported using the following syntax:

```shell
terraform import transcend_data_silo.silo <data_silo_id_from_silo_url>
```