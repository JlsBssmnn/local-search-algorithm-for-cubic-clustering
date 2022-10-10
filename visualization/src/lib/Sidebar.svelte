<script lang="ts">
  import type { Evaluation as EvaluationType } from "src/wasmInterface/types";

	import { createEventDispatcher } from "svelte";
	import Evaluation from "./Evaluation.svelte";
	import Settings from "./Settings.svelte";

	interface DisplayDataState {
		numOfPlanes: number;
		pointsPerPlane: number;
		mean: number;
		stddev: number;
	}

	export let evaluation: EvaluationType;
	const dispatch = createEventDispatcher();
	let selectedNumOfPlanes: number;

	function datagenEvent(event) {
		selectedNumOfPlanes = event.detail.numOfPlanes;
		dispatch("display", event.detail);
	}

	function partitionEvent(event) {
		dispatch("partitionEvent", event.detail);
	}

	const dataSettings = [
		{
			name: "Number of planes",
			field: "numOfPlanes",
			defaultValue: 2,
			minValue: 1,
			maxValue: 10,
		},
		{
			name: "Points per plane",
			field: "pointsPerPlane",
			defaultValue: 5,
			minValue: 1,
			maxValue: 100,
		},
		{
			name: "Noise mean",
			field: "mean",
			defaultValue: 0,
			minValue: -1,
			maxValue: 1,
			step: 0.0001,
		},
		{
			name: "Noise standard deviation",
			field: "stddev",
			defaultValue: 0.0001,
			minValue: 0,
			maxValue: 0.5,
			step: 0.00001,
		},
	];

	const algorithmSettings = [
		{
			name: "Algorithm",
			field: "algorithm",
			defaultValue: "GreedyJoining",
			options: ["GreedyJoining", "GreedyMoving"],
		},
		{
			name: "Threshold",
			field: "threshold",
			defaultValue: 0.0005,
			minValue: 0,
			maxValue: 1,
			step: 0.00001,
		},
		{
			name: "Amplification",
			field: "amplification",
			defaultValue: 1,
			minValue: 1,
			maxValue: 200,
			step: 1,
		},
	];
</script>

<div class="container">
	<Settings
		title="Data generation"
		items={dataSettings}
		submissionText="Display data"
		submissionEventName="datagenEvent"
		on:datagenEvent={datagenEvent}
	/>
	<Settings
		title="Partitioning algorithm"
		items={algorithmSettings}
		submissionText="Compute planes"
		submissionEventName="partitionEvent"
		on:partitionEvent={partitionEvent}
	/>
	{#if evaluation != null}
		<Evaluation  evaluation={evaluation} numOfPlanes={selectedNumOfPlanes}/>
	{/if}
</div>

<style>
	.container {
		display: flex;
		flex-direction: column;
		justify-content: center;
		height: 100%;
		gap: 30px;
		max-width: 600px;
	}
</style>
