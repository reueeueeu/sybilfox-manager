<script lang="ts" module>
	import Frame from 'lucide-svelte/icons/frame';
	import Cog from 'lucide-svelte/icons/cog';

	// This is sample data.
	const data = {
		projects: [
			{
				name: 'Profiles',
				icon: Frame
			},
			{
				name: 'Settings',
				icon: Cog
			}
		]
	};
</script>

<script lang="ts">
	import NavProjects from '$lib/components/nav-projects.svelte';
	import * as Sidebar from '$lib/components/ui/sidebar/index.js';
	import type { ComponentProps } from 'svelte';

	let {
		ref = $bindable(null),
		collapsible = 'icon',
		selectedMenu = $bindable(data.projects[0].name),
		...restProps
	}: ComponentProps<typeof Sidebar.Root> & { selectedMenu: string } = $props();
</script>

<Sidebar.Root bind:ref {collapsible} {...restProps}>
	<Sidebar.Content>
		<NavProjects bind:selected={selectedMenu} projects={data.projects} />
	</Sidebar.Content>
	<Sidebar.Rail />
</Sidebar.Root>
