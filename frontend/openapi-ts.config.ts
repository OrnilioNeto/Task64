import {defineConfig} from '@hey-api/openapi-ts'

const input = process.env.TASK64_OPENAPI_INPUT
if (!input) {
	throw new Error('TASK64_OPENAPI_INPUT must point to the generated temporary spec')
}

export default defineConfig({
	input,
	output: 'src/client/generated',
	plugins: [
		'@hey-api/typescript',
		'@hey-api/sdk',
		{
			name: '@hey-api/client-fetch',
			throwOnError: true,
		},
	],
})
