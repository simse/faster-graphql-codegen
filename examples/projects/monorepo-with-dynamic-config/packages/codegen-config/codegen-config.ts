import type { CodegenConfig } from '@graphql-codegen/cli'
import path from 'node:path';

const createConfig = (): CodegenConfig => ({
schema: ['../../apps/graphql-server' + path.delimiter() + 'schema.graphql'],
documents: null,
generates: {
  '__generated__/baseTypes.ts': {
    plugins: ['typescript'],
    preset: 'client',
  }
}
})

export default createConfig