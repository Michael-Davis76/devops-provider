 terraform {
   required_providers {
     devops= {
       source = "liatr.io/terraform/devops-bootcamp"
     }
   }
 }

 provider "devops" {
   # example configuration here
  host     = "http://localhost:8080"
 }

data "devops_engineer" "example" {}
