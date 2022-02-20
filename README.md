First go at a Terraform provider.

I was annoyed by the fact that the AWS provider DynamoDB table item resource expects the JSON to be formatted in DynamoDB JSON.  This is undesirable when you may have data in a variable format.

Also, I love the AWS SDK for Go, however there is no way to use it to convert JSON into DynamoDB format since they do not export those functions.  Please forgive the copy pasta from some of their unexported functions to achieve this.
