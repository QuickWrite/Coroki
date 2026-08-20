import * as esbuild from "esbuild";

const watch = process.argv.includes("--watch");

const options = {
    entryPoints: ["web/app.ts", "web/app.css"],
    outdir: "web/dist",
    bundle: true,
    sourcemap: true,
    minify: !watch,
};

if (watch) {
    const ctx = await esbuild.context(options);
    await ctx.watch();
    console.log("Watching frontend files...");
} else {
    await esbuild.build(options);
}
