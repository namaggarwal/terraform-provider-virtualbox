terraform {
  required_version = ">= 0.13.0"
  required_providers {
    virtualbox = {
      version = "0.0.1"
      source   = "naman.io/namantest/virtualbox"
    }
  }
}

provider "virtualbox" {
}
