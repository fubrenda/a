# A

A project to produce an easy to use authority server


## Getting Started

This will install dependencies. You'll be able to run tests and
build the codebase afterwards.

```bash
make init
```

Run the tests

```bash
make test
```

To build

```bash
make build
```

## Notes

Filtering the wikidata dump with jq




```
time bzip2 -dc latest-all.json.bz2 |   jq -nc --stream   'fromstream(1|truncate_stream(inputs)) | select(.claims | keys_unsorted | any(contains("P244"), contains("P214"), contains("P4801"), contains("P1014"), contains("P486")))'

time curl "https://dumps.wikimedia.org/wikidatawiki/entities/latest-all.json.bz2" | bzip2 -dc | jq -nc --stream   'fromstream(1|truncate_stream(inputs)) | .claims.P244[].mainsnak.datavalue'
```
