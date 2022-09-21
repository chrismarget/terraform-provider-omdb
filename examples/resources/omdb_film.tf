variable "api_key" {}

provider "omdb" {
  // set with shell command: export TF_VAR_api_key="xxxxxx"
  api_key = var.api_key
}

data "omdb_film_by_id" "my_favorite_film" {
  imdb_id = "tt0080455"
}

resource "omdb_film" "fav" {
  title = data.omdb_film_by_id.my_favorite_film.title
  year = data.omdb_film_by_id.my_favorite_film.year
}
