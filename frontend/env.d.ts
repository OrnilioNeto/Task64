/// <reference types="vite/client" />
/// <reference types="vite-svg-loader" />
/// <reference types="@histoire/plugin-vue/components" />

interface ImportMetaEnv {
	readonly TASK64_API_URL?: string
	readonly TASK64_HTTP_PORT?: number
	readonly TASK64_HTTPS_PORT?: number

	readonly TASK64_SENTRY_ENABLED?: boolean
	readonly TASK64_SENTRY_DSN?: string

	readonly SENTRY_AUTH_TOKEN?: string
	readonly SENTRY_ORG?: string
	readonly SENTRY_PROJECT?: string

	readonly VITE_IS_ONLINE: boolean

	readonly VUE_DEVTOOLS_LAUNCH_EDITOR: VitePluginVueDevToolsOptions.launchEditor
}

interface ImportMeta {
	readonly env: ImportMetaEnv
}
