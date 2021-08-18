# this terraform creates an iam user for github actions to upload 
# to S3 upon a new release
# --
# it must be run locally w/ profile "shannon"

provider "aws" {
  region  = "us-west-2"
  profile = "shannon"
}

terraform {
  backend "s3" {
    bucket  = "shannon-terraform"
    key     = "traverse/terraform.tfstate"
    region  = "us-west-2"
    profile = "shannon"
  }
}

resource "aws_s3_bucket" "traverse_bucket" {
  bucket = "lucassha-traverse-releases"
  acl    = "public-read"

  tags = {
    created_by_terraform = "true"
    owner                = "shannon"
  }
}

resource "aws_s3_bucket_policy" "show_secrets_bucket_policy" {
  bucket = aws_s3_bucket.traverse_bucket.id

  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "Allow Public Access to All Objects",
        "Effect" : "Allow",
        "Principal" : "*",
        "Action" : "s3:GetObject",
        "Resource" : "arn:aws:s3:::${aws_s3_bucket.traverse_bucket.bucket}/*"
      }
    ]
  })
}

resource "aws_iam_user" "s3" {
  name = "s3_traverse_user"
  path = "/system/"

  tags = {
    created_by_terraform = "true"
    owner                = "shannon"
  }
}

resource "aws_iam_user_policy" "s3_user_policy" {
  name = "s3_show_secrets_user_policy"
  user = aws_iam_user.s3.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "S3AllowToReleasesBucket",
      "Action": "s3:*",
      "Effect": "Allow",
      "Resource": [
        "arn:aws:s3:::${aws_s3_bucket.traverse_bucket.bucket}",
        "arn:aws:s3:::${aws_s3_bucket.traverse_bucket.bucket}/*"
      ]
    },
    {
      "Effect": "Allow",
      "Action": [
          "s3:ListAllMyBuckets"
      ],
      "Resource": "arn:aws:s3:::*"
    }
  ]
}
EOF    
}

