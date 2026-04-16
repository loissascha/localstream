<script lang="ts">
	type Props = {
		size?: number | string;
		class?: string;
		title?: string;
	};

	let { size = 24, class: className = '', title = 'Shows' }: Props = $props();
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
	<rect x="3" y="4" width="18" height="12" rx="2.5" />
	<rect x="4.5" y="5.5" width="15" height="9" rx="1.5" opacity="0.15" />
	<path d="M11.2 8.2a1 1 0 0 1 1.5-.86l3.2 1.9a1 1 0 0 1 0 1.72l-3.2 1.9a1 1 0 0 1-1.5-.86V8.2Z" />
	<rect x="4" y="18" width="4" height="2" rx="1" />
	<rect x="10" y="18" width="4" height="2" rx="1" />
	<rect x="16" y="18" width="4" height="2" rx="1" />
</svg>
