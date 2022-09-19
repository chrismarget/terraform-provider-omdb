---
page_title: "Provider: OMDb"
description: |-
  The OMDB provider facilitates interactions with the Open Movie Database
---

# OMDB Provider

The OMDb provider facilitates interactions with the Open Movie Database.

Who would want such a thing? Me, for use as a template when working on more
complicated Terraform providers.

Really, this thing exists only to hone release process, documentation
generation, experiment with new versions of the Terraform Provider Framework
and the like. It's not a real provider.

Use the navigation to the left to read about the available resources.

## Configuration

```terraform
provider "omdb" {
  // set with shell command: export TF_VAR_api_key="xxxxxx"
  api_key = var.api_key
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `api_key` (String) A free OMDb API key can be quickly generated [here](https://www.omdbapi.com/apikey.aspx).

## Environment Variables

```
something something something
```