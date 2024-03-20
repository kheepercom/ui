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

## Tailwind

All components may use [Tailwind](https://tailwindcss.com/) classes.
In development the [Tailwind browser library](https://tailwindcss.com/docs/installation/play-cdn) provides support for all classes.

> [!Important]
> If a component requires a Tailwind class then it must implement the Examples interface and use the class at least once in an example.
> The optimized release CSS will only include Tailwind classes that appear in an example.

## Recommended Workflow

Use [air](https://github.com/cosmtrek/air?tab=readme-ov-file) to watch and rebuild your project on any changes.
