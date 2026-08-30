<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Filter' }: Props = $props();
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
		d="M3.75 5.25A1.25 1.25 0 0 1 5 4h14a1.25 1.25 0 0 1 .98 2.03L14 13.5v5.25a.75.75 0 0 1-1.2.6l-2.5-1.88a.75.75 0 0 1-.3-.6V13.5L4.02 6.03a1.25 1.25 0 0 1-.27-.78Z"
	/>
</svg>
