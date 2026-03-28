<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'View' }: Props = $props();
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
		d="M12 5c4.7 0 8.35 2.95 9.85 6.32a1.5 1.5 0 0 1 0 1.36C20.35 16.05 16.7 19 12 19s-8.35-2.95-9.85-6.32a1.5 1.5 0 0 1 0-1.36C3.65 7.95 7.3 5 12 5Zm0 2c-3.76 0-6.76 2.3-8.2 5 1.44 2.7 4.44 5 8.2 5s6.76-2.3 8.2-5c-1.44-2.7-4.44-5-8.2-5Z"
	/>
	<circle cx="12" cy="12" r="3" />
</svg>
