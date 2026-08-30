<script lang="ts">
	import type { Snippet } from 'svelte';
	import { tick } from 'svelte';
	import { scale } from 'svelte/transition';

	type Anchor = 'left' | 'middle' | 'right';

	interface Props {
		children: Snippet;
		items: Snippet;
		anchor?: Anchor;
	}
	let { children, items, anchor = 'middle' }: Props = $props();

	let open = $state(false);
	let button: HTMLButtonElement | undefined = $state();
	let container: HTMLDivElement | undefined = $state();
	let menu: HTMLDivElement | undefined = $state();

	const menuId = `dropdown-${Math.random().toString(36).slice(2)}`;

	const anchorClasses: Record<Anchor, string> = {
		left: 'left-0 origin-top-left',
		middle: 'left-1/2 -translate-x-1/2 origin-top',
		right: 'right-0 origin-top-right'
	};

	function getMenuItems() {
		return Array.from(menu?.querySelectorAll<HTMLElement>('[role="menuitem"]') ?? []);
	}

	async function focusMenu() {
		await tick();
		menu?.focus();
	}

	function closeDropdown(restoreFocus = false) {
		if (!open) return;

		open = false;

		if (restoreFocus) {
			void tick().then(() => button?.focus());
		}
	}

	function handleDocumentClick(event: MouseEvent) {
		if (!open || !container) return;

		const target = event.target;
		if (target instanceof Node && !container.contains(target)) {
			closeDropdown();
		}
	}

	function handleDocumentKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && open) {
			event.preventDefault();
			closeDropdown(true);
		}
	}

	function handleMenuKeydown(event: KeyboardEvent) {
		const items = getMenuItems();
		const currentIndex = items.findIndex((item) => item === document.activeElement);

		switch (event.key) {
			case 'ArrowDown': {
				event.preventDefault();
				const nextIndex = currentIndex >= 0 ? (currentIndex + 1) % items.length : 0;
				items[nextIndex]?.focus();
				break;
			}
			case 'ArrowUp': {
				event.preventDefault();
				const nextIndex =
					currentIndex >= 0 ? (currentIndex - 1 + items.length) % items.length : items.length - 1;
				items[nextIndex]?.focus();
				break;
			}
			case 'Home':
				event.preventDefault();
				items[0]?.focus();
				break;
			case 'End':
				event.preventDefault();
				items.at(-1)?.focus();
				break;
		}
	}

	$effect(() => {
		if (open) {
			void focusMenu();
		}
	});
</script>

<svelte:document onclick={handleDocumentClick} onkeydown={handleDocumentKeydown} />

<div bind:this={container} class="relative inline-block">
	<button
		bind:this={button}
		type="button"
		aria-controls={menuId}
		aria-expanded={open}
		aria-haspopup="menu"
		onclick={() => {
			open = !open;
		}}
		class="cursor-pointer"
	>
		{@render children()}
	</button>
	{#if open}
		<div
			bind:this={menu}
			id={menuId}
			role="menu"
			tabindex={-1}
			aria-orientation="vertical"
			onclick={() => {
				closeDropdown();
			}}
			onkeydown={handleMenuKeydown}
			transition:scale={{ duration: 140, start: 0.95 }}
			class={`absolute top-full z-40 mt-2 flex w-max min-w-full flex-col gap-1 overflow-hidden rounded-md bg-neutral-800 shadow-lg ${anchorClasses[anchor]}`}
		>
			{@render items()}
		</div>
	{/if}
</div>
