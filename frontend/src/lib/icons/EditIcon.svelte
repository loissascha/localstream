<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Edit' }: Props = $props();
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
	<path d="M4 17.25V20h2.75l8.1-8.1-2.75-2.75-8.1 8.1Z" />
	<path
		d="M13.16 8.84 15.9 11.6l1.94-1.94a1.94 1.94 0 0 0 0-2.74l-1.1-1.1a1.94 1.94 0 0 0-2.74 0l-1.84 1.84Z"
	/>
</svg>
