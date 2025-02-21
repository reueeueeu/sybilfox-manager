<script lang="ts">
	import { useSidebar } from "$lib/components/ui/sidebar/context.svelte.js";
	import * as Sidebar from "$lib/components/ui/sidebar/index.js";
	import Ellipsis from "lucide-svelte/icons/ellipsis";

	let {
		projects,
		selected = $bindable(projects[0].name)
	}: {
		projects: {
			name: string;
			icon: any;
		}[];
		selected: string,
	} = $props();
	const sidebar = useSidebar();
</script>

<Sidebar.Group class="group-data-[collapsible=icon]:hidden">
	<Sidebar.GroupLabel>Sybilfox Manager</Sidebar.GroupLabel>
	<Sidebar.Menu class="mt-2">
		{#each projects as item (item.name)}
			<div class="{selected == item.name.toLowerCase() ? 'text-primary' : 'text-base'} cursor-pointer">
			<Sidebar.MenuItem>
				<Sidebar.MenuButton>
					{#snippet child({ props })}
						<a {...props}  onclick={()=>{
							if (selected == item.name.toLowerCase()) return;
							selected = item.name.toLowerCase()
						}}>
							<item.icon />
							<span>{item.name}</span>
						</a>
					{/snippet}
				</Sidebar.MenuButton>
			</Sidebar.MenuItem>
			</div>
		{/each}
	</Sidebar.Menu>
</Sidebar.Group>
