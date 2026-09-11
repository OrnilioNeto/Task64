import type {OAuthTokens} from '@/types/desktop'

export function isDesktopApp(): boolean {
	return !!window.task64Desktop?.isDesktop
}

export function startDesktopOAuthLogin(apiUrl: string): Promise<void> {
	return window.task64Desktop!.startOAuthLogin(apiUrl)
}

export function listenForDesktopOAuthTokens(callback: (tokens: OAuthTokens) => void): void {
	window.task64Desktop!.onOAuthTokens(callback)
}

export function listenForDesktopOAuthError(callback: (error: string) => void): void {
	window.task64Desktop!.onOAuthError(callback)
}

export function refreshDesktopToken(apiUrl: string, refreshToken: string): Promise<OAuthTokens> {
	return window.task64Desktop!.refreshToken(apiUrl, refreshToken)
}
