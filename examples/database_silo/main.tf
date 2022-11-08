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
    schedule_now               = false
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
