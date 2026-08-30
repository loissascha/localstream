<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Delete' }: Props = $props();
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
		d="M9 3.75A1.75 1.75 0 0 0 7.25 5.5V6H4.5a.75.75 0 0 0 0 1.5h.84l.86 11.17a2.25 2.25 0 0 0 2.24 2.08h7.12a2.25 2.25 0 0 0 2.24-2.08l.86-11.17h.84a.75.75 0 0 0 0-1.5h-2.75v-.5A1.75 1.75 0 0 0 15 3.75H9Zm6.25 2.25h-6.5v-.5a.25.25 0 0 1 .25-.25h6a.25.25 0 0 1 .25.25V6Z"
	/>
	<path
		d="M10 9.25a.75.75 0 0 1 .75.75v6a.75.75 0 0 1-1.5 0v-6a.75.75 0 0 1 .75-.75Zm4 0a.75.75 0 0 1 .75.75v6a.75.75 0 0 1-1.5 0v-6a.75.75 0 0 1 .75-.75Z"
	/>
</svg>
