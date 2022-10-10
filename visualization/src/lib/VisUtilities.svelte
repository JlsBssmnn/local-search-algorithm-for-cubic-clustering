<script lang="ts">
	import type { VisualizationStatus } from "./types";

	export let visualizationStatus: VisualizationStatus;
	export let showPlanes: boolean;
	export let colorizePoints: boolean;
</script>

<div class="border">
	<div class="utilityBox">
		<div class="item" title="Show computed planes">
			<label>
				<input bind:checked={showPlanes} type="checkbox" />
				Show planes</label
			>
		</div>
		<div class="item" title="Color the points such that you can see the correct partitioning">
			<label>
				<input bind:checked={colorizePoints} type="checkbox" />
				Colorize points</label
			>
		</div>
		<div class="item" title="Shows whether or not the algorithms is currently computing a partitioning">
			<div>
				Visualization status:&nbsp;<div
					style={`color: ${visualizationStatus == "ready" ? "green" : "red"}; display: inline;`}
				>
					{visualizationStatus}
				</div>
			</div>
		</div>
	</div>
</div>

<style>
	.utilityBox {
		display: inline-flex;
		justify-content: center;
		margin-top: 5px;
		border: 1px solid black;
		gap: 30px;
		padding: 5px;
	}

	.border {
		display: grid;
		place-items: center;
	}

	.item {
		display: flex;
		justify-content: center;
	}
</style>
