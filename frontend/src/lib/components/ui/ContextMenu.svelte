<script lang="ts">
	import {
		autoUpdate,
		computePosition,
		flip,
		offset as floatingOffset,
		shift
	} from '@floating-ui/dom';
	import type { Placement, Strategy, VirtualElement } from '@floating-ui/dom';
	import type { Snippet } from 'svelte';
	import { tick } from 'svelte';
	import { scale } from 'svelte/transition';

	interface Props {
		children: Snippet;
		items: Snippet<[() => void]>;
		disabled?: boolean;
		closeOnItemClick?: boolean;
		placement?: Placement;
		offset?: number;
		className?: string;
		menuClassName?: string;
	}

	let {
		children,
		items,
		disabled = false,
		closeOnItemClick = true,
		placement = 'bottom-start',
		offset = 6,
		className = '',
		menuClassName = ''
	}: Props = $props();

	let open = $state(false);
	let container: HTMLDivElement | undefined = $state();
	let menu: HTMLDivElement | undefined = $state();
	let cursorX = $state(0);
	let cursorY = $state(0);
	let menuX = $state(0);
	let menuY = $state(0);
	let menuStrategy = $state<Strategy>('fixed');

	const menuId = `context-menu-${Math.random().toString(36).slice(2)}`;
	let cleanupAutoUpdate: (() => void) | undefined;

	function createVirtualReference(): VirtualElement {
		return {
			getBoundingClientRect() {
				return {
					width: 0,
					height: 0,
					x: cursorX,
					y: cursorY,
					top: cursorY,
					right: cursorX,
					bottom: cursorY,
					left: cursorX
				};
			}
		};
	}

	function getMenuItems() {
		return Array.from(menu?.querySelectorAll<HTMLElement>('[role="menuitem"]') ?? []);
	}

	function focusFirstItem() {
		const [firstItem] = getMenuItems();
		firstItem?.focus();
	}

	async function updatePosition() {
		if (!menu) return;

		const { x, y, strategy } = await computePosition(createVirtualReference(), menu, {
			placement,
			strategy: 'fixed',
			middleware: [floatingOffset(offset), flip(), shift({ padding: 8 })]
		});

		menuX = x;
		menuY = y;
		menuStrategy = strategy;
	}

	function cleanupPositioning() {
		cleanupAutoUpdate?.();
		cleanupAutoUpdate = undefined;
	}

	function closeMenu() {
		if (!open) return;

		open = false;
		cleanupPositioning();
	}

	async function openMenu(event: MouseEvent) {
		if (disabled) return;

		event.preventDefault();
		cursorX = event.clientX;
		cursorY = event.clientY;
		open = true;

		await tick();
		if (!menu) return;

		cleanupPositioning();
		cleanupAutoUpdate = autoUpdate(createVirtualReference(), menu, () => {
			void updatePosition();
		});

		await updatePosition();
		focusFirstItem();
	}

	function handleDocumentClick(event: MouseEvent) {
		if (!open || !container) return;

		const target = event.target;
		if (target instanceof Node && !container.contains(target)) {
			closeMenu();
		}
	}

	function handleDocumentKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape' && open) {
			event.preventDefault();
			closeMenu();
		}
	}

	function handleMenuClick(event: MouseEvent) {
		if (!closeOnItemClick) return;
		if (!(event.target instanceof Element)) return;
		if (!event.target.closest('[role="menuitem"]')) return;

		closeMenu();
	}

	function handleMenuKeydown(event: KeyboardEvent) {
		const items = getMenuItems();
		if (items.length === 0) return;

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
			case 'Tab':
				closeMenu();
				break;
		}
	}

	$effect(() => {
		return () => {
			cleanupPositioning();
		};
	});
</script>

<svelte:document onclick={handleDocumentClick} onkeydown={handleDocumentKeydown} />

<!-- svelte-ignore a11y_no_static_element_interactions -->
<div bind:this={container} class={className} oncontextmenu={openMenu}>
	{@render children()}

	{#if open}
		<div
			bind:this={menu}
			id={menuId}
			role="menu"
			tabindex={-1}
			aria-orientation="vertical"
			onclick={handleMenuClick}
			onkeydown={handleMenuKeydown}
			oncontextmenu={(event) => {
				event.preventDefault();
			}}
			transition:scale={{ duration: 140, start: 0.95 }}
			style={`position: ${menuStrategy}; left: ${menuX}px; top: ${menuY}px;`}
			class={`z-50 flex w-max min-w-40 flex-col gap-1 overflow-hidden rounded-md bg-neutral-800 text-white shadow-lg ${menuClassName}`}
		>
			{@render items(closeMenu)}
		</div>
	{/if}
</div>
