<script lang="ts">
	import ChevronDown from 'lucide-svelte/icons/chevron-down';
	import {
		type ColumnDef,
		type ColumnFiltersState,
		type PaginationState,
		type RowSelectionState,
		type SortingState,
		type VisibilityState,
		getCoreRowModel,
		getFilteredRowModel,
		getPaginationRowModel,
		getSortedRowModel
	} from '@tanstack/table-core';
	import { createRawSnippet } from 'svelte';
	import DataTableSortingButton from './data-table/data-table-sorting-button.svelte';
	import DataTableActions from './data-table/data-table-actions.svelte';
	import * as Table from '$lib/components/ui/table/index.js';
	import { Button } from '$lib/components/ui/button/index.js';
	import {
		FlexRender,
		createSvelteTable,
		renderComponent,
		renderSnippet
	} from '$lib/components/ui/data-table/index.js';
	import type { profile } from '$lib/wailsjs/go/models';
	import moment from 'moment';
	let {
		configs,
		onProfileDelete,
		onProfileEdit,
		onProfileRun,
	}: {configs: profile.Config[], onProfileDelete: (id: number) => void, onProfileEdit: (id: number) => void, onProfileRun: (id: number) => void} = $props();
	$effect(()=>{
		console.log(configs)
	})
	const columns: ColumnDef<profile.Config>[] = [
		{
			accessorKey: 'id',
			accessorFn: (row, index) => {
				return row.id
			},
			header: ({ column }) => {
				return renderComponent(DataTableSortingButton, {
					text: 'ID',
					onclick: () => column.toggleSorting(column.getIsSorted() === 'asc')
				})
			},
			cell: ({ row }) => {
				const idSnippet = createRawSnippet<[number]>((getID) => {
					return {
						render: () => `${getID()}`
					};
				});
				return renderSnippet(idSnippet, row.getValue('id'));
			}
		},
		{
			accessorKey: 'name',
			accessorFn: (row) => {
				return row.request.name
			},
			header: ({ column }) =>
				renderComponent(DataTableSortingButton, {
					text: 'Name',
					onclick: () => column.toggleSorting(column.getIsSorted() === 'asc')
				}),
			cell: ({ row }) => {
				const nameSnippet = createRawSnippet<[string]>((getName) => {
					return {
						render: () => `${getName()}`
					};
				});
				return renderSnippet(nameSnippet, row.getValue('name'));
			}
		},
		{
			accessorKey: 'proxy',
			accessorFn: (row)=>{
				return row.fingerprint.fingerprint.ipv4
			},
			header: ({ column }) =>
				renderComponent(DataTableSortingButton, {
					text: 'Proxy',
					onclick: () => column.toggleSorting(column.getIsSorted() === 'asc')
				}),
			cell: ({ row }) => {
				const nameSnippet = createRawSnippet<[string]>((getProxy) => {
					return {
						render: () => `${getProxy()}`
					};
				});
				return renderSnippet(nameSnippet, row.getValue('proxy'));
			}
		},
		{
			accessorKey: 'country',
			accessorFn: (row) =>{
				return row.fingerprint.country + " / " + row.fingerprint.fingerprint.config.timezone.replaceAll("/", " / ")
			},
			header: ({ column }) =>
				renderComponent(DataTableSortingButton, {
					text: 'Country',
					onclick: () => column.toggleSorting(column.getIsSorted() === 'asc')
				}),
			cell: ({ row }) => {
				const nameSnippet = createRawSnippet<[string]>((getCountry) => {
					return {
						render: () => `${getCountry()}`
					};
				});
				return renderSnippet(nameSnippet, row.getValue('country'));
			}
		},
		{
			accessorKey: 'gpu',
			accessorFn: (row) => {
				return row.fingerprint.fingerprint.config['webGl:renderer']
			},
			header: ({ column }) =>
				renderComponent(DataTableSortingButton, {
					text: 'GPU',
					onclick: () => column.toggleSorting(column.getIsSorted() === 'asc')
				}),
			cell: ({ row }) => {
				const gpuSnippet = createRawSnippet<[string]>((getGPU) => {
					return {
						render: () => `${getGPU()}`
					};
				});
				return renderSnippet(gpuSnippet, row.getValue('gpu'));
			}
		},
		{
			accessorKey: 'created',
			accessorFn: (row) => {
				return row.fingerprint.fingerprint.created_at
			},
			header: ({ column }) =>
				renderComponent(DataTableSortingButton, {
					text: 'Created At',
					onclick: () => column.toggleSorting(column.getIsSorted() === 'asc')
				}),
			cell: ({ row }) => {
				const gpuSnippet = createRawSnippet<[string]>((getCreatedAt) => {
					return {
						render: () => `${moment(getCreatedAt()).format("MMM Do YY")}`
					};	
				});
				return renderSnippet(gpuSnippet, row.getValue('created'));
			}
		},
		{
			accessorKey: 'lastUsed',
			accessorFn: (row) => {
				return row.fingerprint.fingerprint.accessed_at
			},
			header: ({ column }) =>
				renderComponent(DataTableSortingButton, {
					text: 'Last Used',
					onclick: () => column.toggleSorting(column.getIsSorted() === 'asc')
				}),
			cell: ({ row }) => {
				const snippet = createRawSnippet<[string]>((getLastUsed) => {
					return {
						render: () => `${getLastUsed() == null ? "Never" : moment(getLastUsed()).format("MMM Do YY")}`
					};
				});
				return renderSnippet(snippet, row.getValue('lastUsed'));
			}
		},
		{
			id: 'actions',
			enableHiding: false,
			header: ({ column }) =>
				renderComponent(DataTableSortingButton, {
					text: 'Actions'
				}),
			cell: ({ row }) => renderComponent(DataTableActions, { 
				id: row.original.id,
				onDelete: onProfileDelete,
				onEdit: onProfileEdit,
				onRun: onProfileRun
			})
		}
	];

	let pagination = $state<PaginationState>({ pageIndex: 0, pageSize: 10 });
	let sorting = $state<SortingState>([]);
	let columnFilters = $state<ColumnFiltersState>([]);
	let rowSelection = $state<RowSelectionState>({});
	let columnVisibility = $state<VisibilityState>({});

	$effect(()=>{
		console.log(configs)
	})
	const table = createSvelteTable({
		get data() {
			return configs;
		},
		columns,
		state: {
			get pagination() {
				return pagination;
			},
			get sorting() {
				return sorting;
			},
			get columnVisibility() {
				return columnVisibility;
			},
			get rowSelection() {
				return rowSelection;
			},
			get columnFilters() {
				return columnFilters;
			}
		},
		getCoreRowModel: getCoreRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		onPaginationChange: (updater) => {
			if (typeof updater === 'function') {
				pagination = updater(pagination);
			} else {
				pagination = updater;
			}
		},
		onSortingChange: (updater) => {
			if (typeof updater === 'function') {
				sorting = updater(sorting);
			} else {
				sorting = updater;
			}
		},
		onColumnFiltersChange: (updater) => {
			if (typeof updater === 'function') {
				columnFilters = updater(columnFilters);
			} else {
				columnFilters = updater;
			}
		},
		onColumnVisibilityChange: (updater) => {
			if (typeof updater === 'function') {
				columnVisibility = updater(columnVisibility);
			} else {
				columnVisibility = updater;
			}
		},
		onRowSelectionChange: (updater) => {
			if (typeof updater === 'function') {
				rowSelection = updater(rowSelection);
			} else {
				rowSelection = updater;
			}
		}
	});

</script>

<div class="flex h-full w-full flex-col px-4 bg-black/60 z-10">
	<div class="flex rounded-md border">
		<Table.Root>
			<Table.Header>
				{#each table.getHeaderGroups() as headerGroup (headerGroup.id)}
					<Table.Row>
						{#each headerGroup.headers as header (header.id)}
							<Table.Head class="[&:has([role=checkbox])]:pl-3">
								{#if !header.isPlaceholder}
									<FlexRender
										content={header.column.columnDef.header}
										context={header.getContext()}
									/>
								{/if}
							</Table.Head>
						{/each}
					</Table.Row>
				{/each}
			</Table.Header>
			<Table.Body>
				{#each table.getRowModel().rows as row (row.id)}
					<Table.Row data-state={row.getIsSelected() && 'selected'}>
						{#each row.getVisibleCells() as cell (cell.id)}
							<Table.Cell class="[&:has([role=checkbox])]:pl-3">
								<FlexRender content={cell.column.columnDef.cell} context={cell.getContext()} />
							</Table.Cell>
						{/each}
					</Table.Row>
				{:else}
					<Table.Row>
						<Table.Cell colspan={columns.length} class="h-24 text-center">No results.</Table.Cell>
					</Table.Row>
				{/each}
			</Table.Body>
		</Table.Root>
	</div>
	<div class="mb-4 mt-auto flex items-center justify-end space-x-2 pt-4">
		<div class="text-muted-foreground flex-1 text-sm">
			Total {table.getFilteredRowModel().rows.length} profiles.
		</div>
		<div class="space-x-2 z-50">
			<Button
				variant="outline"
				size="sm"
				onclick={() => table.previousPage()}
				disabled={!table.getCanPreviousPage()}
			>
				Previous
			</Button>
			<Button
				variant="outline"
				size="sm"
				onclick={() => table.nextPage()}
				disabled={!table.getCanNextPage()}
			>
				Next
			</Button>
		</div>
	</div>
</div>
