<script lang="ts">
	import Plotly from "plotly.js-dist-min";
	import { onMount } from "svelte";

	let div;
	export let data: any[];
	let plotData: any[] = [];

	$: {
		data = data;
		if (div) {
			updatePlot();
		}
	}

	// Plotly needs the same data object for updating
	// but svelte needs a new object for updating
	function updatePlot() {
		plotData.splice(data.length);
		data.forEach((vis, i) => {
			plotData[i] = vis;
		});
		Plotly.redraw(div);
	}

	onMount(() => {
		let layout = {
			margin: {
				l: 0,
				r: 0,
				b: 0,
				t: 0,
			},
			width: 900,
			height: 800,
		};

		Plotly.newPlot(div, plotData, layout);
	});
</script>

<div bind:this={div} />
