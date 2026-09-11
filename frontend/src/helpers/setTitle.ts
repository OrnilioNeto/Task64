export function setTitle(title : undefined | string) {
	document.title = (typeof title === 'undefined' || title === '')
		? 'Task64'
		: `${title} | Task64`
}
