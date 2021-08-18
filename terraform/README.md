## terraform

This Terraform creates the necessary AWS resources to allow for GitHub Actions to push a new release up to S3 upon a new version being tagged. It only needs to be run once for initial creation of the buckets and policies, so no CI/CD is set up for this portion.

It hosts:
* S3 bucket
* S3 bucket policy
* IAM user
* IAM user policy