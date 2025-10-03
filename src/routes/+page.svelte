<script>
  const API =
    (import.meta.env.DEV ? 'http://localhost:3000' : '') + '/api/close-loop';

  let P_tool_6 = [25, 23, 24];
  let P_tool_f = [0.177, 0.884, -0.433];
  let S6_f     = [-0.153, 0.459, 0.875];
  let a67_f    = [-0.153, 0.459, 0.875];

  let out = null;
  let loading = false;
  let err = "";

  async function compute() {
    loading = true; err = ""; out = null;
    try {
      const res = await fetch(API, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ P_tool_6, P_tool_f, S6_f, a67_f })
      });
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      out = await res.json();
    } catch (e) {
      err = (e && e.message) ? e.message : String(e);
    } finally {
      loading = false;
    }
  }
</script>

<style>
  .grid { display: grid; grid-template-columns: repeat(2, minmax(0,1fr)); gap: .75rem; }
  .card { border: 1px solid #e5e7eb; border-radius: .75rem; padding: 1rem; }
  .title { font-weight: 700; font-size: 1.1rem; margin-bottom: .5rem; }
  .row { display: grid; grid-template-columns: 1fr 1fr 1fr; gap: .5rem; margin-bottom: .5rem; }
  input { width: 100%; padding: .4rem; border: 1px solid #d1d5db; border-radius: .5rem; }
  button { padding: .6rem 1rem; border-radius: .6rem; border: 1px solid #111827; background: #111827; color: #fff; cursor: pointer; }
  button:hover { background: #0b0f1a; }
  .outputs span { display: block; margin-bottom: .3rem; }
  .error { color: #b91c1c; margin-top: .5rem; }
</style>

<h1>Inverse Kinematics : Close-the-Loop</h1>
<p>S₁ assumed to be world Z through the origin.</p>

<div class="grid">
  <div class="card">
    <div class="title">^6P<sub>tool</sub> (unused in minimal math)</div>
    <div class="row">
      <input type="number" step="any" bind:value={P_tool_6[0]}>
      <input type="number" step="any" bind:value={P_tool_6[1]}>
      <input type="number" step="any" bind:value={P_tool_6[2]}>
    </div>
  </div>

  <div class="card">
    <div class="title">^F P<sub>tool</sub></div>
    <div class="row">
      <input type="number" step="any" bind:value={P_tool_f[0]}>
      <input type="number" step="any" bind:value={P_tool_f[1]}>
      <input type="number" step="any" bind:value={P_tool_f[2]}>
    </div>
  </div>

  <div class="card">
    <div class="title">^F S6 (unused in minimal math)</div>
    <div class="row">
      <input type="number" step="any" bind:value={S6_f[0]}>
      <input type="number" step="any" bind:value={S6_f[1]}>
      <input type="number" step="any" bind:value={S6_f[2]}>
    </div>
  </div>

  <div class="card">
    <div class="title">^F a<sub>67</sub> (S7 direction)</div>
    <div class="row">
      <input type="number" step="any" bind:value={a67_f[0]}>
      <input type="number" step="any" bind:value={a67_f[1]}>
      <input type="number" step="any" bind:value={a67_f[2]}>
    </div>
  </div>
</div>

<div style="margin-top:1rem;">
  <button on:click={compute} disabled={loading}>{loading ? "Computing..." : "Compute"}</button>
  {#if err}<div class="error">Error: {err}</div>{/if}
</div>

{#if out}
  <div class="card" style="margin-top:1rem;">
    <div class="title">Results</div>
    <div class="outputs">
      <span>a<sub>71</sub> = {out.a71.toFixed(4)}</span>
      <span>S<sub>7</sub>  = {out.S7.toFixed(4)}</span>
      <span>S<sub>1</sub>  = {out.S1.toFixed(4)}</span>
      <span>α<sub>71</sub> = {out.al71_deg.toFixed(4)}°</span>
      <span>θ<sub>7</sub>  = {out.th7_deg.toFixed(4)}°</span>
      <span>γ<sub>1</sub>  = {out.gaml_deg.toFixed(4)}°</span>
    </div>
  </div>
{/if}
