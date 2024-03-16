# Kheeper UI

Tools for authoring server-side components in Go.

## Component CSS & JS

The CSS for all registered components is bundled and served together.
The bundle will be served at `/css/<version>.css`.
The version is the first 7 characters of the sha256 hex digest of the bundle.

Likewise, JS for all components is bundled and served at `/js/<version>.js`.

A link to each bundle is injected into the head of every page.
Additionally, a Link preload header is emitted in the response for every page to enable [Early Hints](https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/103).

## Static

## Recipes

### Add 
