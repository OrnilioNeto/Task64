import type {Prefixes} from './types'

const TASK64_PREFIXES: Prefixes = {
	label: '*',
	project: '+',
	priority: '!',
	assignee: '@',
}

const TODOIST_PREFIXES: Prefixes = {
	label: '@',
	project: '#',
	priority: '!',
	assignee: '+',
}

export enum PrefixMode {
	Disabled = 'disabled',
	Default = 'vikunja',
	Todoist = 'todoist',
}

export const PREFIXES = {
	[PrefixMode.Disabled]: undefined,
	[PrefixMode.Default]: TASK64_PREFIXES,
	[PrefixMode.Todoist]: TODOIST_PREFIXES,
}
