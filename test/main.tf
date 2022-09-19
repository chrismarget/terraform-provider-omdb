terraform {
  required_providers {
    omdb = {
      source = "github.com/chrismarget/omdb"
    }
  }
}

variable "api_key" {}

provider "omdb" {
  // set with shell command: export TF_VAR_api_key="xxxxxx"
  api_key = var.api_key
}

data "omdb_film_by_id" "terminator" {
  imdb_id = "tt0088247"
}

output "terminator" {
  value = data.omdb_film_by_id.terminator
}
