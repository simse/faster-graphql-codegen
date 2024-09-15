# Configuration

::: danger faster-graphql-codegen is in alpha
Config documentation may be incorrect, incomplete or out of date.
:::

Since `faster-graphql-codegen` is a drop-in replacement for `graphql-codegen`, it uses the same config format.

However, there are still some slight differences to be aware of, and these are documented here.

## `codegen.ts` format

Here's a basic `codegen.ts` file:
```ts
import type { CodegenConfig } from '@graphql-codegen/cli'

const config: CodegenConfig = {
  schema: 'schema.graphql',
  documents: null,
  generates: {
    '__generated__/baseTypes.ts': {
      plugins: ['typescript']
    }
  }
}

export default config
```

### `schema`
**Required.** Must be a string or an array of strings. Must point to a local file.

::: code-group

```ts [Single schema]
import type { CodegenConfig } from '@graphql-codegen/cli'

const config: CodegenConfig = {
  schema: 'schema.graphql', // [!code focus]
  documents: null,
  generates: {
    '__generated__/baseTypes.ts': {
      plugins: ['typescript']
    }
  }
}

export default config
```

```ts [Multiple schemas]
import type { CodegenConfig } from '@graphql-codegen/cli'

const config: CodegenConfig = {
  schema: [           // [!code focus]
    'base.graphql',   // [!code focus]
    'search.graphql', // [!code focus]
  ],                  // [!code focus]
  documents: null,
  generates: {
    '__generated__/baseTypes.ts': {
      plugins: ['typescript']
    }
  }
}

export default config
```

:::

### `generates`
An object containg outputs and their configurations. The object key is the name of the output file.

```ts
import type { CodegenConfig } from '@graphql-codegen/cli'

const config: CodegenConfig = {
  schema: 'schema.graphql',
  documents: null,
  generates: { // [!code focus]
    '__generated__/baseTypes.ts': { // [!code focus]
      plugins: ['typescript'] // [!code focus]
    } // [!code focus]
  } // [!code focus]
}

export default config
```

#### `plugins`
A `generates` config must contain a list of plugins. For available plugins please see the [plugin page](../plugins/index).

## Other formats
`faster-graphql-codegen` can also read this configuration from a `.yaml` or `.json` file.

::: code-group

```yaml [codegen.yml]
schema: ["schema.graphql"]
documents: []
overwrite: true
generates:
  'baseTypes.ts':
    plugins: [typescript]
```

```ts [codegen.ts]
import type { CodegenConfig } from '@graphql-codegen/cli'

const config: CodegenConfig = {
  schema: 'schema.graphql',
  documents: [],
  overwrite: true,
  generates: {
    '__generated__/baseTypes.ts': {
      plugins: ['typescript']
    }
  }
}

export default config
```

:::

## A note on performance

::: warning STATIC FILES LOAD FASTER
JSON and YAML config files are faster to load.
:::

Dynamic config files (JS and TS) are supported for compatibility purposes, but they are slower to load, because they have to be interpreted first.

If you are not using any Javascript features, please consider using a static format such as JSON or YAML.