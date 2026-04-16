<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Logout' }: Props = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	viewBox="0 0 24 24"
	fill="currentColor"
	width={size}
	height={size}
	class={className}
	role="img"
	aria-hidden={title ? undefined : 'true'}
>
	{#if title}
		<title>{title}</title>
	{/if}
	<path
		d="M9 4.75A1.75 1.75 0 0 1 10.75 3h6.5A1.75 1.75 0 0 1 19 4.75v14.5A1.75 1.75 0 0 1 17.25 21h-6.5A1.75 1.75 0 0 1 9 19.25V16a.75.75 0 0 1 1.5 0v3.25c0 .14.11.25.25.25h6.5a.25.25 0 0 0 .25-.25V4.75a.25.25 0 0 0-.25-.25h-6.5a.25.25 0 0 0-.25.25V8a.75.75 0 0 1-1.5 0V4.75Z"
	/>
	<path
		d="M11.72 8.22a.75.75 0 0 1 1.06 0l3.25 3.25a.75.75 0 0 1 0 1.06l-3.25 3.25a.75.75 0 1 1-1.06-1.06L13.69 12l-1.97-1.97a.75.75 0 0 1 0-1.06Z"
	/>
	<path d="M4 12a.75.75 0 0 1 .75-.75h10.5a.75.75 0 0 1 0 1.5H4.75A.75.75 0 0 1 4 12Z" />
</svg>
