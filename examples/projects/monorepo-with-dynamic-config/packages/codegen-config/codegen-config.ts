import type { CodegenConfig } from '@graphql-codegen/cli'

const createConfig = (): CodegenConfig => ({
schema: ['../../apps/graphql-server/schema.graphql'],
documents: null,
generates: {
  '__generated__/baseTypes.ts': {
    plugins: ['typescript'],
    preset: 'client',
  }
}
})

export default createConfig