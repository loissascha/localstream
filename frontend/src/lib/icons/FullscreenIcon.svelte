<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Fullscreen' }: Props = $props();
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
		d="M4 10V5a1 1 0 0 1 1-1h5v2H6v4H4Zm10-6h5a1 1 0 0 1 1 1v5h-2V6h-4V4ZM4 14h2v4h4v2H5a1 1 0 0 1-1-1v-5Zm14 0h2v5a1 1 0 0 1-1 1h-5v-2h4v-4Z"
	/>
</svg>
