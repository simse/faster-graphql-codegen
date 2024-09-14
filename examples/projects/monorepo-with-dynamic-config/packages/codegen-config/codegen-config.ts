import type { CodegenConfig } from '@graphql-codegen/cli'

const createConfig = (): CodegenConfig => ({
schema: '../../apps/graphql-server/schema.graphql',
documents: null,
generates: {
  '__generated__/bruh.ts': {
    plugins: ['typescript']
  }
}
})

export default createConfig