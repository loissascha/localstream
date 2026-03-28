<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Settings' }: Props = $props();
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
		fill-rule="evenodd"
		clip-rule="evenodd"
		d="M11.49 3.17a.75.75 0 0 1 1.02 0l.93.86a1.5 1.5 0 0 0 1.46.34l1.23-.4a.75.75 0 0 1 .92.46l.47 1.2c.21.53.66.95 1.2 1.12l1.24.4a.75.75 0 0 1 .5.89l-.24 1.27a1.5 1.5 0 0 0 .38 1.45l.87.9a.75.75 0 0 1 0 1.04l-.87.9a1.5 1.5 0 0 0-.38 1.45l.24 1.27a.75.75 0 0 1-.5.89l-1.24.4a1.5 1.5 0 0 0-1.2 1.12l-.47 1.2a.75.75 0 0 1-.92.46l-1.23-.4a1.5 1.5 0 0 0-1.46.34l-.93.86a.75.75 0 0 1-1.02 0l-.93-.86a1.5 1.5 0 0 0-1.46-.34l-1.23.4a.75.75 0 0 1-.92-.46l-.47-1.2a1.5 1.5 0 0 0-1.2-1.12l-1.24-.4a.75.75 0 0 1-.5-.89l.24-1.27a1.5 1.5 0 0 0-.38-1.45l-.87-.9a.75.75 0 0 1 0-1.04l.87-.9a1.5 1.5 0 0 0 .38-1.45l-.24-1.27a.75.75 0 0 1 .5-.89l1.24-.4a1.5 1.5 0 0 0 1.2-1.12l.47-1.2a.75.75 0 0 1 .92-.46l1.23.4a1.5 1.5 0 0 0 1.46-.34l.93-.86ZM12 9a3 3 0 1 0 0 6 3 3 0 0 0 0-6Z"
	/>
</svg>
