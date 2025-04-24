terraform {
  required_providers {
    devops = {
      source = "liatr.io/terraform/devops-bootcamp"
    }
  }
}

provider "devops" {
  host     = "http://localhost:8080"
}

// data "devops_engineer" "edu" {}

resource "devops_engineer" "edu" {
  name  = "LangBangGang"
  email = "pingponggod@liatriolife.com"
}

resource "devops_engineer" "edu2" {
  name  = "Your_Mom_loves_me"
  email = "yourmom@myhouseagain.com"
}

resource "devops_engineer" "test" {
  name = "testie_mctestface"
  email = "testie@liatriolife.com"
}

data "devops_engineer" "test" {}

output "edu_engineer" {
  value = devops_engineer.edu
}

output "edu_engineer2" {
  value = devops_engineer.edu2
}

output "test_engineer" {
  value = devops_engineer.test
}

output "test_engineer_name" {
  value = data.devops_engineer.test.name
}
