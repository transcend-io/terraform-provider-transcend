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
        Action    = "sts:AssumeRole"
        Effect    = "Allow"
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
  depends_on = [aws_iam_role.iam_role]
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