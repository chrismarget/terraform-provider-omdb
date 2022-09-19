---
page_title: "omdb_film_by_id Data Source - terraform-provider-omdb"
subcategory: ""
description: |-
  This Data Source returns details about a film by its IMDb ID.
---

# omdb_film_by_id (Data Source)

This Data Source returns details about a film by its IMDb ID.

## Example Usage

```terraform
data "omdb_film_by_id" "terminator" {
  imdb_id = "tt0088247"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `imdb_id` (String) Unique ID used by both OMDb and IMDb.

### Read-Only

- `Year` (String) Release year.
- `title` (String) Film title.